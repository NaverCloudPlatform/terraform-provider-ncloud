package ncloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func init() {
	RegisterDatasource("ncloud_network_interfaces", dataSourceNcloudNetworkInterfaces())
}

func dataSourceNcloudNetworkInterfaces() *schema.Resource {
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
			"filter": dataSourceFiltersSchema(),

			"network_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(resourceNcloudNetworkInterface()),
			},
		},
	}
}

func dataSourceNcloudNetworkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcNetworkInterfaceListFiltered(d, config)
	} else {
		return NotSupportClassic("Network Interface")
	}

	if err != nil {
		return err
	}

	if resources == nil || len(resources) == 0 {
		return errors.New("no matching Network Interfaces found")
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("network_interfaces", resources); err != nil {
		return fmt.Errorf("error setting Network Interfaces: %s", err)
	}

	return nil
}
