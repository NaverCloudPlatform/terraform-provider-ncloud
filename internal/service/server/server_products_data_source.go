package server

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudServerProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerProductsRead,

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_image_product_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(DataSourceNcloudServerProduct()),
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": DataSourceFiltersSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Deprecated
			"exclusion_product_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use filter instead",
			},
		},
	}
}

func dataSourceNcloudServerProductsRead(d *schema.ResourceData, meta interface{}) error {
	var resources []map[string]interface{}
	var err error

	resources, err = getVpcServerProductList(d, meta.(*conn.ProviderConfig))
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudServerProduct().Schema)
	}

	if len(resources) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return serverProductsAttributes(d, resources)
}

func serverProductsAttributes(d *schema.ResourceData, serverProduct []map[string]interface{}) error {
	var ids []string

	for _, r := range serverProduct {
		for k, v := range r {
			if k == "id" {
				ids = append(ids, v.(string))
			}
		}
	}

	d.SetId(DataResourceIdHash(ids))
	d.Set("ids", ids)
	d.Set("server_products", serverProduct)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return WriteToFile(output.(string), d.Get("server_products"))
	}

	return nil
}
