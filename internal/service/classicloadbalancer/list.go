package classicloadbalancer

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
)

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
