package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

func validElem(i interface{}) bool {
	return reflect.ValueOf(i).Elem().IsValid()
}

func ValidField(f reflect.Value) bool {
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

	if f := reflect.ValueOf(i).Elem().FieldByName("RequestId"); ValidField(f) {
		requestId = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ReturnCode"); ValidField(f) {
		returnCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("ReturnMessage"); ValidField(f) {
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

func GetRegion(i interface{}) *conn.Region {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return &conn.Region{}
	}
	var regionNo *string
	var regionCode *string
	var regionName *string
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); ValidField(f) {
		regionNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionCode"); ValidField(f) {
		regionCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionName"); ValidField(f) {
		regionName = StringField(f)
	}

	return &conn.Region{
		RegionNo:   regionNo,
		RegionCode: regionCode,
		RegionName: regionName,
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
		switch v := v.(type) {
		case *string:
			vs = append(vs, v)
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
	_ = json.Unmarshal(b, &m)

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
	_ = json.Unmarshal(b, &m)

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

// Int64ValueFromInt32 converts an int32 pointer to a Framework Int64 value.
// A nil int32 pointer is converted to a null Int64.
func Int64ValueFromInt32(value *int32) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Null()
	}
	return basetypes.NewInt64Value(int64(*value))
}

// Int64FromInt32OrDefault converts an int32 pointer to a Framework Int64 value.
// A nil int32 pointer is converted to a zero Int64.
// Used when the optional and computed attribute have no response value
func Int64FromInt32OrDefault(value *int32) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Value(0)
	}
	return basetypes.NewInt64Value(int64(*value))
}

// StringFrameworkOrDefault converts a Framework StringValue struct to a same Framework StringValue.
// A null or unknown state is converted to a default(not aloocated) string.
// Used when the optional and computed attribute have no response value
func StringFrameworkOrDefault(value types.String) basetypes.StringValue {
	if value.IsNull() || value.IsUnknown() {
		return types.StringValue("not allocated")
	}
	return value
}

func ConvertToStringList(values basetypes.ListValue, attrValue string) []string {
	result := make([]string, 0, len(values.Elements()))

	for _, v := range values.Elements() {
		obj := v.(types.Object)
		attrs := obj.Attributes()

		name := attrs[attrValue].(types.String).ValueString()
		result = append(result, name)
	}

	return result
}
