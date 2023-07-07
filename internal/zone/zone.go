package zone

import (
	"fmt"
	"reflect"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
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

func ParseZoneNoParameter(config *conn.ProviderConfig, d *schema.ResourceData) (*string, error) {
	if zoneCode, zoneCodeOk := d.GetOk("zone"); zoneCodeOk {
		zoneNo := GetZoneNoByCode(config, zoneCode.(string))
		if zoneNo == "" {
			return nil, fmt.Errorf("no zone data for zone_code `%s`. please change zone_code and try again", zoneCode.(string))
		}
		return ncloud.String(zoneNo), nil

	}
	return nil, nil
}

func GetZoneNoByCode(config *conn.ProviderConfig, code string) string {
	if zoneNo := zoneCache[code]; zoneNo != "" {
		return zoneNo
	}
	if zone, err := GetZoneByCode(config, code); err == nil && zone != nil {
		zoneCache[code] = *zone.ZoneNo
		return *zone.ZoneNo
	}
	return ""
}

func getZoneCodeByNo(config *conn.ProviderConfig, no string) string {
	if zoneCode := zoneCache[no]; zoneCode != "" {
		return zoneCode
	}
	if zone, err := getZoneByNo(config, no); err == nil && zone != nil {
		zoneCache[no] = *zone.ZoneCode
		return *zone.ZoneCode
	}
	return ""
}

func getZoneByNo(config *conn.ProviderConfig, no string) (*Zone, error) {
	zonesList, err := GetZones(config)
	if err != nil {
		return nil, err
	}

	var filteredZone *Zone
	for _, zone := range zonesList {
		if zone.ZoneNo != nil && no == *zone.ZoneNo {
			filteredZone = zone
			break
		}
	}
	return filteredZone, nil
}

func GetZoneByCode(config *conn.ProviderConfig, code string) (*Zone, error) {
	zonesList, err := GetZones(config)
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

func GetZones(config *conn.ProviderConfig) ([]*Zone, error) {
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

func GetZone(i interface{}) *Zone {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return &Zone{}
	}
	var zoneNo *string
	var zoneDescription *string
	var zoneName *string
	var zoneCode *string
	var regionNo *string
	var regionCode *string

	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneNo"); ValidField(f) {
		zoneNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneName"); ValidField(f) {
		zoneName = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneCode"); ValidField(f) {
		zoneCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneDescription"); ValidField(f) {
		zoneDescription = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); ValidField(f) {
		regionNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionCode"); ValidField(f) {
		regionCode = StringField(f)
	}

	return &Zone{
		ZoneNo:          zoneNo,
		ZoneName:        zoneName,
		ZoneCode:        zoneCode,
		ZoneDescription: zoneDescription,
		RegionNo:        regionNo,
		RegionCode:      regionCode,
	}
}
