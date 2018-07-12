package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
)

func TestAccDataSourceNcloudZonesBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_zones.zones"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudZonesByRegionCode(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudZonesByRegionCodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_zones.zones"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudZonesByInvalidRegionCode(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceNcloudZonesByInvalidRegionCodeConfig,
				ExpectError: regexp.MustCompile(`no region data for (.*) please change region_code and try again`),
			},
		},
	})
}

var testAccDataSourceNcloudZonesConfig = `
data "ncloud_zones" "zones" {}
`

var testAccDataSourceNcloudZonesByRegionCodeConfig = `
data "ncloud_zones" "zones" {
	"region_code" = "JPN"
}
`

var testAccDataSourceNcloudZonesByInvalidRegionCodeConfig = `
data "ncloud_zones" "zones" {
	"region_code" = "INVALID"
}
`
