package server

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudBlockStorage() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"server_instance_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudBlockStorage(), fieldMap, dataSourceNcloudBlockStorageRead)
}

func dataSourceNcloudBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instances, err := getBlockStorageList(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudBlockStorage().Schema)
	}

	if err := ValidateOneResult(len(resources)); err != nil {
		return err
	}

	d.SetId(resources[0]["block_storage_no"].(string))
	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getBlockStorageList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*BlockStorage, error) {
	if config.SupportVPC {
		return getVpcBlockStorageList(d, config)
	}

	return getClassicBlockStorageList(d, config)
}

func getClassicBlockStorageList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*BlockStorage, error) {
	regionNo, err := conn.ParseRegionNoParameter(d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.GetBlockStorageInstanceListRequest{
		RegionNo:         regionNo,
		ServerInstanceNo: StringPtrOrNil(d.GetOk("server_instance_no")),
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.BlockStorageInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getClassicBlockStorageList", reqParams)

	resp, err := config.Client.Server.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicBlockStorageList", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicBlockStorageList", resp)

	var list []*BlockStorage
	for _, r := range resp.BlockStorageInstanceList {
		instance := &BlockStorage{
			BlockStorageInstanceNo:  r.BlockStorageInstanceNo,
			ServerInstanceNo:        r.ServerInstanceNo,
			ServerName:              r.ServerName,
			BlockStorageType:        r.BlockStorageType.Code,
			BlockStorageName:        r.BlockStorageName,
			BlockStorageSize:        ncloud.Int64(*r.BlockStorageSize / GIGABYTE),
			DeviceName:              r.DeviceName,
			BlockStorageProductCode: r.BlockStorageProductCode,
			Status:                  r.BlockStorageInstanceStatus.Code,
			Description:             r.BlockStorageInstanceDescription,
			DiskType:                r.DiskType.Code,
			DiskDetailType:          r.DiskDetailType.Code,
		}

		list = append(list, instance)
	}

	return list, nil
}

func getVpcBlockStorageList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*BlockStorage, error) {
	reqParams := &vserver.GetBlockStorageInstanceListRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: StringPtrOrNil(d.GetOk("server_instance_no")),
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.BlockStorageInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getVpcBlockStorageList", reqParams)

	resp, err := config.Client.Vserver.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcBlockStorage", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcBlockStorageList", resp)

	var list []*BlockStorage
	for _, r := range resp.BlockStorageInstanceList {
		instance := &BlockStorage{
			BlockStorageInstanceNo:  r.BlockStorageInstanceNo,
			ServerInstanceNo:        r.ServerInstanceNo,
			BlockStorageType:        r.BlockStorageType.Code,
			BlockStorageName:        r.BlockStorageName,
			BlockStorageSize:        ncloud.Int64(*r.BlockStorageSize / GIGABYTE),
			DeviceName:              r.DeviceName,
			BlockStorageProductCode: r.BlockStorageProductCode,
			Status:                  r.BlockStorageInstanceStatus.Code,
			Description:             r.BlockStorageDescription,
			DiskType:                r.BlockStorageDiskType.Code,
			DiskDetailType:          r.BlockStorageDiskDetailType.Code,
			ZoneCode:                r.ZoneCode,
			MaxIops:                 r.MaxIopsThroughput,
			IsEncryptedVolume:       r.IsEncryptedVolume,
			IsReturnProtection:      r.IsReturnProtection,
			VolumeType:              r.BlockStorageVolumeType.Code,
			HypervisorType:          r.HypervisorType.Code,
		}

		list = append(list, instance)
	}

	return list, nil
}
