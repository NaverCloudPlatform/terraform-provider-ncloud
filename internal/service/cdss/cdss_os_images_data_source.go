package cdss

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudCDSSOsImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSOsImagesRead,
		Schema: map[string]*schema.Schema{
			"os_images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(DataSourceNcloudCDSSOsImage()),
			},
		},
	}
}

func dataSourceNcloudCDSSOsImagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_cdss_os_images`")
	}

	resources, err := getCDSSOsProducts(config)
	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("os_images", resources); err != nil {
		return fmt.Errorf("Error setting os images: %s", err)
	}

	return nil
}
