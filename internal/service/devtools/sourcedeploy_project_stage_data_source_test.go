package devtools_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceDeployStage(t *testing.T) {
	stageNameSvr := getTestSourceDeployStageName() + "svr"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployStageConfig(stageNameSvr),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcedeploy_project_stage.stage"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployStageConfig(stageNameSvr string) string {
	return fmt.Sprintf(`
data "ncloud_server" "server" {
	filter {
		name   = "name"
		values = ["%[1]s"]
	}
}

resource "ncloud_sourcedeploy_project" "sd_project" {
	name = "tf-test-project2"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  = ncloud_sourcedeploy_project.sd_project.id
	name        = "%[2]s"
	target_type = "Server"
	config {
		server{
			id = data.ncloud_server.server.id
		} 
	}
}

data "ncloud_sourcedeploy_project_stage" "stage"{
	project_id = ncloud_sourcedeploy_project.sd_project.id
	id         = ncloud_sourcedeploy_project_stage.svr_stage.id
}
`, TF_TEST_SD_SERVER_NAME, stageNameSvr)
}
