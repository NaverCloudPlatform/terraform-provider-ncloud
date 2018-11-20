package ncloud

import (
	"reflect"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
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

func flattenZones(zones []*Zone) []map[string]interface{} {
	var s []map[string]interface{}

	for _, zone := range zones {
		mapping := flattenZone(zone)
		s = append(s, mapping)
	}

	return s
}

func flattenZone(i interface{}) map[string]interface{} {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return map[string]interface{}{}
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

	return map[string]interface{}{
		"zone_no":          ncloud.StringValue(zoneNo),
		"zone_code":        ncloud.StringValue(zoneCode),
		"zone_name":        ncloud.StringValue(zoneName),
		"zone_description": ncloud.StringValue(zoneDescription),
		"region_no":        ncloud.StringValue(regionNo),
	}
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

func flattenServerImages(serverImages []*server.Product) []map[string]interface{} {
	var s []map[string]interface{}
	for _, product := range serverImages {
		mapping := map[string]interface{}{
			"product_code":            ncloud.StringValue(product.ProductCode),
			"product_name":            ncloud.StringValue(product.ProductName),
			"product_type":            flattenCommonCode(product.ProductType),
			"product_description":     ncloud.StringValue(product.ProductDescription),
			"infra_resource_type":     flattenCommonCode(product.InfraResourceType),
			"cpu_count":               int(ncloud.Int32Value(product.CpuCount)),
			"memory_size":             int(ncloud.Int64Value(product.MemorySize)),
			"base_block_storage_size": int(ncloud.Int64Value(product.BaseBlockStorageSize)),
			"platform_type":           flattenCommonCode(product.PlatformType),
			"os_information":          ncloud.StringValue(product.OsInformation),
			"add_block_storage_size":  int(ncloud.Int64Value(product.AddBlockStorageSize)),
		}

		s = append(s, mapping)
	}

	return s
}

func flattenNasVolumeInstances(nasVolumeInstances []*server.NasVolumeInstance) []map[string]interface{} {
	var s []map[string]interface{}

	for _, nasVolume := range nasVolumeInstances {
		mapping := map[string]interface{}{
			"nas_volume_instance_no":         ncloud.StringValue(nasVolume.NasVolumeInstanceNo),
			"nas_volume_instance_status":     flattenCommonCode(nasVolume.NasVolumeInstanceStatus),
			"create_date":                    ncloud.StringValue(nasVolume.CreateDate),
			"nas_volume_description":         ncloud.StringValue(nasVolume.NasVolumeInstanceDescription),
			"volume_allotment_protocol_type": flattenCommonCode(nasVolume.VolumeAllotmentProtocolType),
			"volume_name":                    ncloud.StringValue(nasVolume.VolumeName),
			"volume_total_size":              int(ncloud.Int64Value(nasVolume.VolumeTotalSize)),
			"volume_size":                    int(ncloud.Int64Value(nasVolume.VolumeSize)),
			"volume_use_size":                int(ncloud.Int64Value(nasVolume.VolumeUseSize)),
			"volume_use_ratio":               ncloud.Float32Value(nasVolume.VolumeUseRatio),
			"snapshot_volume_size":           ncloud.Int64Value(nasVolume.SnapshotVolumeSize),
			"snapshot_volume_use_size":       ncloud.Int64Value(nasVolume.SnapshotVolumeUseSize),
			"snapshot_volume_use_ratio":      ncloud.Float32Value(nasVolume.SnapshotVolumeUseRatio),
			"is_snapshot_configuration":      ncloud.BoolValue(nasVolume.IsSnapshotConfiguration),
			"is_event_configuration":         ncloud.BoolValue(nasVolume.IsEventConfiguration),
			"zone":                           flattenZone(nasVolume.Zone),
			"region":                         flattenRegion(nasVolume.Region),
		}
		if len(nasVolume.NasVolumeInstanceCustomIpList) > 0 {
			mapping["nas_volume_instance_custom_ip_list"] = flattenCustomIPList(nasVolume.NasVolumeInstanceCustomIpList)
		}

		s = append(s, mapping)
	}

	return s
}

func expandLoadBalancerRuleParams(list []interface{}) ([]*loadbalancer.LoadBalancerRuleParameter, error) {
	lbRuleList := make([]*loadbalancer.LoadBalancerRuleParameter, 0, len(list))

	for _, v := range list {
		lbRule := new(loadbalancer.LoadBalancerRuleParameter)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "protocol_type_code":
				lbRule.ProtocolTypeCode = ncloud.String(value.(string))
			case "load_balancer_port":
				lbRule.LoadBalancerPort = ncloud.Int32(int32(value.(int)))
			case "server_port":
				lbRule.ServerPort = ncloud.Int32(int32(value.(int)))
			case "l7_health_check_path":
				lbRule.L7HealthCheckPath = ncloud.String(value.(string))
			case "certificate_name":
				lbRule.CertificateName = ncloud.String(value.(string))
			case "proxy_protocol_use_yn":
				lbRule.ProxyProtocolUseYn = ncloud.String(value.(string))
			}
		}
		lbRuleList = append(lbRuleList, lbRule)
	}

	return lbRuleList, nil
}

func flattenLoadBalancerRuleList(lbRuleList []*loadbalancer.LoadBalancerRule) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(lbRuleList))

	for _, r := range lbRuleList {
		rule := map[string]interface{}{
			"protocol_type":         flattenCommonCode(r.ProtocolType),
			"load_balancer_port":    ncloud.Int32Value(r.LoadBalancerPort),
			"server_port":           ncloud.Int32Value(r.ServerPort),
			"l7_health_check_path":  ncloud.StringValue(r.L7HealthCheckPath),
			"certificate_name":      ncloud.StringValue(r.CertificateName),
			"proxy_protocol_use_yn": ncloud.StringValue(r.ProxyProtocolUseYn),
		}
		list = append(list, rule)
	}

	return list
}

func flattenLoadBalancedServerInstanceList(loadBalancedServerInstanceList []*loadbalancer.LoadBalancedServerInstance) []string {
	list := make([]string, 0, len(loadBalancedServerInstanceList))

	for _, instance := range loadBalancedServerInstanceList {
		list = append(list, ncloud.StringValue(instance.ServerInstance.ServerInstanceNo))
	}

	return list
}

func expandTagListParams(tl []interface{}) ([]*server.InstanceTagParameter, error) {
	tagList := make([]*server.InstanceTagParameter, 0, len(tl))

	for _, v := range tl {
		tag := new(server.InstanceTagParameter)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "tag_key":
				tag.TagKey = ncloud.String(value.(string))
			case "tag_value":
				tag.TagValue = ncloud.String(value.(string))
			}
		}
		tagList = append(tagList, tag)
	}

	return tagList, nil
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
