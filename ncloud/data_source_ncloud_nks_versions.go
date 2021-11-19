package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_nks_version", dataSourceNcloudNKSVersion())
}

func dataSourceNcloudNKSVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudVersionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"versions": {
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

func dataSourceNcloudVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_nks_version`")
	}

	resources, err := getNKSVersion(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudNKSVersion().Schema["versions"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("versions", resources); err != nil {
		return fmt.Errorf("Error setting Versions: %s", err)
	}

	return nil

}

func getNKSVersion(config *ProviderConfig) ([]map[string]interface{}, error) {

	logCommonRequest("GetNKSVersion", "")
	resp, err := config.Client.vnks.V2Api.OptionVersionGet(context.Background())

	if err != nil {
		logErrorResponse("GetNKSVersion", err, "")
		return nil, err
	}

	logResponse("GetNKSVersion", resp)

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
