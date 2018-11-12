package ncloud

import (
	"reflect"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
)

func TestExpandStringInterfaceList(t *testing.T) {
	initialList := []string{"1111", "2222", "3333"}
	l := make([]interface{}, len(initialList))
	for i, v := range initialList {
		l[i] = v
	}
	stringList := expandStringInterfaceList(l)
	expected := []*string{
		ncloud.String("1111"),
		ncloud.String("2222"),
		ncloud.String("3333"),
	}
	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func TestFlattenCommonCode(t *testing.T) {
	expanded := &server.CommonCode{
		Code:     ncloud.String("code"),
		CodeName: ncloud.String("codename"),
	}

	result := flattenCommonCode(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if result["code"] != "code" {
		t.Fatalf("expected result code to be code, but was %s", result["code"])
	}

	if result["code_name"] != "codename" {
		t.Fatalf("expected result code_name to be codename, but was %s", result["code_name"])
	}
}

func TestFlattenAccessControlRules(t *testing.T) {
	expected := []*server.AccessControlRule{
		{
			AccessControlRuleConfigurationNo: ncloud.String("25363"),
			ProtocolType: &server.CommonCode{
				Code:     ncloud.String("TCP"),
				CodeName: ncloud.String("tcp"),
			},
			SourceIp:                               ncloud.String("0.0.0.0/0"),
			SourceAccessControlRuleConfigurationNo: ncloud.String("4964"),
			SourceAccessControlRuleName:            ncloud.String("ncloud-default-acg"),
			DestinationPort:                        ncloud.String("1-65535"),
			AccessControlRuleDescription:           ncloud.String("for test"),
		},
	}

	result := flattenAccessControlRules(expected)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 1 {
		t.Fatalf("expected result had %d elements, but got %d", 1, len(result))
	}

	if result[0]["access_control_rule_configuration_no"] != "25363" {
		t.Fatalf("expected access_control_rule_configuration_no to be 25363, but was %s", result[0]["access_control_rule_configuration_no"])
	}

	if result[0]["source_ip"] != "0.0.0.0/0" {
		t.Fatalf("expected source_ip to be 0.0.0.0/0, but was %s", result[0]["source_ip"])
	}

	if result[0]["destination_port"] != "1-65535" {
		t.Fatalf("expected destination_port to be 1-65535, but was %s", result[0]["destination_port"])
	}

	if result[0]["source_access_control_rule_configuration_no"] != "4964" {
		t.Fatalf("expected source_access_control_rule_configuration_no to be 4964, but was %s", result[0]["source_access_control_rule_configuration_no"])
	}

	if result[0]["source_access_control_rule_name"] != "ncloud-default-acg" {
		t.Fatalf("expected source_access_control_rule_name to be ncloud-default-acg, but was %s", result[0]["source_access_control_rule_name"])
	}

	if result[0]["access_control_rule_description"] != "for test" {
		t.Fatalf("expected access_control_rule_description to be 'for test', but was %s", result[0]["access_control_rule_description"])
	}
}

func TestFlattenAccessControlGroups(t *testing.T) {
	expected := []*server.AccessControlGroup{
		{
			AccessControlGroupConfigurationNo: ncloud.String("4964"),
			AccessControlGroupName:            ncloud.String("ncloud-default-acg"),
			AccessControlGroupDescription:     ncloud.String("for test"),
			IsDefaultGroup:                    ncloud.Bool(true),
			CreateDate:                        ncloud.String("2017-02-23T10:25:39+0900"),
		},
		{
			AccessControlGroupConfigurationNo: ncloud.String("30067"),
			AccessControlGroupName:            ncloud.String("httpdport"),
			AccessControlGroupDescription:     ncloud.String("for httpd test"),
			IsDefaultGroup:                    ncloud.Bool(false),
			CreateDate:                        ncloud.String("2018-01-07T10:17:14+0900"),
		},
	}

	result := flattenAccessControlGroups(expected)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	if result[0]["access_control_group_configuration_no"] != "4964" {
		t.Fatalf("expected access_control_group_configuration_no to be 4964, but was %s", result[0]["access_control_group_configuration_no"])
	}

	if result[0]["access_control_group_name"] != "ncloud-default-acg" {
		t.Fatalf("expected access_control_group_name to be ncloud-default-acg, but was %s", result[0]["access_control_group_name"])
	}

	if result[0]["access_control_group_description"] != "for test" {
		t.Fatalf("expected access_control_group_description to be 'for test', but was %s", result[0]["access_control_group_description"])
	}

	if result[0]["is_default_group"] != true {
		t.Fatalf("expected is_default_group to be true, but was %b", result[0]["is_default_group"])
	}

	if result[0]["create_date"] != "2017-02-23T10:25:39+0900" {
		t.Fatalf("expected create_date to be 2017-02-23T10:25:39+0900, but was %s", result[0]["create_date"])
	}
}

func TestFlattenRegion(t *testing.T) {
	expanded := &server.Region{
		RegionNo:   ncloud.String("1"),
		RegionCode: ncloud.String("KR"),
		RegionName: ncloud.String("Korea"),
	}

	result := flattenRegion(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if result["region_no"] != "1" {
		t.Fatalf("expected result region_no to be 1, but was %s", result["region_no"])
	}

	if result["region_code"] != "KR" {
		t.Fatalf("expected result region_code to be KR, but was %s", result["region_code"])
	}

	if result["region_name"] != "Korea" {
		t.Fatalf("expected result region_name to be Korea, but was %s", result["region_name"])
	}
}

func TestFlattenZone(t *testing.T) {
	expanded := &server.Zone{
		ZoneNo:          ncloud.String("3"),
		ZoneName:        ncloud.String("KR-2"),
		ZoneCode:        ncloud.String("KR-2"),
		ZoneDescription: ncloud.String("평촌 zone"),
		RegionNo:        ncloud.String("1"),
	}

	result := flattenZone(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if result["zone_no"] != "3" {
		t.Fatalf("expected result zone_no to be 3, but was %s", result["zone_no"])
	}

	if result["zone_code"] != "KR-2" {
		t.Fatalf("expected result zone_code to be KR-2, but was %s", result["zone_code"])
	}

	if result["zone_name"] != "KR-2" {
		t.Fatalf("expected result zone_name to be KR-2, but was %s", result["zone_name"])
	}

	if result["zone_description"] != "평촌 zone" {
		t.Fatalf("expected result zone_description to be 평촌 zone, but was %s", result["zone_description"])
	}

	if result["region_no"] != "1" {
		t.Fatalf("expected result region_no to be 1, but was %s", result["region_no"])
	}
}

func TestFlattenMemberServerImages(t *testing.T) {
}

func TestFlattenCustomIPList(t *testing.T) {
}

func TestFlattenServerImages(t *testing.T) {
}

func TestFlattenNasVolumeInstances(t *testing.T) {
}
func TestExpandLoadBalancerRuleParams(t *testing.T) {
}

func TestFlattenLoadBalancerRuleList(t *testing.T) {
}

func TestFlattenLoadBalancedServerInstanceList(t *testing.T) {
}

func TestExpandTagListParams(t *testing.T) {
}

func TestFlattenInstanceTagList(t *testing.T) {
}
