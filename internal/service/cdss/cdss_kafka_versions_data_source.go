package cdss

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudCDSSKafkaVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSKafkaVersionsRead,
		Schema: map[string]*schema.Schema{
			"kafka_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(DataSourceNcloudCDSSKafkaVersion()),
			},
		},
	}
}

func dataSourceNcloudCDSSKafkaVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
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
