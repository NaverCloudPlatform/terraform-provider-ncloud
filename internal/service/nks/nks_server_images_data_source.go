package nks

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudNKSServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNKSServerImagesRead,

		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"hypervisor_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_nks_node_pool_server_images`")
	}

	resources, err := getNKSServerImages(config, d)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudNKSServerImages().Schema["images"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("images", resources); err != nil {
		return fmt.Errorf("Error setting Codes: %s", err)
	}

	return nil

}

func getNKSServerImages(config *conn.ProviderConfig, d *schema.ResourceData) ([]map[string]interface{}, error) {

	LogCommonRequest("GetNKSServerImages", "")
	hypervisorCode := StringPtrOrNil(d.GetOk("hypervisor_code"))

	opt := make(map[string]interface{})
	if hypervisorCode != nil {
		opt["hypervisorCode"] = hypervisorCode
	}

	resp, err := config.Client.Vnks.V2Api.OptionServerImageGet(context.Background(), opt)

	if err != nil {
		LogErrorResponse("GetNKSServerImages", err, "")
		return nil, err
	}

	LogResponse("GetNKSServerImages", resp)

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
