package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_nas_volume", dataSourceNcloudNasVolume())
}

func dataSourceNcloudNasVolume() *schema.Resource {
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
			ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"NFS", "CIFS"}, false)),
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
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudNasVolume(), fieldMap, dataSourceNcloudNasVolumeRead)
}

func dataSourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instances, err := getNasVolumeList(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudNasVolume().Schema)
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	d.SetId(resources[0]["nas_volume_no"].(string))
	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getNasVolumeList(d *schema.ResourceData, config *ProviderConfig) ([]*NasVolume, error) {
	if config.SupportVPC {
		return getVpcNasVolumeList(d, config)
	} else {
		return getClassicNasVolumeList(d, config)
	}
}

func getClassicNasVolumeList(d *schema.ResourceData, config *ProviderConfig) ([]*NasVolume, error) {
	client := config.Client

	regionNo, err := parseRegionNoParameter(d)
	if err != nil {
		return nil, err
	}

	zoneNo, err := parseZoneNoParameter(config, d)
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

	logCommonRequest("getClassicNasVolumeList", reqParams)

	resp, err := client.server.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("getClassicNasVolumeList", err, reqParams)
		return nil, err
	}
	logResponse("getClassicNasVolumeList", resp)

	var list []*NasVolume
	for _, r := range resp.NasVolumeInstanceList {
		list = append(list, convertClassicNasVolume(r))
	}

	return list, nil
}

func getVpcNasVolumeList(d *schema.ResourceData, config *ProviderConfig) ([]*NasVolume, error) {
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

	logCommonRequest("getVpcNasVolumeList", reqParams)

	resp, err := client.vnas.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("getVpcNasVolumeList", err, reqParams)
		return nil, err
	}
	logResponse("getVpcNasVolumeList", resp)

	var list []*NasVolume
	for _, r := range resp.NasVolumeInstanceList {
		list = append(list, convertVpcNasVolume(r))
	}

	return list, nil
}
