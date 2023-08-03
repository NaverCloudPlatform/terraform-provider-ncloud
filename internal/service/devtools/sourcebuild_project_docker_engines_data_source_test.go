package devtools_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceBuildDockerEngines(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildDockerEnginesConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcebuild_project_docker_engines.docker_engines"),
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
