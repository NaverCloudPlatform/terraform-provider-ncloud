package ncloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDatasource("ncloud_server_images", dataSourceNcloudServerImages())
}

func dataSourceNcloudServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImagesRead,

		Schema: map[string]*schema.Schema{
			"product_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to apply to the server image list returned by ncloud.",
				Deprecated:   "use filter instead",
			},
			"exclusion_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Product code you want to exclude from the list.",
			},
			"product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Product code you want to view on the list. Use this when searching for 1 product.",
			},
			"platform_type_code_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Values required for identifying platforms in list-type.",
			},
			"block_storage_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{50, 100}),
				Description:  "Block storage size.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
				Deprecated:  "use region attribute of provider instead",
			},
			"infra_resource_detail_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "infra resource detail type code.",
			},
			"server_images": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of server image product code.",
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

func dataSourceNcloudServerImagesRead(d *schema.ResourceData, meta interface{}) error {
	var resources []map[string]interface{}
	var err error

	if meta.(*ProviderConfig).SupportVPC == true {
		resources, err = getVpcServerImageProductList(d, meta.(*ProviderConfig))
	} else {
		resources, err = getClassicServerImageProductList(d, meta.(*ProviderConfig))
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudServerImage().Schema)
	}

	if len(resources) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return serverImagesAttributes(d, resources)
}

func serverImagesAttributes(d *schema.ResourceData, serverImages []map[string]interface{}) error {
	var ids []string

	for _, r := range serverImages {
		for k, v := range r {
			if k == "id" {
				ids = append(ids, v.(string))
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("server_images", ids)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), d.Get("server_images"))
	}

	return nil
}
