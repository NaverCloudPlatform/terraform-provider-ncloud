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

func DataSourceNcloudNKSVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudVersionsRead,

		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"hypervisor_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
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

func dataSourceNcloudVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_nks_versions`")
	}

	resources, err := getNKSVersion(config, d)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudNKSVersions().Schema["versions"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("versions", resources); err != nil {
		return fmt.Errorf("Error setting Versions: %s", err)
	}

	return nil

}

func getNKSVersion(config *conn.ProviderConfig, d *schema.ResourceData) ([]map[string]interface{}, error) {

	LogCommonRequest("GetNKSVersion", "")
	hypervisorCode := StringPtrOrNil(d.GetOk("hypervisor_code"))

	opt := make(map[string]interface{})
	if hypervisorCode != nil {
		opt["hypervisorCode"] = hypervisorCode
	}

	resp, err := config.Client.Vnks.V2Api.OptionVersionGet(context.Background(), opt)

	if err != nil {
		LogErrorResponse("GetNKSVersion", err, "")
		return nil, err
	}

	LogResponse("GetNKSVersion", resp)

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
