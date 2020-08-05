package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudNasVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNasVolumeRead,

		Schema: map[string]*schema.Schema{
			"volume_allotment_protocol_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NFS", "CIFS"}, false),
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
			"no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Zone code. Get available values using the `data ncloud_zones`.",
			},

			"instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_total_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_use_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_use_ratio": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"snapshot_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_volume_use_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_volume_use_ratio": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"instance_custom_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

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

	var nasVolumeInstance *server.NasVolumeInstance
	nasVolumeInstances := resp.NasVolumeInstanceList
	if err := validateOneResult(len(nasVolumeInstances)); err != nil {
		return err
	}
	nasVolumeInstance = nasVolumeInstances[0]

	return nasVolumeInstanceAttributes(d, nasVolumeInstance)
}

func nasVolumeInstanceAttributes(d *schema.ResourceData, nasVolume *server.NasVolumeInstance) error {
	d.Set("instance_no", nasVolume.NasVolumeInstanceNo)
	d.Set("description", nasVolume.NasVolumeInstanceDescription)
	d.Set("volume_name", nasVolume.VolumeName)
	d.Set("volume_total_size", nasVolume.VolumeTotalSize)
	d.Set("volume_size", nasVolume.VolumeSize)
	d.Set("volume_use_size", nasVolume.VolumeUseSize)
	d.Set("volume_use_ratio", nasVolume.VolumeUseRatio)
	d.Set("snapshot_volume_size", nasVolume.SnapshotVolumeSize)
	d.Set("snapshot_volume_use_size", nasVolume.SnapshotVolumeUseSize)
	d.Set("snapshot_volume_use_ratio", nasVolume.SnapshotVolumeUseRatio)
	d.Set("is_snapshot_configuration", nasVolume.IsSnapshotConfiguration)
	d.Set("is_event_configuration", nasVolume.IsEventConfiguration)

	if instanceStatus := flattenCommonCode(nasVolume.NasVolumeInstanceStatus); instanceStatus["code"] != nil {
		d.Set("instance_status", instanceStatus["code"])
	}

	if typeCode := flattenCommonCode(nasVolume.VolumeAllotmentProtocolType); typeCode["code"] != nil {
		d.Set("volume_allotment_protocol_type_code", typeCode["code"])
	}

	if len(nasVolume.NasVolumeInstanceCustomIpList) > 0 {
		d.Set("instance_custom_ip_list", flattenCustomIPList(nasVolume.NasVolumeInstanceCustomIpList))
	}

	if zone := flattenZone(nasVolume.Zone); zone["zone_code"] != nil {
		d.Set("zone", zone["zone_code"])
	}

	if region := flattenRegion(nasVolume.Region); region["region_code"] != nil {
		d.Set("region", region["region_code"])
	}

	d.SetId(ncloud.StringValue(nasVolume.NasVolumeInstanceNo))

	return nil
}
