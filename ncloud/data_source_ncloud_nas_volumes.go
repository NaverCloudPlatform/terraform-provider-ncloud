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

			"nas_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "NAS Volume Instance list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"is_snapshot_configuration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_event_configuration": {
							Type:     schema.TypeBool,
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
		writeToFile(output.(string), d.Get("nas_volumes"))
	}

	return nil
}
