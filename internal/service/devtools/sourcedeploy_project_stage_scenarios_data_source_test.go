package devtools_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceDeployScenarios(t *testing.T) {
	stageNameSvr := GetTestSourceDeployScenarioName() + "svr"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployScenariosConfig(stageNameSvr, productCode),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcedeploy_project_stage_scenarios.scenarios"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployScenariosConfig(stageNameSvr string, productCode string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[2]s-key"
}

resource "ncloud_vpc" "vpc" {
	ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}


resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.subnet.id
	name = "%[2]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
	server_product_code = "%[3]s"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_sourcedeploy_project" "sd_project" {
	name = "%[2]s"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  = ncloud_sourcedeploy_project.sd_project.id
	name		= "%[2]s"
	target_type = "Server"
	config {
		server {
			id = ncloud_server.server.id
		}
	}
}

data "ncloud_sourcedeploy_project_stage_scenarios" "scenarios"{
	project_id = ncloud_sourcedeploy_project.sd_project.id
	stage_id   = ncloud_sourcedeploy_project_stage.svr_stage.id
}
`, TF_TEST_SD_SERVER_NAME, stageNameSvr, productCode)
}
