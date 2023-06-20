package ses

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func init() {
	RegisterDataSource("ncloud_ses_node_os_images", dataSourceNcloudSESNodeOsImage())
}

func dataSourceNcloudSESNodeOsImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudSESNodeOsImageRead,

		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSESNodeOsImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_ses_node_os_image`")
	}

	resources, err := getSESNodeOsImage(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSESNodeOsImage().Schema["images"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("images", resources); err != nil {
		return fmt.Errorf("Error setting Codes: %s", err)
	}

	return nil

}

func getSESNodeOsImage(config *ProviderConfig) ([]map[string]interface{}, error) {

	LogCommonRequest("GetSESNodeOsImage", "")
	resp, _, err := config.Client.Vses.V2Api.GetOsProductListUsingGET(context.Background())

	if err != nil {
		LogErrorResponse("GetSESNodeOsImage", err, "")
		return nil, err
	}

	LogResponse("GetSESNodeOsImage", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.ProductList {
		instance := map[string]interface{}{
			"id":   ncloud.StringValue(&r.ProductCode),
			"name": ncloud.StringValue(&r.ProductEnglishDesc),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
