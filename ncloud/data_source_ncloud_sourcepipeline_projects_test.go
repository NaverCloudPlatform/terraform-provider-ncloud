package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourcePipelineProjects_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineProjectsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcepipeline_projects.projects"),
				),
			},
		},
	})
}
func TestAccDataSourceNcloudSourcePipelineProjects_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineProjectsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcepipeline_projects.projects"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourcePipelineProjectsConfig() string {
	return fmt.Sprintf(`
data "ncloud_sourcepipeline_projects" "projects" {
}
`)
}
