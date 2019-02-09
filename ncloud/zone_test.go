package ncloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

func testZoneSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"zone_no": {
			Type: schema.TypeString,
		},
		"zone_code": {
			Type: schema.TypeString,
		},
	}
	return s
}

func TestParseZoneNoParameter_basic(t *testing.T) {
	testZoneCode := "KR-2"
	client, _ := testApiClient(t)

	s := testZoneSchema()
	d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
		"zone_code": testZoneCode,
	})
	if zoneNo, _ := parseZoneNoParameter(client, d); zoneNo == nil {
		t.Fatalf("zone_no should be returned when input zoneCode. input: %s", testZoneCode)
	}
}

func TestParseZoneNoParameter_inputNil(t *testing.T) {
	client, _ := testApiClient(t)

	if zoneNo, _ := parseZoneNoParameter(client, &schema.ResourceData{}); zoneNo != nil {
		t.Fatalf("zone_no should be return nil when input empty resource data. actual: %s", *zoneNo)
	}
}

func TestParseZoneNoParameter_inputZoneNo(t *testing.T) {
	testZoneNo := "1"
	client, _ := testApiClient(t)

	s := testZoneSchema()
	d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
		"zone_no": testZoneNo,
	})
	if zoneNo, _ := parseZoneNoParameter(client, d); *zoneNo != testZoneNo {
		t.Fatalf("Expected: %s, Actual: %s", testZoneNo, *zoneNo)
	}
}

func TestParseZoneNoParameter_inputUnknownZoneCode(t *testing.T) {
	client, _ := testApiClient(t)

	s := testZoneSchema()
	d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
		"zone_code": "unknown-zone-code",
	})
	if zoneNo, err := parseZoneNoParameter(client, d); err == nil || zoneNo != nil {
		t.Fatalf("Unknown zone code must throw error. zone_no: %s", *zoneNo)
	}
}

func TestGetZoneNoByCode_basic(t *testing.T) {
	testZoneCode := "KR-2"
	client, _ := testApiClient(t)
	if zoneNo := getZoneNoByCode(client, testZoneCode); zoneNo == "" {
		t.Fatalf("No zone data for zone_code: %s", testZoneCode)
	}
}

func TestGetZoneNoByCode_inputUnknownZoneCode(t *testing.T) {
	testZoneCode := "unknown-zone-code"
	client, _ := testApiClient(t)
	if zoneNo := getZoneNoByCode(client, testZoneCode); zoneNo != "" {
		t.Fatalf("Unknown zone code must return nil. zone_code: %s", testZoneCode)
	}
}

func TestGetZoneByCode_basic(t *testing.T) {
	testZoneCode := "KR-2"
	client, _ := testApiClient(t)
	if zone, err := getZoneByCode(client, testZoneCode); err != nil || zone == nil {
		t.Fatalf("No zone data for zone_code: %s, %#v", testZoneCode, err)
	}
}

func TestGetZoneByCode_inputUnknownZoneCode(t *testing.T) {
	testZoneCode := "unknown-zone-code"
	client, _ := testApiClient(t)
	if zone, _ := getZoneByCode(client, testZoneCode); zone != nil {
		t.Fatalf("Unknown zone code must return nil. zone: %#v", zone)
	}
}

func TestGetZones_basic(t *testing.T) {
	client, _ := testApiClient(t)
	if zones, err := getZones(client); err != nil || zones == nil || len(zones) == 0 {
		t.Fatalf("No zone data")
	}
}
