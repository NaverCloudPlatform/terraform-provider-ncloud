package devtools_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourcePipelineProjects_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineProjectsConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcepipeline_projects.projects"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourcePipelineProjectsConfig() string {
	return `
data "ncloud_sourcepipeline_projects" "projects" {
}
`
}
