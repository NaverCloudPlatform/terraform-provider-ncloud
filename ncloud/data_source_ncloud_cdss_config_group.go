package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_cdss_config_group", dataSourceNcloudCDSSConfigGroup())
}

func dataSourceNcloudCDSSConfigGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSConfigGroupRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kafka_version_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudCDSSConfigGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("dataSource `ncloud_cdss_config_group`")
	}

	resources, err := getCDSSConfigGroups(config, *StringPtrOrNil(d.GetOk("kafka_version_code")))
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

func getCDSSConfigGroups(config *ProviderConfig, kafkaVersionCode string) ([]map[string]interface{}, error) {
	logCommonRequest("GetCDSSConfigGroups", "")
	resp, _, err := config.Client.vcdss.V1Api.ConfigGroupGetKafkaVersionConfigGroupListPost(context.Background(), vcdss.GetKafkaVersionConfigGroupListRequest{
		KafkaVersionCode: kafkaVersionCode,
	})

	if err != nil {
		logErrorResponse("GetCDSSConfigGroups", err, "")
		return nil, err
	}

	logResponse("GetCDSSConfigGroups", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.Result.KafkaConfigGroupList {
		instance := map[string]interface{}{
			"id":                 *ncloud.Int32String(r.ConfigGroupNo),
			"config_group_no":    *ncloud.Int32String(r.ConfigGroupNo),
			"name":               ncloud.StringValue(&r.ConfigGroupName),
			"description":        ncloud.StringValue(&r.Description),
			"kafka_version_code": ncloud.StringValue(&r.KafkaVersionCode),
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
