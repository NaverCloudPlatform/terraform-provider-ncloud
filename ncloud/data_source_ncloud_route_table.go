package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudRouteTable() *schema.Resource {
	fieldMap := make(map[string]*schema.Schema)
	fieldMap["route_table_no"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
	fieldMap["vpc_no"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
	fieldMap["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
	fieldMap["supported_subnet_type"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
	fieldMap["filter"] = dataSourceFiltersSchema()

	return GetSingularDataSourceItemSchema(resourceNcloudRouteTable(), fieldMap, dataSourceNcloudRouteTableRead)
}

func dataSourceNcloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	resources, err := getRouteTableListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	for k, v := range resources[0] {
		if k == "id" {
			d.SetId(v.(string))
			continue
		}
		d.Set(k, v)
	}

	return nil
}
