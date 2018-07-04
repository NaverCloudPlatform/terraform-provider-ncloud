package ncloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"os"
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

var regionCache = make(map[string]string)

func parseRegionNoParameter(conn *sdk.Conn, d *schema.ResourceData) string {
	if paramRegionNo, regionNoOk := d.GetOk("region_no"); regionNoOk {
		return paramRegionNo.(string)
	}

	if regionCode, regionCodeOk := d.GetOk("region_code"); regionCodeOk {
		return getRegionNoByCode(conn, regionCode.(string))
	}

	// provider region
	if regionCode := os.Getenv("NCLOUD_REGION"); regionCode != "" {
		return getRegionNoByCode(conn, regionCode)
	}

	return ""
}

func getRegionNoByCode(conn *sdk.Conn, code string) string {
	if regionNo := regionCache[code]; regionNo != "" {
		return regionNo
	}
	if region, err := getRegionByCode(conn, code); err != nil {
		regionCache[code] = region.RegionNo
		return region.RegionNo
	}
	return ""
}

func getRegionByCode(conn *sdk.Conn, code string) (*common.Region, error) {
	resp, err := conn.GetRegionList()
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("no matching regions found")
	}
	regionList := resp.RegionList

	var filteredRegion common.Region
	for _, region := range regionList {
		if code == region.RegionCode {
			filteredRegion = region
			break
		}
	}

	return &filteredRegion, nil
}

var zoneCache = make(map[string]string)

func parseZoneNoParameter(conn *sdk.Conn, d *schema.ResourceData) string {
	if zoneNo, zoneNoOk := d.GetOk("zone_no"); zoneNoOk {
		return zoneNo.(string)
	}

	if zoneCode, zoneCodeOk := d.GetOk("zone_code"); zoneCodeOk {
		return getZoneNoByCode(conn, zoneCode.(string))
	}
	return ""
}

func getZoneNoByCode(conn *sdk.Conn, code string) string {
	if zoneNo := zoneCache[code]; zoneNo != "" {
		return zoneNo
	}
	if zone, err := getZoneByCode(conn, code); err != nil {
		zoneCache[code] = zone.ZoneNo
		return zone.ZoneNo
	}
	return ""
}

func getZoneByCode(conn *sdk.Conn, code string) (*common.Zone, error) {
	zonesList, err := getZones(conn)
	if err != nil {
		return nil, err
	}

	var filteredZone common.Zone
	for _, zone := range zonesList {
		if code == zone.ZoneCode {
			filteredZone = zone
			break
		}
	}
	return &filteredZone, nil
}

func getZones(conn *sdk.Conn) ([]common.Zone, error) {
	resp, err := conn.GetZoneList("")
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching zones found")
	}

	return resp.Zone, nil
}
