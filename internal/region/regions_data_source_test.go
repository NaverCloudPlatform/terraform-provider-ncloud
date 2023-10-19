package region_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudRegions_classic_basic(t *testing.T) {
	testAccDataSourceNcloudRegionsBasic(t, false)
}

func TestAccDataSourceNcloudRegions_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudRegionsBasic(t, true)
}

func testAccDataSourceNcloudRegionsBasic(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRegionsConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_regions.regions"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudRegionsConfig = `
data "ncloud_regions" "regions" {}
`
