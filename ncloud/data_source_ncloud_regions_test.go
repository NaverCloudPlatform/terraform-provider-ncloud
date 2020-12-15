package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudRegions_classic_basic(t *testing.T) {
	testAccDataSourceNcloudRegionsBasic(t, false)
}

func TestAccDataSourceNcloudRegions_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudRegionsBasic(t, true)
}

func testAccDataSourceNcloudRegionsBasic(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRegionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_regions.regions"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudRegionsConfig = `
data "ncloud_regions" "regions" {}
`
