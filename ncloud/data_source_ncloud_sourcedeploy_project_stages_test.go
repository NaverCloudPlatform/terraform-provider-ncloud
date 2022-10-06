package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceDeployStages(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployStagesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcedeploy_project_stages.stages"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployStagesConfig() string {
	return fmt.Sprintf(`
resource "ncloud_sourcedeploy_project" "sd_project" {
	name = "tf-test-project"
}

data "ncloud_sourcedeploy_project_stages" stages{
	project_id = ncloud_sourcedeploy_project.sd_project.id
}
`)
}
