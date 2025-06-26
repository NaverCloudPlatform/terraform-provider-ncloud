package zone

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"filter": DataSourceFiltersSchema(),
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

	zones, err = getVpcZones(meta.(*conn.ProviderConfig))
	if err != nil {
		return err
	}

	resources := flattenZones(zones)

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudZones().Schema["zones"].Elem.(*schema.Resource).Schema)
	}

	if err := d.Set("zones", resources); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return WriteToFile(output.(string), d.Get("zones"))
	}

	return nil
}

func getVpcZones(config *conn.ProviderConfig) ([]*Zone, error) {
	client := config.Client
	regionCode := config.RegionCode

	resp, err := client.Vserver.V2Api.GetZoneList(&vserver.GetZoneListRequest{RegionCode: &regionCode})
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
