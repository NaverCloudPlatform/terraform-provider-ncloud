package ncloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRouteTablesRead,

		Schema: map[string]*schema.Schema{
			"vpc_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"supported_subnet_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": dataSourceFiltersSchema(),
			"route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(resourceNcloudRouteTable()),
			},
		},
	}
}

func dataSourceNcloudRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	resources, err := getRouteTableListFiltered(d, config)

	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("route_tables", resources); err != nil {
		return fmt.Errorf("Error setting route table ids: %s", err)
	}

	return nil
}
