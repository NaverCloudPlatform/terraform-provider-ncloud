package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testZoneSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"zone": {
			Type: schema.TypeString,
		},
	}
	return s
}

func TestParseZoneNoParameterBasic(t *testing.T) {
	testZoneCode := "KR-2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: zonesConfig,
				Check: func(*terraform.State) error {
					config := testAccClassicProvider.Meta().(*ProviderConfig)
					s := testZoneSchema()
					d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
						"zone": testZoneCode,
					})
					if zoneNo, _ := parseZoneNoParameter(config, d); zoneNo == nil {
						t.Fatalf("zone_no should be returned when input zoneCode. input: %s", testZoneCode)
					}
					return nil
				},
			},
		},
	})
}

func TestParseZoneNoParameterInputNil(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: zonesConfig,
				Check: func(*terraform.State) error {
					config := testAccProvider.Meta().(*ProviderConfig)
					if zoneNo, _ := parseZoneNoParameter(config, &schema.ResourceData{}); zoneNo != nil {
						t.Fatalf("zone_no should be return nil when input empty resource data. actual: %s", *zoneNo)
					}
					return nil
				},
			},
		},
	})
}

func TestParseZoneNoParameterInputUnknownZoneCode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					config := testAccProvider.Meta().(*ProviderConfig)
					s := testZoneSchema()
					d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
						"zone": "unknown-zone-code",
					})
					if zoneNo, err := parseZoneNoParameter(config, d); err == nil || zoneNo != nil {
						t.Fatalf("Unknown zone code must throw error. zone_no: %s", *zoneNo)
					}
					return nil
				},
			},
		},
	})
}

func TestGetZoneNoByCodeBasic(t *testing.T) {
	testZoneCode := "KR-2"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: zonesConfig,
				Check: func(*terraform.State) error {
					config := testAccClassicProvider.Meta().(*ProviderConfig)
					if zoneNo := getZoneNoByCode(config, testZoneCode); zoneNo == "" {
						t.Fatalf("No zone data for zone_code: %s", testZoneCode)
					}
					return nil
				},
			},
		},
	})
}

func TestGetZoneNoByCodeInputUnknownZoneCode(t *testing.T) {
	testZoneCode := "unknown-zone-code"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: zonesConfig,
				Check: func(*terraform.State) error {
					config := testAccClassicProvider.Meta().(*ProviderConfig)
					if zoneNo := getZoneNoByCode(config, testZoneCode); zoneNo != "" {
						t.Fatalf("Unknown zone code must return nil. zone_code: %s", testZoneCode)
					}
					return nil
				},
			},
		},
	})
}

func TestGetZoneByCodeBasic(t *testing.T) {
	testZoneCode := "KR-2"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: zonesConfig,
				Check: func(*terraform.State) error {
					config := testAccProvider.Meta().(*ProviderConfig)
					if zone, err := getZoneByCode(config, testZoneCode); err != nil || zone == nil {
						t.Fatalf("No zone data for zone_code: %s, %#v", testZoneCode, err)
					}
					return nil
				},
			},
		},
	})
}

func TestGetZoneByCodeInputUnknownZoneCode(t *testing.T) {
	testZoneCode := "unknown-zone-code"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: zonesConfig,
				Check: func(*terraform.State) error {
					config := testAccProvider.Meta().(*ProviderConfig)
					if zone, _ := getZoneByCode(config, testZoneCode); zone != nil {
						t.Fatalf("Unknown zone code must return nil. zone: %#v", zone)
					}
					return nil
				},
			},
		},
	})
}

func TestGetZonesBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: zonesConfig,
				Check: func(*terraform.State) error {
					config := testAccProvider.Meta().(*ProviderConfig)
					if zones, err := getZones(config); err != nil || zones == nil || len(zones) == 0 {
						t.Fatalf("No zone data")
					}
					return nil
				},
			},
		},
	})
}

var zonesConfig = `data "ncloud_zones" "zones" {}`
