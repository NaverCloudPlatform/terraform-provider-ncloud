package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNetworkInterfacesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter": DataSourceFiltersSchema(),

			"network_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(ResourceNcloudNetworkInterface()),
			},
		},
	}
}

func dataSourceNcloudNetworkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcNetworkInterfaceListFiltered(d, config)
	} else {
		return NotSupportClassic("data source `ncloud_network_interface`")
	}

	if err != nil {
		return err
	}

	if len(resources) == 0 {
		return errors.New("no matching Network Interfaces found")
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("network_interfaces", resources); err != nil {
		return fmt.Errorf("error setting Network Interfaces: %s", err)
	}

	return nil
}
