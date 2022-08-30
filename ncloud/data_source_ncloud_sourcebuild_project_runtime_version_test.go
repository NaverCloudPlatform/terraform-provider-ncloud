package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceBuildRuntimeVersion(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildRuntimeVersionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcebuild_project_runtime_version.runtime_version"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceBuildRuntimeVersionConfig() string {
	return `
data "ncloud_sourcebuild_project_os" "os" {
}

data "ncloud_sourcebuild_project_runtime" "runtime" {
	os_id = data.ncloud_sourcebuild_project_os.os.os[0].id
}

data "ncloud_sourcebuild_project_runtime_version" "runtime_version" {
	os_id      = data.ncloud_sourcebuild_project_os.os.os[0].id
	runtime_id = data.ncloud_sourcebuild_project_runtime.runtime.runtime[0].id
}
`
}
