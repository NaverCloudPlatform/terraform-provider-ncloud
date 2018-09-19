package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudNasVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNasVolumesRead,

		Schema: map[string]*schema.Schema{
			"volume_allotment_protocol_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIncludeValues([]string{"NFS", "CIFS"}),
			},
			"is_event_configuration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_snapshot_configuration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nas_volume_instance_no_list": {
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

			"nas_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "NAS Volume Instance list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nas_volume_instance_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nas_volume_instance_status": {
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
						"is_snapshot_configuration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_event_configuration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nas_volume_instance_custom_ip_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"nas_volume_description": {
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
				},
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
		VolumeAllotmentProtocolTypeCode: ncloud.String(d.Get("volume_allotment_protocol_type_code").(string)),
		NasVolumeInstanceNoList:         ncloud.StringInterfaceList(d.Get("nas_volume_instance_no_list").([]interface{})),
		RegionNo:                        regionNo,
		ZoneNo:                          zoneNo,
	}
	if isEventConfiguration, ok := d.GetOk("is_event_configuration"); ok {
		reqParams.IsEventConfiguration = ncloud.Bool(isEventConfiguration.(bool))
	}
	if isSnapshotConfiguration, ok := d.GetOk("is_snapshot_configuration"); ok {
		reqParams.IsSnapshotConfiguration = ncloud.Bool(isSnapshotConfiguration.(bool))
	}
	resp, err := client.server.V2Api.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetNasVolumeInstanceList", err, reqParams)
		return err
	}
	logCommonResponse("GetNasVolumeInstanceList", reqParams, GetCommonResponse(resp))

	nasVolumeInstances := resp.NasVolumeInstanceList
	if len(nasVolumeInstances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}
	return nasVolumeInstancesAttributes(d, nasVolumeInstances)
}

func nasVolumeInstancesAttributes(d *schema.ResourceData, nasVolumeInstances []*server.NasVolumeInstance) error {
	var ids []string
	var s []map[string]interface{}

	for _, nasVolume := range nasVolumeInstances {
		mapping := map[string]interface{}{
			"nas_volume_instance_no":         *nasVolume.NasVolumeInstanceNo,
			"nas_volume_instance_status":     setCommonCode(nasVolume.NasVolumeInstanceStatus),
			"create_date":                    *nasVolume.CreateDate,
			"nas_volume_description":         *nasVolume.NasVolumeInstanceDescription,
			"volume_allotment_protocol_type": setCommonCode(nasVolume.VolumeAllotmentProtocolType),
			"volume_name":                    *nasVolume.VolumeName,
			"volume_total_size":              int(*nasVolume.VolumeTotalSize),
			"volume_size":                    int(*nasVolume.VolumeSize),
			"volume_use_size":                int(*nasVolume.VolumeUseSize),
			"volume_use_ratio":               *nasVolume.VolumeUseRatio,
			"snapshot_volume_size":           *nasVolume.SnapshotVolumeSize,
			"snapshot_volume_use_size":       *nasVolume.SnapshotVolumeUseSize,
			"snapshot_volume_use_ratio":      *nasVolume.SnapshotVolumeUseRatio,
			"is_snapshot_configuration":      *nasVolume.IsSnapshotConfiguration,
			"is_event_configuration":         *nasVolume.IsEventConfiguration,
			"zone":   setZone(nasVolume.Zone),
			"region": setRegion(nasVolume.Region),
		}
		if len(nasVolume.NasVolumeInstanceCustomIpList) > 0 {
			mapping["nas_volume_instance_custom_ip_list"] = customIPList(nasVolume.NasVolumeInstanceCustomIpList)
		}

		ids = append(ids, *nasVolume.NasVolumeInstanceNo)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("nas_volumes", s); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
