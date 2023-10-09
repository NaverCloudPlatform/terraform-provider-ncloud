package devtools_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceDeployStages(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployStagesConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcedeploy_project_stages.stages"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployStagesConfig() string {
	return `
resource "ncloud_sourcedeploy_project" "sd_project" {
	name = "tf-test-project"
}

data "ncloud_sourcedeploy_project_stages" stages{
	project_id = ncloud_sourcedeploy_project.sd_project.id
}
`
}
