package devtools_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourcePipelineProjects_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ClassicProtoV5ProviderFactories,
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
func TestAccDataSourceNcloudSourcePipelineProjects_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
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
	return fmt.Sprintf(`
data "ncloud_sourcepipeline_projects" "projects" {
}
`)
}
