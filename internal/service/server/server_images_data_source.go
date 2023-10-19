package server

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServerImagesRead,

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"platform_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"filter": DataSourceFiltersSchema(),

			"server_images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(DataSourceNcloudServerImage()),
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Deprecated
			"product_name_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsValidRegExp),
				Deprecated:       "use `filter` instead",
			},
			"exclusion_product_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field no longer support",
			},
			"block_storage_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{50, 100})),
				Deprecated:       "use `filter` instead",
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
	resources, err := getServerImageProductListFiltered(d, meta.(*conn.ProviderConfig))

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

	d.SetId(DataResourceIdHash(ids))
	d.Set("ids", ids)
	d.Set("server_images", resources)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return WriteToFile(output.(string), d.Get("server_images"))
	}

	return nil
}
