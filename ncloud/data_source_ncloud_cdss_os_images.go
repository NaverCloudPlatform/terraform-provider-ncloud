package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func init() {
	RegisterDataSource("ncloud_cdss_os_images", dataSourceNcloudCDSSOsImages())
}

func dataSourceNcloudCDSSOsImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSOsImagesRead,
		Schema: map[string]*schema.Schema{
			"os_images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(dataSourceNcloudCDSSOsImage()),
			},
		},
	}
}

func dataSourceNcloudCDSSOsImagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
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
