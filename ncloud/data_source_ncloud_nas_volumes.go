package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceNcloudNasVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNasVolumesRead,

		Schema: map[string]*schema.Schema{
			"volume_allotment_protocol_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NFS", "CIFS"}, false),
			},
			"is_event_configuration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_snapshot_configuration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone code. Get available values using the `data ncloud_zones`.",
			},
			"nas_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of NAS Volume Instance no",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudNasVolumesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := &server.GetNasVolumeInstanceListRequest{
		RegionNo: regionNo,
		ZoneNo:   zoneNo,
	}

	if volumeAllotmentProtocolTypeCode, ok := d.GetOk("volume_allotment_protocol_type_code"); ok {
		reqParams.VolumeAllotmentProtocolTypeCode = ncloud.String(volumeAllotmentProtocolTypeCode.(string))
	}

	if noList, ok := d.GetOk("no_list"); ok {
		reqParams.NasVolumeInstanceNoList = expandStringInterfaceList(noList.([]interface{}))
	}

	if isEventConfiguration, ok := d.GetOk("is_event_configuration"); ok {
		reqParams.IsEventConfiguration = ncloud.Bool(isEventConfiguration.(bool))
	}

	if isSnapshotConfiguration, ok := d.GetOk("is_snapshot_configuration"); ok {
		reqParams.IsSnapshotConfiguration = ncloud.Bool(isSnapshotConfiguration.(bool))
	}

	logCommonRequest("GetNasVolumeInstanceList", reqParams)

	resp, err := client.server.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetNasVolumeInstanceList", err, reqParams)
		return err
	}
	logCommonResponse("GetNasVolumeInstanceList", GetCommonResponse(resp))

	nasVolumeInstances := resp.NasVolumeInstanceList
	if len(nasVolumeInstances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}
	return nasVolumeInstancesAttributes(d, nasVolumeInstances)
}

func nasVolumeInstancesAttributes(d *schema.ResourceData, nasVolumeInstances []*server.NasVolumeInstance) error {
	var ids []string

	for _, nasVolume := range nasVolumeInstances {
		ids = append(ids, ncloud.StringValue(nasVolume.NasVolumeInstanceNo))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("nas_volumes", flattenNasVolumeInstances(nasVolumeInstances)); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), d.Get("nas_volumes"))
	}

	return nil
}
