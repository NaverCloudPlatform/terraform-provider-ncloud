package ncloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceNcloudServerProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerProductsRead,

		Schema: map[string]*schema.Schema{
			"product_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to apply to the Server Product list returned.",
				Deprecated:   "use filter instead",
			},
			"exclusion_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enter a product code to exclude from the list.",
				Deprecated:  "use filter instead",
			},
			"product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enter a product code to search from the list. Use it for a single search.",
			},
			"server_image_product_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "You can get one from `data ncloud_server_images`. This is a required value, and each available server's specification varies depending on the server image product.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
				Deprecated:  "use region attribute of provider instead",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. You can get one by calling `data ncloud_zones`. default : Select the first Zone in the specific region",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"server_products": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of server product code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"filter": dataSourceFiltersSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudServerProductsRead(d *schema.ResourceData, meta interface{}) error {
	var resources []map[string]interface{}
	var err error

	if meta.(*ProviderConfig).SupportVPC == true || meta.(*ProviderConfig).Site == "fin" {
		resources, err = getVpcServerProductList(d, meta.(*ProviderConfig))
	} else {
		resources, err = getClassicServerProductList(d, meta.(*ProviderConfig))
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudServerProduct().Schema)
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

	d.SetId(dataResourceIdHash(ids))
	d.Set("server_products", ids)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), d.Get("server_products"))
	}

	return nil
}
