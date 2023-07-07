package server

import (
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

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

	if result[0] != "25363" {
		t.Fatalf("expected configuration_no to be 25363, but was %s", result[0])
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

	result := zone.FlattenZone(expanded)

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

func TestExpandTagListParams(t *testing.T) {
	lbrulelist := []interface{}{
		map[string]interface{}{
			"tag_key":   "dev",
			"tag_value": "web",
		},
		map[string]interface{}{
			"tag_key":   "prod",
			"tag_value": "auth",
		},
	}

	result, _ := expandTagListParams(lbrulelist)

	if result == nil {
		t.Fatal("result was nil")
	}

	if len(result) != 2 {
		t.Fatalf("expected result had %d elements, but got %d", 2, len(result))
	}

	tag := result[0]
	if *tag.TagKey != "dev" {
		t.Fatalf("expected result ProtocolTypeCode to be dev, but was %s", *tag.TagKey)
	}

	if *tag.TagValue != "web" {
		t.Fatalf("expected result ProtocolTypeCode to be web, but was %s", *tag.TagValue)
	}

	tag = result[1]
	if *tag.TagKey != "prod" {
		t.Fatalf("expected result ProtocolTypeCode to be prod, but was %s", *tag.TagKey)
	}

	if *tag.TagValue != "auth" {
		t.Fatalf("expected result ProtocolTypeCode to be auth, but was %s", *tag.TagValue)
	}
}

func TestFlattenMapByKey(t *testing.T) {
	expanded := &server.CommonCode{
		Code: ncloud.String("test"),
	}

	result := flattenMapByKey(expanded, "code")

	if result == nil {
		t.Fatal("result was nil")
	}

	if *result != "test" {
		t.Fatalf("result expected 'test' but was %s", *result)
	}
}
