package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vcdss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func init() {
	RegisterDataSource("ncloud_cdss_node_products", dataSourceNcloudCDSSNodeProducts())
}

func dataSourceNcloudCDSSNodeProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudCDSSNodeProductsRead,
		Schema: map[string]*schema.Schema{
			"os_image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudCDSSNodeProductsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("datasource `ncloud_cdss_node_products`")
	}

	reqParams := vcdss.NodeProduct{
		SoftwareProductCode: *StringPtrOrNil(d.GetOk("os_image")),
		SubnetNo:            *getInt32FromString(d.GetOk("subnet_no")),
	}

	resources, err := getCDSSNodeProducts(config, reqParams)
	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("node_products", resources); err != nil {
		return fmt.Errorf("Error setting node products: %s", err)
	}

	return nil
}
