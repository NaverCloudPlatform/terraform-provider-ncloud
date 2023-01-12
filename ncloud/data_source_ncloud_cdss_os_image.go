package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_cdss_os_image", dataSourceNcloudCDSSOsImage())
}

func dataSourceNcloudCDSSOsImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSOsProductRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudCDSSOsProductRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_cdss_node_os_image`")
	}

	resources, err := getCDSSOsProducts(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudCDSSOsImage().Schema)
	}

	if len(resources) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	for k, v := range resources[0] {
		if k == "id" {
			d.SetId(v.(string))
		}
		d.Set(k, v)
	}

	return nil
}

func getCDSSOsProducts(config *ProviderConfig) ([]map[string]interface{}, error) {
	logCommonRequest("GetOsProductList", "")
	resp, _, err := config.Client.vcdss.V1Api.ClusterGetOsProductListGet(context.Background())

	if err != nil {
		logErrorResponse("GetOsProductList", err, "")
		return nil, err
	}

	logResponse("GetOsProductList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.ProductList {
		instance := map[string]interface{}{
			"id":         ncloud.StringValue(&r.ProductCode),
			"image_name": ncloud.StringValue(&r.ProductEnglishDesc),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
