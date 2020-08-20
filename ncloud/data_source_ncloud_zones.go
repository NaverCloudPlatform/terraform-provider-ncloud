package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"filter": dataSourceFiltersSchema(),
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
	d.SetId(time.Now().UTC().String())

	var zones []*Zone
	var err error

	if meta.(*ProviderConfig).SupportVPC == true || meta.(*ProviderConfig).Site == "fin" {
		zones, err = getVpcZones(d, meta.(*ProviderConfig))
	} else {
		zones, err = getClassicZones(d, meta.(*ProviderConfig))
	}

	if err != nil {
		return err
	}

	resources := flattenZones(zones)

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudZones().Schema["zones"].Elem.(*schema.Resource).Schema)
	}

	if err := d.Set("zones", resources); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), d.Get("zones"))
	}

	return nil
}

func getClassicZones(d *schema.ResourceData, config *ProviderConfig) ([]*Zone, error) {
	client := config.Client
	regionNo := config.RegionNo

	resp, err := client.server.V2Api.GetZoneList(&server.GetZoneListRequest{RegionNo: &regionNo})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching zones found")
	}

	var zones []*Zone

	for _, zone := range resp.ZoneList {
		zones = append(zones, GetZone(zone))
	}

	return zones, nil
}

func getVpcZones(d *schema.ResourceData, config *ProviderConfig) ([]*Zone, error) {
	client := config.Client
	regionCode := config.RegionCode

	resp, err := client.vserver.V2Api.GetZoneList(&vserver.GetZoneListRequest{RegionCode: &regionCode})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching zones found")
	}

	var zones []*Zone

	for _, zone := range resp.ZoneList {
		zones = append(zones, GetZone(zone))
	}

	return zones, nil
}
