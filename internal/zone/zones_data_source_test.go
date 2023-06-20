package zone

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudZones_classic_basic(t *testing.T) {
	testAccDataSourceNcloudZonesBasic(t, false)
}

func TestAccDataSourceNcloudZones_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudZonesBasic(t, true)
}

func testAccDataSourceNcloudZonesBasic(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_zones.zones"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudZonesConfig = `
data "ncloud_zones" "zones" {}
`
