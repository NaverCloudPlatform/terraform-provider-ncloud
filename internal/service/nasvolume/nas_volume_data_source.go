package nasvolume

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func DataSourceNcloudNasVolume() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"volume_allotment_protocol_type_code": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NFS", "CIFS"}, false)),
		},
		"is_event_configuration": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"is_snapshot_configuration": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Zone code. Get available values using the `data ncloud_zones`.",
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudNasVolume(), fieldMap, dataSourceNcloudNasVolumeRead)
}

func dataSourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instances, err := getNasVolumeList(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudNasVolume().Schema)
	}

	if err := ValidateOneResult(len(resources)); err != nil {
		return err
	}

	d.SetId(resources[0]["nas_volume_no"].(string))
	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getNasVolumeList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*NasVolume, error) {
	if config.SupportVPC {
		return getVpcNasVolumeList(d, config)
	} else {
		return getClassicNasVolumeList(d, config)
	}
}

func getClassicNasVolumeList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*NasVolume, error) {
	client := config.Client

	regionNo, err := conn.ParseRegionNoParameter(d)
	if err != nil {
		return nil, err
	}

	zoneNo, err := zone.ParseZoneNoParameter(config, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.GetNasVolumeInstanceListRequest{
		VolumeAllotmentProtocolTypeCode: StringPtrOrNil(d.GetOk("volume_allotment_protocol_type_code")),
		IsEventConfiguration:            BoolPtrOrNil(d.GetOk("is_event_configuration")),
		IsSnapshotConfiguration:         BoolPtrOrNil(d.GetOk("is_snapshot_configuration")),
		RegionNo:                        regionNo,
		ZoneNo:                          zoneNo,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.NasVolumeInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getClassicNasVolumeList", reqParams)

	resp, err := client.Server.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getClassicNasVolumeList", err, reqParams)
		return nil, err
	}
	LogResponse("getClassicNasVolumeList", resp)

	var list []*NasVolume
	for _, r := range resp.NasVolumeInstanceList {
		list = append(list, convertClassicNasVolume(r))
	}

	return list, nil
}

func getVpcNasVolumeList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*NasVolume, error) {
	client := config.Client

	reqParams := &vnas.GetNasVolumeInstanceListRequest{
		RegionCode:                      &config.RegionCode,
		VolumeAllotmentProtocolTypeCode: StringPtrOrNil(d.GetOk("volume_allotment_protocol_type_code")),
		IsEventConfiguration:            BoolPtrOrNil(d.GetOk("is_event_configuration")),
		IsSnapshotConfiguration:         BoolPtrOrNil(d.GetOk("is_snapshot_configuration")),
		ZoneCode:                        StringPtrOrNil(d.GetOk("zone")),
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.NasVolumeInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getVpcNasVolumeList", reqParams)

	resp, err := client.Vnas.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcNasVolumeList", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcNasVolumeList", resp)

	var list []*NasVolume
	for _, r := range resp.NasVolumeInstanceList {
		list = append(list, convertVpcNasVolume(r))
	}

	return list, nil
}
