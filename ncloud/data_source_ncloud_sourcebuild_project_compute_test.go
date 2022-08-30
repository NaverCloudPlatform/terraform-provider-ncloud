package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceBuildCompute(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildComputeConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcebuild_project_compute.compute"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceBuildComputeConfig() string {
	return `
data "ncloud_sourcebuild_project_compute" "compute" {}
`
}
