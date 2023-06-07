package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_nks_server_images", dataSourceNcloudNKSServerImages())
}

func dataSourceNcloudNKSServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNKSServerImagesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"images": {
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

func dataSourceNcloudNKSServerImagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_nks_node_pool_server_images`")
	}

	resources, err := getNKSServerImages(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudNKSServerImages().Schema["images"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("images", resources); err != nil {
		return fmt.Errorf("Error setting Codes: %s", err)
	}

	return nil

}

func getNKSServerImages(config *ProviderConfig) ([]map[string]interface{}, error) {

	logCommonRequest("GetNKSServerImages", "")
	resp, err := config.Client.vnks.V2Api.OptionServerImageGet(context.Background())

	if err != nil {
		logErrorResponse("GetNKSServerImages", err, "")
		return nil, err
	}

	logResponse("GetNKSServerImages", resp)

	resources := []map[string]interface{}{}

	for _, r := range *resp {
		instance := map[string]interface{}{
			"value": ncloud.StringValue(r.Value),
			"label": ncloud.StringValue(r.Label),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
