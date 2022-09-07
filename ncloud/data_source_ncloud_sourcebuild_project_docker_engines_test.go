package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceBuildDockerEngines(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildDockerEnginesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcebuild_project_docker_engines.docker_engines"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceBuildDockerEnginesConfig() string {
	return `
data "ncloud_sourcebuild_project_docker_engines" "docker_engines" {}
`
}
