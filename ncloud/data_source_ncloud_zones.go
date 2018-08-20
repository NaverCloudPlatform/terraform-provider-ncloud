package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceNcloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudZonesRead,

		Schema: map[string]*schema.Schema{
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
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     zoneSchemaResource,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	d.SetId(time.Now().UTC().String())

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	resp, err := client.server.V2Api.GetZoneList(&server.GetZoneListRequest{RegionNo: regionNo})
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("no matching zones found")
	}

	var zones []*Zone

	for _, zone := range resp.ZoneList {
		zones = append(zones, GetZone(zone))
	}

	if len(zones) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return zonesAttributes(d, zones)
}

func zonesAttributes(d *schema.ResourceData, zones []*Zone) error {

	var ids []string
	var s []map[string]interface{}
	for _, zone := range zones {
		mapping := setZone(zone)
		ids = append(ids, *zone.ZoneNo)
		s = append(s, mapping)
	}
	log.Printf("[DEBUG] zones: %#v", s)
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
