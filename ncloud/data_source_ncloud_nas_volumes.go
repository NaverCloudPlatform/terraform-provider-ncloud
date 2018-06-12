package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateBoolValue,
			},
			"is_snapshot_configuration": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateBoolValue,
			},
			"nas_volume_instance_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_no": {
				Type:     schema.TypeString,
				Optional: true,
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_no": {
										Type: schema.TypeString,
									},
									"zone_name": {
										Type: schema.TypeString,
									},
									"zone_description": {
										Type: schema.TypeString,
									},
								},
							},
						},
						"region": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_no": {
										Type: schema.TypeString,
									},
									"region_code": {
										Type: schema.TypeString,
									},
									"region_name": {
										Type: schema.TypeString,
									},
								},
							},
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
	conn := meta.(*NcloudSdk).conn

	reqParams := &sdk.RequestGetNasVolumeInstanceList{
		VolumeAllotmentProtocolTypeCode: d.Get("volume_allotment_protocol_type_code").(string),
		IsEventConfiguration:            d.Get("is_event_configuration").(string),
		IsSnapshotConfiguration:         d.Get("is_snapshot_configuration").(string),
		NasVolumeInstanceNoList:         StringList(d.Get("nas_volume_instance_no_list").([]interface{})),
		RegionNo:                        d.Get("region_no").(string),
		ZoneNo:                          d.Get("zone_no").(string),
	}
	resp, err := conn.GetNasVolumeInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetNasVolumeInstanceList", err, reqParams)
		return err
	}
	logCommonResponse("GetNasVolumeInstanceList", reqParams, resp.CommonResponse)

	nasVolumeInstances := resp.NasVolumeInstanceList
	if len(nasVolumeInstances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}
	return nasVolumeInstancesAttributes(d, nasVolumeInstances)
}

func nasVolumeInstancesAttributes(d *schema.ResourceData, nasVolumeInstances []sdk.NasVolumeInstance) error {
	var ids []string
	var s []map[string]interface{}

	for _, nasVolume := range nasVolumeInstances {
		mapping := map[string]interface{}{
			"nas_volume_instance_no": nasVolume.NasVolumeInstanceNo,
			"nas_volume_instance_status": map[string]interface{}{
				"code":      nasVolume.NasVolumeInstanceStatus.Code,
				"code_name": nasVolume.NasVolumeInstanceStatus.CodeName,
			},
			"create_date":            nasVolume.CreateDate,
			"nas_volume_description": nasVolume.NasVolumeInstanceDescription,
			"volume_allotment_protocol_type": map[string]interface{}{
				"code":      nasVolume.VolumeAllotmentProtocolType.Code,
				"code_name": nasVolume.VolumeAllotmentProtocolType.CodeName,
			},
			"volume_name":                        nasVolume.VolumeName,
			"volume_total_size":                  nasVolume.VolumeTotalSize,
			"volume_size":                        nasVolume.VolumeSize,
			"volume_use_size":                    nasVolume.VolumeUseSize,
			"volume_use_ratio":                   nasVolume.VolumeUseRatio,
			"snapshot_volume_size":               nasVolume.SnapshotVolumeSize,
			"snapshot_volume_use_size":           nasVolume.SnapshotVolumeUseSize,
			"snapshot_volume_use_ratio":          nasVolume.SnapshotVolumeUseRatio,
			"is_snapshot_configuration":          nasVolume.IsSnapshotConfiguration,
			"is_event_configuration":             nasVolume.IsEventConfiguration,
			"nas_volume_instance_custom_ip_list": nasVolume.NasVolumeInstanceCustomIpList,
			"zone": map[string]interface{}{
				"zone_no":          nasVolume.Zone.ZoneNo,
				"zone_name":        nasVolume.Zone.ZoneName,
				"zone_description": nasVolume.Zone.ZoneDescription,
			},
			"region": map[string]interface{}{
				"region_no":   nasVolume.Region.RegionNo,
				"region_code": nasVolume.Region.RegionCode,
				"region_name": nasVolume.Region.RegionName,
			},
		}

		ids = append(ids, nasVolume.NasVolumeInstanceNo)
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
