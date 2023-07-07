package server

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudServersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": DataSourceFiltersSchema(),
		},
	}
}

func dataSourceNcloudServersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instances, err := getServerList(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. there is no available server resource")
	}

	if values, ok := d.GetOk("ids"); ok {
		return readServersIDs(d, values.(*schema.Set).List(), instances)
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudServer().Schema)
	}

	if len(resources) == 0 {
		return fmt.Errorf("no results with filter. there is no available server resource")
	}

	var ids []string
	for _, r := range resources {
		for k, v := range r {
			if k == "instance_no" {
				ids = append(ids, v.(string))
			}
		}
	}

	d.SetId(DataResourceIdHash(ids))
	d.Set("ids", ids)
	return nil
}

func readServersIDs(d *schema.ResourceData, values []interface{}, serverInstances []*ServerInstance) error {
	var ids []string
	for _, id := range values {
		for _, s := range serverInstances {
			if *s.ServerInstanceNo == id.(string) {
				ids = append(ids, id.(string))
				break
			}
		}
	}

	if len(values) != len(ids) {
		return fmt.Errorf("invalid server id specified")
	}

	d.SetId(DataResourceIdHash(ids))
	d.Set("ids", ids)
	return nil
}
