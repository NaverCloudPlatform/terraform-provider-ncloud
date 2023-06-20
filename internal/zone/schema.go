package zone

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var zoneSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"zone_no": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"zone_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"zone_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"zone_description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region_no": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}
