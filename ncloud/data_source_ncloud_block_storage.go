package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_block_storage", dataSourceNcloudBlockStorage())
}

func dataSourceNcloudBlockStorage() *schema.Resource {
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
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudBlockStorage(), fieldMap, dataSourceNcloudBlockStorageRead)
}

func dataSourceNcloudBlockStorageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instances, err := getBlockStorageList(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudBlockStorage().Schema)
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	d.SetId(resources[0]["block_storage_no"].(string))
	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getBlockStorageList(d *schema.ResourceData, config *ProviderConfig) ([]*BlockStorage, error) {
	if config.SupportVPC {
		return getVpcBlockStorageList(d, config)
	}

	return getClassicBlockStorageList(d, config)
}

func getClassicBlockStorageList(d *schema.ResourceData, config *ProviderConfig) ([]*BlockStorage, error) {
	regionNo, err := parseRegionNoParameter(config.Client, d)
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

	logCommonRequest("getClassicBlockStorageList", reqParams)

	resp, err := config.Client.server.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		logErrorResponse("getClassicBlockStorageList", err, reqParams)
		return nil, err
	}
	logResponse("getClassicBlockStorageList", resp)

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

func getVpcBlockStorageList(d *schema.ResourceData, config *ProviderConfig) ([]*BlockStorage, error) {
	reqParams := &vserver.GetBlockStorageInstanceListRequest{
		RegionCode:       &config.RegionCode,
		ServerInstanceNo: StringPtrOrNil(d.GetOk("server_instance_no")),
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.BlockStorageInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getVpcBlockStorageList", reqParams)

	resp, err := config.Client.vserver.V2Api.GetBlockStorageInstanceList(reqParams)
	if err != nil {
		logErrorResponse("getVpcBlockStorage", err, reqParams)
		return nil, err
	}
	logResponse("getVpcBlockStorageList", resp)

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
		}

		list = append(list, instance)
	}

	return list, nil
}
