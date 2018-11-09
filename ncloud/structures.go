package ncloud

import (
	"reflect"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
)

func expandStringInterfaceList(i []interface{}) []*string {
	vs := make([]*string, 0, len(i))
	for _, v := range i {
		switch v.(type) {
		case *string:
			vs = append(vs, v.(*string))
		default:
			vs = append(vs, ncloud.String(v.(string)))
		}

	}
	return vs
}

func flattenCommonCode(i interface{}) map[string]interface{} {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return map[string]interface{}{}
	}

	var code *string
	var codeName *string

	if f := reflect.ValueOf(i).Elem().FieldByName("Code"); validField(f) {
		code = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("CodeName"); validField(f) {
		codeName = StringField(f)
	}

	return map[string]interface{}{
		"code":      ncloud.StringValue(code),
		"code_name": ncloud.StringValue(codeName),
	}
}

func flattenAccessControlRules(accessControlRules []*server.AccessControlRule) []map[string]interface{} {
	var s []map[string]interface{}

	for _, accessControlRule := range accessControlRules {
		mapping := map[string]interface{}{
			"access_control_rule_configuration_no":        ncloud.StringValue(accessControlRule.AccessControlRuleConfigurationNo),
			"protocol_type":                               flattenCommonCode(accessControlRule.ProtocolType),
			"source_ip":                                   ncloud.StringValue(accessControlRule.SourceIp),
			"destination_port":                            ncloud.StringValue(accessControlRule.DestinationPort),
			"source_access_control_rule_configuration_no": ncloud.StringValue(accessControlRule.SourceAccessControlRuleConfigurationNo),
			"source_access_control_rule_name":             ncloud.StringValue(accessControlRule.SourceAccessControlRuleName),
			"access_control_rule_description":             ncloud.StringValue(accessControlRule.AccessControlRuleDescription),
		}

		s = append(s, mapping)
	}

	return s
}

func flattenAccessControlGroups(accessControlGroups []*server.AccessControlGroup) []map[string]interface{} {
	var s []map[string]interface{}
	for _, accessControlGroup := range accessControlGroups {
		mapping := map[string]interface{}{
			"access_control_group_configuration_no": ncloud.StringValue(accessControlGroup.AccessControlGroupConfigurationNo),
			"access_control_group_name":             ncloud.StringValue(accessControlGroup.AccessControlGroupName),
			"access_control_group_description":      ncloud.StringValue(accessControlGroup.AccessControlGroupDescription),
			"is_default_group":                      ncloud.BoolValue(accessControlGroup.IsDefaultGroup),
			"create_date":                           ncloud.StringValue(accessControlGroup.CreateDate),
		}

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

	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); validField(f) {
		regionNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionCode"); validField(f) {
		regionCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionName"); validField(f) {
		regionName = StringField(f)
	}

	m := map[string]interface{}{
		"region_no":   ncloud.StringValue(regionNo),
		"region_code": ncloud.StringValue(regionCode),
		"region_name": ncloud.StringValue(regionName),
	}

	return m
}

func flattenMemberServerImages(memberServerImages []*server.MemberServerImage) []map[string]interface{} {
	var s []map[string]interface{}
	for _, m := range memberServerImages {
		mapping := map[string]interface{}{
			"member_server_image_no":                       ncloud.StringValue(m.MemberServerImageNo),
			"member_server_image_name":                     ncloud.StringValue(m.MemberServerImageName),
			"member_server_image_description":              ncloud.StringValue(m.MemberServerImageDescription),
			"original_server_instance_no":                  ncloud.StringValue(m.OriginalServerInstanceNo),
			"original_server_product_code":                 ncloud.StringValue(m.OriginalServerProductCode),
			"original_server_name":                         ncloud.StringValue(m.OriginalServerName),
			"original_base_block_storage_disk_type":        flattenCommonCode(m.OriginalBaseBlockStorageDiskType),
			"original_server_image_product_code":           ncloud.StringValue(m.OriginalServerImageProductCode),
			"original_os_information":                      ncloud.StringValue(m.OriginalOsInformation),
			"original_server_image_name":                   ncloud.StringValue(m.OriginalServerImageName),
			"member_server_image_status_name":              ncloud.StringValue(m.MemberServerImageStatusName),
			"member_server_image_status":                   flattenCommonCode(m.MemberServerImageStatus),
			"member_server_image_operation":                flattenCommonCode(m.MemberServerImageOperation),
			"member_server_image_platform_type":            flattenCommonCode(m.MemberServerImagePlatformType),
			"create_date":                                  ncloud.StringValue(m.CreateDate),
			"region":                                       flattenRegion(m.Region),
			"member_server_image_block_storage_total_rows": int(ncloud.Int32Value(m.MemberServerImageBlockStorageTotalRows)),
			"member_server_image_block_storage_total_size": int(ncloud.Int64Value(m.MemberServerImageBlockStorageTotalSize)),
		}

		s = append(s, mapping)
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
