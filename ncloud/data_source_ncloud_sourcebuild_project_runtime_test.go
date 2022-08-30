package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceBuildRuntime(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildRuntimeConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcebuild_project_runtime.runtime"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceBuildRuntimeConfig() string {
	return `
data "ncloud_sourcebuild_project_os" "os" {
}

data "ncloud_sourcebuild_project_runtime" "runtime" {
	os_id = data.ncloud_sourcebuild_project_os.os.os[0].id
}
`
}
