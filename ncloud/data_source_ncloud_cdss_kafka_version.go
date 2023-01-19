package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_cdss_kafka_version", dataSourceNcloudCDSSKafkaVersion())
}

func dataSourceNcloudCDSSKafkaVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSKafkaVersionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudCDSSKafkaVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_cdss_kafka_version`")
	}

	resources, err := getCDSSKafkaVersions(config)
	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudCDSSKafkaVersion().Schema)
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

func getCDSSKafkaVersions(config *ProviderConfig) ([]map[string]interface{}, error) {
	logCommonRequest("GetCDSSVersionList", "")
	resp, _, err := config.Client.vcdss.V1Api.ClusterGetCDSSVersionListGet(context.Background())

	if err != nil {
		logErrorResponse("GetCDSSVersionList", err, "")
		return nil, err
	}

	logResponse("GetCDSSVersionList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.KafkaVersionList {
		instance := map[string]interface{}{
			"id":   ncloud.StringValue(&r.KafkaVersionCode),
			"name": ncloud.StringValue(&r.KafkaVersionName),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
