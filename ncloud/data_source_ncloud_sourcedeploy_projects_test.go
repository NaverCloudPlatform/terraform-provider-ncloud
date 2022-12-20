package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceDeployProjects(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployProjectsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcedeploy_projects.projects"),
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
