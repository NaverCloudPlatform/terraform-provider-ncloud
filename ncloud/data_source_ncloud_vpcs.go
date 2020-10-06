package ncloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_vpcs", dataSourceNcloudVpcs())
}

func dataSourceNcloudVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudVpcsRead,

		Schema: map[string]*schema.Schema{
			"vpc_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"INIT", "CREATING", "RUN", "TERMTING"}, false),
			},
			"filter": dataSourceFiltersSchema(),
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(resourceNcloudVpc()),
			},
		},
	}
}

func dataSourceNcloudVpcsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_vpcs`")
	}

	resources, err := getVpcListFiltered(d, config)

	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("vpcs", resources); err != nil {
		return fmt.Errorf("Error setting vpcs: %s", err)
	}

	return nil
}
