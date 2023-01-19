package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_ses_versions", dataSourceNcloudSESVersions())
}

func dataSourceNcloudSESVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudSESVersionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"versions": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudSESVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_ses_versions`")
	}

	resources, err := getSESVersion(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudSESVersions().Schema["versions"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("versions", resources); err != nil {
		return fmt.Errorf("Error setting Versions: %s", err)
	}

	return nil

}

func getSESVersion(config *ProviderConfig) ([]map[string]interface{}, error) {

	logCommonRequest("GetSESVersion", "")
	resp, _, err := config.Client.vses.V2Api.GetSearchEngineVersionListUsingGET(context.Background())

	if err != nil {
		logErrorResponse("GetSESVersion", err, "")
		return nil, err
	}

	logResponse("GetSESVersion", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.SearchEngineVersionList {
		instance := map[string]interface{}{
			"id":      ncloud.StringValue(&r.SearchEngineVersionCode),
			"name":    ncloud.StringValue(&r.SearchEngineVersionName),
			"type":    ncloud.StringValue(&r.Type_),
			"version": ncloud.StringValue(&strings.Split(r.SearchEngineVersionName, " ")[1]),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
