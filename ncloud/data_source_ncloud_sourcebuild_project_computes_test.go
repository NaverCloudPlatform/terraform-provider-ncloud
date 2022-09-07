package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceBuildComputes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildComputesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcebuild_project_computes.computes"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceBuildComputesConfig() string {
	return `
data "ncloud_sourcebuild_project_computes" "computes" {}
`
}
