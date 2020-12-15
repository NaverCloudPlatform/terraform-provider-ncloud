package ncloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_server_images", dataSourceNcloudServerImages())
}

func dataSourceNcloudServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImagesRead,

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"platform_type_code_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"infra_resource_detail_type_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": dataSourceFiltersSchema(),

			"server_images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(dataSourceNcloudServerImage()),
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Deprecated
			"product_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Deprecated:   "use `filter` instead",
			},
			"exclusion_product_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field no longer support",
			},
			"block_storage_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{50, 100}),
				Deprecated:   "use `filter` instead",
			},
			"region": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use provider config instead",
			},
		},
	}
}

func dataSourceNcloudServerImagesRead(d *schema.ResourceData, meta interface{}) error {
	resources, err := getServerImageProductListFiltered(d, meta.(*ProviderConfig))

	if err != nil {
		return err
	}

	if len(resources) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return serverImagesAttributes(d, resources)
}

func serverImagesAttributes(d *schema.ResourceData, resources []map[string]interface{}) error {
	var ids []string

	for _, r := range resources {
		for k, v := range r {
			if k == "id" {
				ids = append(ids, v.(string))
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)
	d.Set("server_images", resources)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), d.Get("server_images"))
	}

	return nil
}
