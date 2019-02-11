package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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

func TestParseZoneNoParameterBasic(t *testing.T) {
	testZoneCode := "KR-2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					s := testZoneSchema()
					d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
						"zone_code": testZoneCode,
					})
					if zoneNo, _ := parseZoneNoParameter(client, d); zoneNo == nil {
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
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					if zoneNo, _ := parseZoneNoParameter(client, &schema.ResourceData{}); zoneNo != nil {
						t.Fatalf("zone_no should be return nil when input empty resource data. actual: %s", *zoneNo)
					}
					return nil
				},
			},
		},
	})
}

func TestParseZoneNoParameterInputZoneNo(t *testing.T) {
	testZoneNo := "1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					s := testZoneSchema()
					d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
						"zone_no": testZoneNo,
					})
					if zoneNo, _ := parseZoneNoParameter(client, d); *zoneNo != testZoneNo {
						t.Fatalf("Expected: %s, Actual: %s", testZoneNo, *zoneNo)
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
					client := testAccProvider.Meta().(*NcloudAPIClient)
					s := testZoneSchema()
					d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
						"zone_code": "unknown-zone-code",
					})
					if zoneNo, err := parseZoneNoParameter(client, d); err == nil || zoneNo != nil {
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
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					if zoneNo := getZoneNoByCode(client, testZoneCode); zoneNo == "" {
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
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					if zoneNo := getZoneNoByCode(client, testZoneCode); zoneNo != "" {
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
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					if zone, err := getZoneByCode(client, testZoneCode); err != nil || zone == nil {
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
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					if zone, _ := getZoneByCode(client, testZoneCode); zone != nil {
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
				Config: `data "ncloud_zones" "zones" {}`,
				Check: func(*terraform.State) error {
					client := testAccProvider.Meta().(*NcloudAPIClient)
					if zones, err := getZones(client); err != nil || zones == nil || len(zones) == 0 {
						t.Fatalf("No zone data")
					}
					return nil
				},
			},
		},
	})
}
