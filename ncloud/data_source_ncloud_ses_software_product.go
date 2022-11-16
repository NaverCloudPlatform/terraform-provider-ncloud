package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_ses_software_product", dataSourceNcloudSESSoftwareProduct())
}

func dataSourceNcloudSESSoftwareProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudSESSoftwareProductRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSESSoftwareProductRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_ses_software_product`")
	}

	resources, err := getSESSoftwareProduct(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSESSoftwareProduct().Schema["codes"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("codes", resources); err != nil {
		return fmt.Errorf("Error setting Codes: %s", err)
	}

	return nil

}

func getSESSoftwareProduct(config *ProviderConfig) ([]map[string]interface{}, error) {

	logCommonRequest("GetSESSoftwareProduct", "")
	resp, _, err := config.Client.vses.V2Api.GetOsProductListUsingGET(context.Background())

	if err != nil {
		logErrorResponse("GetSESSoftwareProduct", err, "")
		return nil, err
	}

	logResponse("GetSESSoftwareProduct", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.ProductList {
		instance := map[string]interface{}{
			"value": ncloud.StringValue(&r.ProductCode),
			"label": ncloud.StringValue(&r.ProductEnglishDesc),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
