package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"reflect"
	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcebuild"
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

func flattenAccessControlRules(accessControlRules []*server.AccessControlRule) []string {
	var s []string

	for _, accessControlRule := range accessControlRules {
		s = append(s, ncloud.StringValue(accessControlRule.AccessControlRuleConfigurationNo))
	}

	return s
}

func flattenAccessControlGroups(accessControlGroups []*server.AccessControlGroup) []string {
	var s []string
	for _, accessControlGroup := range accessControlGroups {
		s = append(s, ncloud.StringValue(accessControlGroup.AccessControlGroupConfigurationNo))
	}

	return s
}

func flattenRegions(regions []*Region) []map[string]interface{} {
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

	return map[string]interface{}{
		"zone_no":          ncloud.StringValue(zoneNo),
		"zone_code":        ncloud.StringValue(zoneCode),
		"zone_name":        ncloud.StringValue(zoneName),
		"zone_description": ncloud.StringValue(zoneDescription),
		"region_no":        ncloud.StringValue(regionNo),
		"region_code":      ncloud.StringValue(regionCode),
	}
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

func expandLoadBalancerRuleParams(list []interface{}) ([]*loadbalancer.LoadBalancerRuleParameter, error) {
	lbRuleList := make([]*loadbalancer.LoadBalancerRuleParameter, 0, len(list))

	for _, v := range list {
		lbRule := new(loadbalancer.LoadBalancerRuleParameter)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "protocol_type":
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
			"protocol_type":         ncloud.StringValue(r.ProtocolType.Code),
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

func flattenMapByKey(i interface{}, key string) *string {
	m := ConvertToMap(i)
	if m[key] != nil {
		return ncloud.String(m[key].(string))
	} else {
		return nil
	}
}

func flattenArrayStructByKey(list interface{}, key string) []*string {
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

func getInt32FromString(v interface{}, ok bool) *int32 {
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

func expandStringInterfaceListToInt32List(list []interface{}) (res []*int32) {
	for _, v := range list {
		intV, err := strconv.Atoi(v.(string))
		if err == nil {
			res = append(res, ncloud.Int32(int32(intV)))
		}
	}
	return res
}

func flattenInt32ListToStringList(list []*int32) (res []*string) {
	for _, v := range list {
		res = append(res, ncloud.IntString(int(ncloud.Int32Value(v))))
	}
	return
}

func flattenNKSClusterLogInput(logInput *vnks.ClusterLogInput) []map[string]interface{} {
	if logInput == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"audit": ncloud.BoolValue(logInput.Audit),
		},
	}
}
func expandNKSClusterLogInput(logList []interface{}) *vnks.ClusterLogInput {
	if len(logList) == 0 {
		return nil
	}
	log := logList[0].(map[string]interface{})
	return &vnks.ClusterLogInput{
		Audit: ncloud.Bool(log["audit"].(bool)),
	}
}

func flattenNKSNodePoolAutoScale(ao *vnks.AutoscaleOption) (res []map[string]interface{}) {
	if ao == nil {
		return
	}
	m := map[string]interface{}{
		"enabled": ncloud.BoolValue(ao.Enabled),
		"min":     ncloud.Int32Value(ao.Min),
		"max":     ncloud.Int32Value(ao.Max),
	}
	res = append(res, m)
	return
}

func expandNKSNodePoolAutoScale(as []interface{}) *vnks.AutoscalerUpdate {
	if len(as) == 0 {
		return nil
	}
	autoScale := as[0].(map[string]interface{})
	return &vnks.AutoscalerUpdate{
		Enabled: ncloud.Bool(autoScale["enabled"].(bool)),
		Min:     ncloud.Int32(int32(autoScale["min"].(int))),
		Max:     ncloud.Int32(int32(autoScale["max"].(int))),
	}
}

func flattenNKSWorkerNodes(wns []*vnks.WorkerNode) (res []map[string]interface{}) {
	if wns == nil {
		return
	}
	for _, wn := range wns {
		m := map[string]interface{}{
			"name":              ncloud.StringValue(wn.Name),
			"instance_no":       ncloud.Int32Value(wn.Id),
			"spec":              ncloud.StringValue(wn.ServerSpec),
			"private_ip":        ncloud.StringValue(wn.PrivateIp),
			"public_ip":         ncloud.StringValue(wn.PublicIp),
			"node_status":       ncloud.StringValue(wn.K8sStatus),
			"container_version": ncloud.StringValue(wn.DockerVersion),
			"kernel_version":    ncloud.StringValue(wn.KernelVersion),
		}
		res = append(res, m)
	}

	return
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
