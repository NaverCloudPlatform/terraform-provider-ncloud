package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccResourceNcloudLbTargetGroupAttachment_basic(t *testing.T) {
	var target vloadbalancer.Target
	targetGroupName := fmt.Sprintf("terraform-testacc-tga-%s", acctest.RandString(5))
	testServerName := getTestServerName()
	resourceName := "ncloud_lb_target_group_attachment.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLbTargetGroupAttachmentDestroy(state, testAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudLbTargetGroupAttachmentConfig(targetGroupName, testServerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLbTargetGroupAttachmentExists(resourceName, &target, testAccProvider),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_no"),
					resource.TestCheckResourceAttrSet(resourceName, "target_no"),
				),
			},
		},
	})
}

func testAccCheckLbTargetGroupAttachmentExists(n string, t *vloadbalancer.Target, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Target ID is set: %s", n)
		}

		config := provider.Meta().(*ProviderConfig)
		resp, err := config.Client.vloadbalancer.V2Api.GetTargetList(&vloadbalancer.GetTargetListRequest{
			RegionCode:    &config.RegionCode,
			TargetGroupNo: ncloud.String(rs.Primary.Attributes["target_group_no"]),
		})
		if err != nil {
			return err
		}

		if len(resp.TargetList) < 1 {
			return fmt.Errorf("Not found Target : %s", rs.Primary.ID)
		}

		*t = *resp.TargetList[0]
		return nil
	}
}

func testAccCheckLbTargetGroupAttachmentDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_lb_target_group_attachment" {
			continue
		}

		resp, err := config.Client.vloadbalancer.V2Api.GetTargetList(&vloadbalancer.GetTargetListRequest{
			RegionCode:    &config.RegionCode,
			TargetGroupNo: ncloud.String(rs.Primary.Attributes["target_group_no"]),
		})
		if err != nil {
			return err
		}

		if len(resp.TargetList) > 0 {
			return fmt.Errorf("Target (%s) still exists in Target Group (%s)", ncloud.StringValue(resp.TargetList[0].TargetNo), rs.Primary.Attributes["target_group_no"])
		}
	}
	return nil
}

func testAccResourceNcloudLbTargetGroupAttachmentConfig(targetGroupName string, serverName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_login_key" "test" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "test" {
	subnet_no = ncloud_subnet.test.subnet_no
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.test.key_name
}

resource "ncloud_lb_target_group" "test" {
  vpc_no   = ncloud_vpc.test.vpc_no
  protocol = "HTTP"
  target_type = "VSVR"
  port        = 8080
  name        = "%[2]s"
  description = "for test"

  health_check {
	protocol = "HTTP"
    http_method = "GET"
    port           = 8080
    url_path       = "/monitor/l7check"
    cycle          = 30
    up_threshold   = 2 
    down_threshold = 2 
  }

  algorithm_type = "RR"
  use_sticky_session = true
}

resource "ncloud_lb_target_group_attachment" "test" {
  target_group_no = ncloud_lb_target_group.test.target_group_no
  target_no = ncloud_server.test.instance_no
}

`, serverName, targetGroupName)
}
