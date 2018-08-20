package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"reflect"
)

func validElem(i interface{}) bool {
	return reflect.ValueOf(i).Elem().IsValid()
}

func validField(f reflect.Value) bool {
	return (!f.CanAddr() || f.CanAddr() && !f.IsNil()) && f.IsValid()
}

func StringField(f reflect.Value) *string {
	if f.Kind() == reflect.Ptr && f.Type().String() == "*string" {
		return f.Interface().(*string)
	} else if f.Kind() == reflect.Slice && f.Type().String() == "string" {
		return ncloud.String(f.Interface().(string))
	}
	return nil
}

func GetCommonResponse(i interface{}) *CommonResponse {
	if i == nil || !validElem(i) {
		return &CommonResponse{}
	}
	var requestId *string
	var returnCode *string
	var returnMessage *string

	if f := reflect.ValueOf(i).Elem().FieldByName("RequestId"); validField(f) {
		requestId = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ReturnCode"); validField(f) {
		returnCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ReturnMessage"); validField(f) {
		returnMessage = StringField(f)
	}
	return &CommonResponse{
		RequestId:     requestId,
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
	}
}

func GetCommonCode(i interface{}) *CommonCode {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return &CommonCode{}
	}

	var code *string
	var codeName *string
	if f := reflect.ValueOf(i).Elem().FieldByName("Code"); validField(f) {
		code = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("CodeName"); validField(f) {
		codeName = StringField(f)
	}

	return &CommonCode{
		Code:     code,
		CodeName: codeName,
	}
}

func GetRegion(i interface{}) *Region {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return &Region{}
	}
	var regionNo *string
	var regionCode *string
	var regionName *string
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); validField(f) {
		regionNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionCode"); validField(f) {
		regionCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionName"); validField(f) {
		regionName = StringField(f)
	}

	return &Region{
		RegionNo:   regionNo,
		RegionCode: regionCode,
		RegionName: regionName,
	}
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
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneNo"); validField(f) {
		zoneNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneName"); validField(f) {
		zoneName = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneCode"); validField(f) {
		zoneCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneDescription"); validField(f) {
		zoneDescription = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); validField(f) {
		regionNo = StringField(f)
	}

	return &Zone{
		ZoneNo:          zoneNo,
		ZoneName:        zoneName,
		ZoneCode:        zoneCode,
		ZoneDescription: zoneDescription,
		RegionNo:        regionNo,
	}
}
