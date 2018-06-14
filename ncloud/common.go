package ncloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/hashicorp/terraform/helper/schema"
)

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

func logErrorResponse(tag string, err error, args interface{}) {
	param, _ := json.Marshal(args)
	log.Printf("[ERROR] %s error params=%s, err=%s", tag, param, err)
}

func logCommonResponse(tag string, args interface{}, commonResponse common.CommonResponse) {
	param, _ := json.Marshal(args)
	result := fmt.Sprintf("RequestID: %s, ReturnCode: %d, ReturnMessage: %s", commonResponse.RequestID, commonResponse.ReturnCode, commonResponse.ReturnMessage)
	log.Printf("[DEBUG] %s success params=%s, response=%s", tag, param, result)
}

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

func setCommonCode(cc common.CommonCode) map[string]interface{} {
	m := map[string]interface{}{
		"code":      cc.Code,
		"code_name": cc.CodeName,
	}

	return m
}

func setZone(zone common.Zone) map[string]interface{} {
	m := map[string]interface{}{
		"zone_no":          zone.ZoneNo,
		"zone_name":        zone.ZoneName,
		"zone_description": zone.ZoneDescription,
	}

	return m
}

func setRegion(region common.Region) map[string]interface{} {
	m := map[string]interface{}{
		"region_no":   region.RegionNo,
		"region_code": region.RegionCode,
		"region_name": region.RegionName,
	}

	return m
}
