package zone

import (
	"reflect"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
)

func flattenZones(zones []*Zone) []map[string]interface{} {
	var s []map[string]interface{}

	for _, zone := range zones {
		mapping := FlattenZone(zone)
		s = append(s, mapping)
	}

	return s
}

func FlattenZone(i interface{}) map[string]interface{} {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return map[string]interface{}{}
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

	return map[string]interface{}{
		"zone_no":          ncloud.StringValue(zoneNo),
		"zone_code":        ncloud.StringValue(zoneCode),
		"zone_name":        ncloud.StringValue(zoneName),
		"zone_description": ncloud.StringValue(zoneDescription),
		"region_no":        ncloud.StringValue(regionNo),
		"region_code":      ncloud.StringValue(regionCode),
	}
}
