package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceZones_basic(t *testing.T) {
	//t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudZonesDataSourceID("data.ncloud_zones.zones"),
				),
			},
		},
	})
}

func testAccCheckNcloudZonesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Region data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("zone data source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudZonesConfig = `
data "ncloud_zones" "zones" {}
`
