package common

import (
	"reflect"
	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcebuild"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ExpandStringInterfaceList(i []interface{}) []*string {
	vs := make([]*string, 0, len(i))
	for _, v := range i {
		if v == nil {
			vs = append(vs, nil)
			continue
		}

		switch v := v.(type) {
		case *string:
			vs = append(vs, v)
		default:
			vs = append(vs, ncloud.String(v.(string)))
		}

	}
	return vs
}

func FlattenCommonCode(i interface{}) map[string]interface{} {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return map[string]interface{}{}
	}

	var code *string
	var codeName *string

	if f := reflect.ValueOf(i).Elem().FieldByName("Code"); ValidField(f) {
		code = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("CodeName"); ValidField(f) {
		codeName = StringField(f)
	}

	return map[string]interface{}{
		"code":      ncloud.StringValue(code),
		"code_name": ncloud.StringValue(codeName),
	}
}

func flattenAccessControlGroups(accessControlGroups []*server.AccessControlGroup) []string {
	var s []string
	for _, accessControlGroup := range accessControlGroups {
		s = append(s, ncloud.StringValue(accessControlGroup.AccessControlGroupConfigurationNo))
	}

	return s
}

func flattenMemberServerImages(memberServerImages []*server.MemberServerImage) []string {
	var s []string

	for _, m := range memberServerImages {
		s = append(s, ncloud.StringValue(m.MemberServerImageNo))
	}

	return s
}

func flattenCustomIPList(customIPList []*server.NasVolumeInstanceCustomIp) []string {
	var a []string

	for _, v := range customIPList {
		a = append(a, ncloud.StringValue(v.CustomIp))
	}

	return a
}

func flattenServerProducts(serverProduct []*server.Product) []string {
	var s []string

	for _, product := range serverProduct {
		s = append(s, ncloud.StringValue(product.ProductCode))
	}

	return s
}

func flattenNasVolumeInstances(nasVolumeInstances []*server.NasVolumeInstance) []string {
	var s []string

	for _, nasVolume := range nasVolumeInstances {
		s = append(s, ncloud.StringValue(nasVolume.NasVolumeInstanceNo))
	}

	return s
}

func FlattenRegions(regions []*conn.Region) []map[string]interface{} {
	var s []map[string]interface{}

	for _, region := range regions {
		mapping := flattenRegion(region)
		s = append(s, mapping)
	}

	return s
}

func flattenRegion(i interface{}) map[string]interface{} {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return map[string]interface{}{}
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

	m := map[string]interface{}{
		"region_no":   ncloud.StringValue(regionNo),
		"region_code": ncloud.StringValue(regionCode),
		"region_name": ncloud.StringValue(regionName),
	}

	return m
}

func flattenInstanceTagList(tagList []*server.InstanceTag) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(tagList))

	for _, r := range tagList {
		tag := map[string]interface{}{
			"tag_key":   ncloud.StringValue(r.TagKey),
			"tag_value": ncloud.StringValue(r.TagValue),
		}
		list = append(list, tag)
	}

	return list
}

func FlattenArrayStructByKey(list interface{}, key string) []*string {
	s := make([]*string, 0)

	if list == nil {
		return s
	}
	arr := ConvertToArrayMap(list)
	for _, v := range arr {
		s = append(s, ncloud.String(v[key].(string)))
	}

	return s
}

func GetInt32FromString(v interface{}, ok bool) *int32 {
	if !ok {
		return nil
	}

	intV, err := strconv.Atoi(v.(string))
	if err == nil {
		return ncloud.Int32(int32(intV))
	} else {
		return nil
	}
}

func ExpandStringInterfaceListToInt32List(list []interface{}) []*int32 {
	res := make([]*int32, 0)
	for _, v := range list {
		if v == nil {
			res = append(res, nil)
			continue
		}

		intV, err := strconv.Atoi(v.(string))
		if err == nil {
			res = append(res, ncloud.Int32(int32(intV)))
		}
	}
	return res
}

func expandSourceBuildEnvVarsParams(eVars []interface{}) ([]*sourcebuild.ProjectEnvEnvVars, error) {
	envVars := make([]*sourcebuild.ProjectEnvEnvVars, 0, len(eVars))

	for _, v := range eVars {
		env := new(sourcebuild.ProjectEnvEnvVars)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "key":
				env.Key = ncloud.String(value.(string))
			case "value":
				env.Value = ncloud.String(value.(string))
			}
		}
		envVars = append(envVars, env)
	}

	return envVars, nil
}
