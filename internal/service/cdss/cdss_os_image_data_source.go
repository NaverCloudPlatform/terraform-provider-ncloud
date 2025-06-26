package cdss

import (
	"context"
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudCDSSOsImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSOsProductRead,
		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
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
	config := meta.(*conn.ProviderConfig)

	resources, err := getCDSSOsProducts(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudCDSSOsImage().Schema)
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

func getCDSSOsProducts(config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	LogCommonRequest("GetOsProductList", "")
	resp, _, err := config.Client.Vcdss.V1Api.ClusterGetOsProductListGet(context.Background())

	if err != nil {
		LogErrorResponse("GetOsProductList", err, "")
		return nil, err
	}

	LogResponse("GetOsProductList", resp)

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
