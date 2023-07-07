package region

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var regionSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"region_no": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}
