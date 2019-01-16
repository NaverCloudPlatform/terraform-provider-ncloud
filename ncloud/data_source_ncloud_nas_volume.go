package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudNasVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNasVolumeRead,

		Schema: map[string]*schema.Schema{
			"volume_allotment_protocol_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIncludeValues([]string{"NFS", "CIFS"}),
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
			"region_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region code. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_no"},
			},
			"region_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region number. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_code"},
			},
			"zone_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone code",
				ConflictsWith: []string{"zone_no"},
			},
			"zone_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone number",
				ConflictsWith: []string{"zone_code"},
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
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_allotment_protocol_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
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
			"zone": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     zoneSchemaResource,
			},
			"region": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     regionSchemaResource,
			},
		},
	}
}

func dataSourceNcloudNasVolumeRead(d *schema.ResourceData, meta interface{}) error {
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
		VolumeAllotmentProtocolTypeCode: ncloud.String(d.Get("volume_allotment_protocol_type_code").(string)),
		NasVolumeInstanceNoList:         expandStringInterfaceList(d.Get("no_list").([]interface{})),
		RegionNo:                        regionNo,
		ZoneNo:                          zoneNo,
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
	if len(nasVolumeInstances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}
	nasVolumeInstance = nasVolumeInstances[0]

	return nasVolumeInstanceAttributes(d, nasVolumeInstance)
}

func nasVolumeInstanceAttributes(d *schema.ResourceData, nasVolume *server.NasVolumeInstance) error {
	d.Set("instance_no", nasVolume.NasVolumeInstanceNo)
	d.Set("create_date", nasVolume.CreateDate)
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

	if err := d.Set("instance_status", flattenCommonCode(nasVolume.NasVolumeInstanceStatus)); err != nil {
		return err
	}

	if err := d.Set("volume_allotment_protocol_type", flattenCommonCode(nasVolume.VolumeAllotmentProtocolType)); err != nil {
		return err
	}

	if len(nasVolume.NasVolumeInstanceCustomIpList) > 0 {
		d.Set("instance_custom_ip_list", flattenCustomIPList(nasVolume.NasVolumeInstanceCustomIpList))
	}

	if err := d.Set("zone", flattenZone(nasVolume.Zone)); err != nil {
		return err
	}

	if err := d.Set("region", flattenRegion(nasVolume.Region)); err != nil {
		return err
	}

	d.SetId(ncloud.StringValue(nasVolume.NasVolumeInstanceNo))

	return nil
}
