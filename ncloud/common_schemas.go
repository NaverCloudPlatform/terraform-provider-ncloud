package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform/helper/schema"
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
	},
}

var commonCodeSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"code": {
			Type: schema.TypeString,
		},
		"code_name": {
			Type: schema.TypeString,
		},
	},
}

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

func setZone(i interface{}) map[string]interface{} {
	zone := GetZone(i)
	m := map[string]interface{}{
		"zone_no":          ncloud.StringValue(zone.ZoneNo),
		"zone_code":        ncloud.StringValue(zone.ZoneCode),
		"zone_name":        ncloud.StringValue(zone.ZoneName),
		"zone_description": ncloud.StringValue(zone.ZoneDescription),
		"region_no":        ncloud.StringValue(zone.RegionNo),
	}

	return m
}
