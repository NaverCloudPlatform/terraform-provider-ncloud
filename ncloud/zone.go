package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func parseZoneNoParameter(client *NcloudAPIClient, d *schema.ResourceData) (*string, error) {
	if zoneCode, zoneCodeOk := d.GetOk("zone"); zoneCodeOk {
		zoneNo := getZoneNoByCode(client, zoneCode.(string))
		if zoneNo == "" {
			return nil, fmt.Errorf("no zone data for zone_code `%s`. please change zone_code and try again", zoneCode.(string))
		}
		return ncloud.String(zoneNo), nil

	}
	return nil, nil
}

func getZoneNoByCode(client *NcloudAPIClient, code string) string {
	if zoneNo := zoneCache[code]; zoneNo != "" {
		return zoneNo
	}
	if zone, err := getZoneByCode(client, code); err == nil && zone != nil {
		zoneCache[code] = *zone.ZoneNo
		return *zone.ZoneNo
	}
	return ""
}

func getZoneByCode(client *NcloudAPIClient, code string) (*Zone, error) {
	zonesList, err := getZones(client)
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

func getZones(client *NcloudAPIClient) ([]*Zone, error) {
	resp, err := client.server.V2Api.GetZoneList(&server.GetZoneListRequest{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("no matching zones found")
	}

	var zones []*Zone
	for _, zone := range resp.ZoneList {
		zones = append(zones, GetZone(zone))
	}

	return zones, nil
}
