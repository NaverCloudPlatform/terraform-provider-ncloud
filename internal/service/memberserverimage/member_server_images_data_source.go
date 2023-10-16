package memberserverimage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudMemberServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudMemberServerImagesRead,

		Schema: map[string]*schema.Schema{
			"no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of member server images to view",
			},
			"platform_type_code_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of platform codes of server images to view",
			},
			"filter": DataSourceFiltersSchema(),
			"name_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsValidRegExp),
				Deprecated:       "use filter instead",
				Description:      "A regex string to apply to the member server image list returned by ncloud",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "use region attribute of provider instead",
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},

			"member_server_images": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of Member server image no",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of file that can save data source after running `terraform plan`.",
			},
		},
	}
}

func dataSourceNcloudMemberServerImagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcMemberServerImage(d, config)
	} else {
		resources, err = getClassicMemberServerImage(d, config)
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudMemberServerImage().Schema)
	}

	if len(resources) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return memberServerImagesAttributes(d, resources)
}

func memberServerImagesAttributes(d *schema.ResourceData, memberServerImages []map[string]interface{}) error {
	var ids []string

	for _, r := range memberServerImages {
		ids = append(ids, r["id"].(string))
	}

	d.SetId(DataResourceIdHash(ids))
	_ = d.Set("member_server_images", ids)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return WriteToFile(output.(string), d.Get("member_server_images"))
	}

	return nil
}
