package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"reflect"
)

func StringField(f reflect.Value) *string {
	if f.Kind() == reflect.Ptr && f.Type().String() == "*string" {
		return f.Interface().(*string)
	} else if f.Kind() == reflect.Slice && f.Type().String() == "string" {
		return ncloud.String(f.Interface().(string))
	}
	return nil
}

func GetCommonResponse(i interface{}) *CommonResponse {
	if i == nil {
		return &CommonResponse{}
	}
	var requestId *string
	var returnCode *string
	var returnMessage *string
	if f := reflect.ValueOf(i).Elem().FieldByName("RequestId"); !f.IsNil() && f.IsValid() {
		requestId = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ReturnCode"); !f.IsNil() && f.IsValid() {
		returnCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ReturnMessage"); !f.IsNil() && f.IsValid() {
		returnMessage = StringField(f)
	}
	return &CommonResponse{
		RequestId:     requestId,
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
	}
}

func GetCommonCode(i interface{}) *CommonCode {
	if i == nil {
		return &CommonCode{}
	}
	var code *string
	var codeName *string
	if f := reflect.ValueOf(i).Elem().FieldByName("Code"); !f.IsNil() && f.IsValid() {
		code = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("CodeName"); !f.IsNil() && f.IsValid() {
		codeName = StringField(f)
	}

	return &CommonCode{
		Code:     code,
		CodeName: codeName,
	}
}

func GetRegion(i interface{}) *Region {
	if i == nil {
		return &Region{}
	}
	var regionNo *string
	var regionCode *string
	var regionName *string
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); !f.IsNil() && f.IsValid() {
		regionNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionCode"); !f.IsNil() && f.IsValid() {
		regionCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionName"); !f.IsNil() && f.IsValid() {
		regionName = StringField(f)
	}

	return &Region{
		RegionNo:   regionNo,
		RegionCode: regionCode,
		RegionName: regionName,
	}
}

func GetZone(i interface{}) *Zone {
	if i == nil {
		return &Zone{}
	}
	var zoneNo *string
	var zoneDescription *string
	var zoneName *string
	var zoneCode *string
	var regionNo *string
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneNo"); !f.IsNil() && f.IsValid() {
		zoneNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneName"); !f.IsNil() && f.IsValid() {
		zoneName = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneCode"); !f.IsNil() && f.IsValid() {
		zoneCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ZoneDescription"); !f.IsNil() && f.IsValid() {
		zoneDescription = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); !f.IsNil() && f.IsValid() {
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
