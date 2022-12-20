package ncloud

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
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

// GetCommonErrorBody parse common error message
func GetCommonErrorBody(err error) (*CommonError, error) {
	sa := strings.Split(err.Error(), "Body: ")
	var errMsg string

	if len(sa) != 2 {
		return nil, fmt.Errorf("error body is incorrect: %s", err)
	}

	errMsg = sa[1]

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(errMsg), &m); err != nil {
		return nil, err
	}

	e := m["responseError"].(map[string]interface{})

	return &CommonError{
		ReturnCode:    e["returnCode"].(string),
		ReturnMessage: e["returnMessage"].(string),
	}, nil
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
	var regionCode *string

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
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionCode"); validField(f) {
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

// StringPtrOrNil return *string from interface{}
func StringPtrOrNil(v interface{}, ok bool) *string {
	if !ok {
		return nil
	}
	return ncloud.String(v.(string))
}

// Int32PtrOrNil return *int32 from interface{}
func Int32PtrOrNil(v interface{}, ok bool) *int32 {
	if !ok {
		return nil
	}

	switch i := v.(type) {
	case int:
		return ncloud.Int32(int32(i))
	case int32:
		return ncloud.Int32(i)
	case int64:
		return ncloud.Int32(int32(i))
	default:
		return ncloud.Int32(i.(int32))
	}
}

// BoolPtrOrNil return *bool from interface{}
func BoolPtrOrNil(v interface{}, ok bool) *bool {
	if !ok {
		return nil
	}
	return ncloud.Bool(v.(bool))
}

// StringListPtrOrNil Convert from interface to []*string
func StringListPtrOrNil(i interface{}, ok bool) []*string {
	if !ok {
		return nil
	}

	// Handling when not slice type
	if r := reflect.ValueOf(i); r.Kind() != reflect.Slice {
		tmp := []interface{}{r.String()}
		i = tmp
	}

	il := i.([]interface{})
	vs := make([]*string, 0, len(il))
	for _, v := range il {
		switch v.(type) {
		case *string:
			vs = append(vs, v.(*string))
		default:
			// TODO: if the value is "" in list, occur crash error.
			vs = append(vs, ncloud.String(v.(string)))
		}
	}
	return vs
}

// StringOrEmpty Get string from *pointer
func StringOrEmpty(v *string) string {
	if v != nil {
		return *v
	}

	return ""
}

// StringPtrArrToStringArr Convert []*string to []string
func StringPtrArrToStringArr(ptrArray []*string) []string {
	var arr []string
	for _, v := range ptrArray {
		arr = append(arr, *v)
	}

	return arr
}

// SetStringIfNotNilAndEmpty set value map[key] if *string pointer is not nil and not empty
func SetStringIfNotNilAndEmpty(m map[string]interface{}, k string, v *string) {
	if v != nil && len(*v) > 0 {
		m[k] = *v
	}
}

// ConvertToMap convert interface{} to map[string]interface{}
func ConvertToMap(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)

	return m
}

// ConvertToArrayMap convert interface{} to map[string]interface{}
func ConvertToArrayMap(i interface{}) []map[string]interface{} {
	if i == nil {
		return nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	var m []map[string]interface{}
	json.Unmarshal(b, &m)

	return m
}

// ExpandStringSet Takes the result of schema.Set of strings and returns a []*string
func ExpandStringSet(configured *schema.Set) []*string {
	return ExpandStringList(configured.List())
}

// ExpandStringList Takes the result of flatmap.Expand for an array of strings and returns a []*string
func ExpandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, ncloud.String(v.(string)))
		}
	}
	return vs
}
