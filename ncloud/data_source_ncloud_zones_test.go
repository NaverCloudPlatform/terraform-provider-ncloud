package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudZones_classic_basic(t *testing.T) {
	testAccDataSourceNcloudZonesBasic(t, false)
}

func TestAccDataSourceNcloudZones_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudZonesBasic(t, true)
}

func testAccDataSourceNcloudZonesBasic(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
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

var testAccDataSourceNcloudZonesConfig = `
data "ncloud_zones" "zones" {}
`
