package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/hashicorp/terraform/helper/schema"
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
	conn := meta.(*NcloudSdk).conn

	d.SetId(time.Now().UTC().String())

	resp, err := conn.GetZoneList(parseRegionNoParameter(conn, d))
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
		mapping := setZone(zone)
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
