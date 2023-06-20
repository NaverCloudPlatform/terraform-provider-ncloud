package classicloadbalancer

import (
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
)

func TestExpandLoadBalancerRuleParams(t *testing.T) {
	lbrulelist := []interface{}{
		map[string]interface{}{
			"protocol_type":        "HTTP",
			"load_balancer_port":   80,
			"server_port":          80,
			"l7_health_check_path": "/monitor/l7check",
		},
		map[string]interface{}{
			"protocol_type":        "HTTPS",
			"load_balancer_port":   443,
			"server_port":          443,
			"l7_health_check_path": "/monitor/l7check",
			"certificate_name":     "aaa",
		},
	}

	result, _ := expandLoadBalancerRuleParams(lbrulelist)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	r := result[0]
	if *r.ProtocolTypeCode != "HTTP" {
		t.Fatalf("expected result ProtocolTypeCode to be HTTP, but was %s", *r.ProtocolTypeCode)
	}

	if *r.LoadBalancerPort != 80 {
		t.Fatalf("expected result LoadBalancerPort to be 80, but was %d", *r.LoadBalancerPort)
	}

	if *r.ServerPort != 80 {
		t.Fatalf("expected result ServerPort to be 80, but was %d", *r.ServerPort)
	}

	if *r.L7HealthCheckPath != "/monitor/l7check" {
		t.Fatalf("expected result L7HealthCheckPath to be '/monitor/l7check', but was %s", *r.L7HealthCheckPath)
	}

	if r.CertificateName != nil {
		t.Fatalf("expected result CertificateName to be nil, but was %s", *r.CertificateName)
	}

	r = result[1]
	if *r.ProtocolTypeCode != "HTTPS" {
		t.Fatalf("expected result ProtocolTypeCode to be HTTPS, but was %s", *r.ProtocolTypeCode)
	}

	if *r.LoadBalancerPort != 443 {
		t.Fatalf("expected result LoadBalancerPort to be 443, but was %d", *r.LoadBalancerPort)
	}

	if *r.ServerPort != 443 {
		t.Fatalf("expected result ServerPort to be 443, but was %d", *r.ServerPort)
	}

	if *r.L7HealthCheckPath != "/monitor/l7check" {
		t.Fatalf("expected result L7HealthCheckPath to be '/monitor/l7check', but was %s", *r.L7HealthCheckPath)
	}

	if *r.CertificateName != "aaa" {
		t.Fatalf("expected result CertificateName to be aaa, but was %s", *r.CertificateName)
	}
}

func TestFlattenLoadBalancerRuleList(t *testing.T) {
	expanded := []*loadbalancer.LoadBalancerRule{
		{
			ProtocolType: &loadbalancer.CommonCode{
				Code:     ncloud.String("HTTP"),
				CodeName: ncloud.String("http"),
			},
			LoadBalancerPort:   ncloud.Int32(80),
			ServerPort:         ncloud.Int32(80),
			L7HealthCheckPath:  ncloud.String("/monitor/l7check"),
			CertificateName:    ncloud.String("aaa"),
			ProxyProtocolUseYn: ncloud.String("Y"),
		},
	}

	result := flattenLoadBalancerRuleList(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	r := result[0]

	if r["load_balancer_port"] != int32(80) {
		t.Fatalf("expected result load_balancer_port to be 80, but was %d", r["load_balancer_port"])
	}

	if r["server_port"] != int32(80) {
		t.Fatalf("expected result server_port to be 80, but was %d", r["server_port"])
	}

	if r["l7_health_check_path"] != "/monitor/l7check" {
		t.Fatalf("expected result l7_health_check_path to be /monitor/l7check, but was %s", r["l7_health_check_path"])
	}

	if r["certificate_name"] != "aaa" {
		t.Fatalf("expected result certificate_name to be aaa, but was %s", r["certificate_name"])
	}

	if r["proxy_protocol_use_yn"] != "Y" {
		t.Fatalf("expected result proxy_protocol_use_yn to be Y, but was %s", r["proxy_protocol_use_yn"])
	}

}

func TestFlattenLoadBalancedServerInstanceList(t *testing.T) {
	expanded := []*loadbalancer.LoadBalancedServerInstance{
		{
			ServerInstance: &loadbalancer.ServerInstance{
				ServerInstanceNo: ncloud.String("123456"),
			},
		},
		{
			ServerInstance: &loadbalancer.ServerInstance{
				ServerInstanceNo: ncloud.String("234567"),
			},
		},
	}

	result := flattenLoadBalancedServerInstanceList(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	if result[0] != "123456" {
		t.Fatalf("expected result load_balancer_port to be '123456', but was %s", result[0])
	}

	if result[1] != "234567" {
		t.Fatalf("expected result load_balancer_port to be '234567', but was %s", result[1])
	}
}
