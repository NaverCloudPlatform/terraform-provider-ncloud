package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Zone struct {
	ZoneNo          *string `json:"zoneNo,omitempty"`
	ZoneName        *string `json:"zoneName,omitempty"`
	ZoneCode        *string `json:"zoneCode,omitempty"`
	ZoneDescription *string `json:"zoneDescription,omitempty"`
	RegionNo        *string `json:"regionNo,omitempty"`
	RegionCode      *string `json:"regionCode,omitempty"`
}

var zoneCache = make(map[string]string)

func parseZoneNoParameter(config *ProviderConfig, d *schema.ResourceData) (*string, error) {
	if zoneCode, zoneCodeOk := d.GetOk("zone"); zoneCodeOk {
		zoneNo := getZoneNoByCode(config, zoneCode.(string))
		if zoneNo == "" {
			return nil, fmt.Errorf("no zone data for zone_code `%s`. please change zone_code and try again", zoneCode.(string))
		}
		return ncloud.String(zoneNo), nil

	}
	return nil, nil
}

func getZoneNoByCode(config *ProviderConfig, code string) string {
	if zoneNo := zoneCache[code]; zoneNo != "" {
		return zoneNo
	}
	if zone, err := getZoneByCode(config, code); err == nil && zone != nil {
		zoneCache[code] = *zone.ZoneNo
		return *zone.ZoneNo
	}
	return ""
}

func getZoneByCode(config *ProviderConfig, code string) (*Zone, error) {
	zonesList, err := getZones(config)
	if err != nil {
		return nil, err
	}

	var filteredZone *Zone
	for _, zone := range zonesList {
		if zone.ZoneCode != nil && code == *zone.ZoneCode {
			filteredZone = zone
			break
		}
	}
	return filteredZone, nil
}

func getZones(config *ProviderConfig) ([]*Zone, error) {
	var zones []*Zone
	var err error

	if config.SupportVPC == true {
		zones, err = getVpcZones(config)
	} else {
		zones, err = getClassicZones(config)
	}

	if err != nil {
		return nil, err
	}

	if len(zones) == 0 {
		return nil, fmt.Errorf("no matching zones found")
	}

	return zones, nil
}
