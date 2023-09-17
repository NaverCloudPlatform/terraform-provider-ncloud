package devtools_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceDeployProjects(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployProjectsConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcedeploy_projects.projects"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployProjectsConfig() string {
	return fmt.Sprintf(`
data "ncloud_sourcedeploy_projects" "projects" {
}
`)
}
