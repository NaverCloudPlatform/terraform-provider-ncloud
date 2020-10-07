package ncloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_subnets", dataSourceNcloudSubnets())
}

func dataSourceNcloudSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudSubnetsRead,

		Schema: map[string]*schema.Schema{
			"subnet_no": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of subnet ID to retrieve",
			},
			"vpc_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC ID that you want to filter from",
			},
			"subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CIDR block for the subnet.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Available Zone. Get available values using the `data ncloud_zones`.",
			},
			"network_acl_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network ACL No. Get available values using the `default_network_acl_no` from Resource `ncloud_vpc` or Data source `data.ncloud_network_acls`.",
			},
			"subnet_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false),
				Description:  "Internet Gateway Only. PUBLC(Yes/Public), PRIVATE(No/Private).",
			},
			"usage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"GEN", "LOADB", "BM"}, false),
				Description:  "Usage type. GEN(Normal), LOADB(Load Balance), BM(BareMetal). default : GEN(Normal).",
			},
			"filter": dataSourceFiltersSchema(),
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(resourceNcloudSubnet()),
			},
		},
	}
}

func dataSourceNcloudSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_subnets`")
	}

	resources, err := getSubnetListFiltered(d, config)

	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("subnets", resources); err != nil {
		return fmt.Errorf("Error setting Subnets: %s", err)
	}

	return nil
}
