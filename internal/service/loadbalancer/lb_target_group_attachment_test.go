package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/loadbalancer"
)

func TestAccResourceNcloudLbTargetGroupAttachment_basic(t *testing.T) {
	t.Skip()
	{
		// Skip: deprecated server_image_product_code
	}

	var target string
	targetGroupName := fmt.Sprintf("terraform-testacc-tga-%s", acctest.RandString(5))
	testServerName := GetTestServerName()
	resourceName := "ncloud_lb_target_group_attachment.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLbTargetGroupAttachmentDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudLbTargetGroupAttachmentConfig(targetGroupName, testServerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLbTargetGroupAttachmentExists(resourceName, &target, GetTestProvider(true)),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_no"),
					resource.TestCheckResourceAttr(resourceName, "target_no_list.#", "1"),
				),
			},
		},
	})
}

func testAccCheckLbTargetGroupAttachmentExists(n string, t *string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Target ID is set: %s", n)
		}

		config := provider.Meta().(*conn.ProviderConfig)
		targetNoList, err := loadbalancer.GetVpcLoadBalancerTargetGroupAttachment(config, rs.Primary.Attributes["target_group_no"], []string{rs.Primary.Attributes["target_no_list.0"]})

		if err != nil {
			return err
		}

		if targetNoList == nil {
			return fmt.Errorf("Not found Target : %s, %s", rs.Primary.ID, rs.Primary.Attributes["target_no_list.1"])
		}

		t = ncloud.String(targetNoList[0])
		return nil
	}
}

func testAccCheckLbTargetGroupAttachmentDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_lb_target_group_attachment" {
			continue
		}

		targetNoList, err := loadbalancer.GetVpcLoadBalancerTargetGroupAttachment(config, rs.Primary.Attributes["target_group_no"], []string{rs.Primary.Attributes["target_no_list.0"]})

		if err != nil {
			return err
		}

		if targetNoList != nil {
			return fmt.Errorf("Target (%s) still exists in Target Group (%s)", targetNoList, rs.Primary.Attributes["target_group_no.0"])
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
  target_no_list = [ncloud_server.test.instance_no]
}

`, serverName, targetGroupName)
}
