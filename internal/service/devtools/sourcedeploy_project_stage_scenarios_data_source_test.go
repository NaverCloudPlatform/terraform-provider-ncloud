package devtools_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceDeployScenarios(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployScenariosConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcedeploy_project_stage_scenarios.scenarios"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployScenariosConfig() string {
	return fmt.Sprintf(`
data "ncloud_server" "server" {
	filter {
		name   = "name"
		values = ["%[1]s"]
	}
}

resource "ncloud_sourcedeploy_project" "sd_project" {
	name = "tf-test-project"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  = ncloud_sourcedeploy_project.sd_project.id
	name		= "svr"
	target_type = "Server"
	config {
		server {
			id = data.ncloud_server.server.id
		}
	}
}

data "ncloud_sourcedeploy_project_stage_scenarios" "scenarios"{
	project_id = ncloud_sourcedeploy_project.sd_project.id
	stage_id   = ncloud_sourcedeploy_project_stage.svr_stage.id
}
`, TF_TEST_SD_SERVER_NAME)
}
