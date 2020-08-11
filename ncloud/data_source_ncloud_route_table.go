package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRouteTableRead,

		Schema: map[string]*schema.Schema{
			"route_table_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_subnet_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getRouteTableInstance(config, d.Get("route_table_no").(string))

	if err != nil {
		return err
	}

	d.SetId(*instance.RouteTableNo)
	d.Set("route_table_no", instance.RouteTableNo)
	d.Set("name", instance.RouteTableName)
	d.Set("description", instance.RouteTableDescription)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("supported_subnet_type", instance.SupportedSubnetType.Code)
	d.Set("is_default", instance.IsDefault)
	d.Set("status", instance.RouteTableStatus.Code)

	return nil
}
