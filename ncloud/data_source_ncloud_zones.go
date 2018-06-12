package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudZonesRead,

		Schema: map[string]*schema.Schema{
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1",
			},
			"zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_no": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceNcloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	log.Printf("[DEBUG] Get Zone List")
	d.SetId(time.Now().UTC().String())

	regionNo, _ := d.GetOk("region_no")

	resp, err := conn.GetZoneList(regionNo.(string))
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("no matching zones found")
	}

	var zones []common.Zone

	for _, zone := range resp.Zone {
		zones = append(zones, zone)
	}

	if len(zones) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return zonesAttributes(d, zones)
}

func zonesAttributes(d *schema.ResourceData, zones []common.Zone) error {

	var ids []string
	var s []map[string]interface{}
	for _, zone := range zones {
		mapping := map[string]interface{}{
			"zone_no":          zone.ZoneNo,
			"zone_name":        zone.ZoneName,
			"zone_description": zone.ZoneDescription,
		}

		log.Printf("[DEBUG] ncloud_regions - adding region mapping: %v", mapping)
		ids = append(ids, string(zone.ZoneNo))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
