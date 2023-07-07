package server

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
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
			// Deprecated
			"internet_line_type_code": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
				Deprecated:       "This parameter is no longer used.",
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
			"product_name_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringIsValidRegExp),
				Deprecated:       "use filter instead",
			},
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

	if meta.(*conn.ProviderConfig).SupportVPC == true {
		resources, err = getVpcServerProductList(d, meta.(*conn.ProviderConfig))
	} else {
		resources, err = getClassicServerProductList(d, meta.(*conn.ProviderConfig))
	}

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
