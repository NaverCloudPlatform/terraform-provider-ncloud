package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func init() {
	RegisterDataSource("ncloud_cdss_kafka_versions", dataSourceNcloudCDSSKafkaVersions())
}

func dataSourceNcloudCDSSKafkaVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSKafkaVersionsRead,
		Schema: map[string]*schema.Schema{
			"kafka_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(dataSourceNcloudCDSSKafkaVersion()),
			},
		},
	}
}

func dataSourceNcloudCDSSKafkaVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_cdss_kafka_versions`")
	}

	resources, err := getCDSSKafkaVersions(config)
	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("kafka_versions", resources); err != nil {
		return fmt.Errorf("Error setting node products: %s", err)
	}

	return nil
}
