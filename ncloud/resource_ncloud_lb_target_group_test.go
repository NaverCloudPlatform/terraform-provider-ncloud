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

func TestAccResourceNcloudLbTargetGroup_basic(t *testing.T) {
	var tg vloadbalancer.TargetGroup
	name := fmt.Sprintf("terraform-testacc-tg-%s", acctest.RandString(5))
	resourceName := "ncloud_lb_target_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLbTargetGroupDestroy(state, testAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudLbTargetGroupConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLbTargetGroupExists(resourceName, &tg, testAccProvider),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "target_type", "VSVR"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "for test"),
					resource.TestCheckResourceAttr(resourceName, "health_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.url_path", "/monitor/l7check"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.cycle", "30"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.up_threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.down_threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "algorithm_type", "RR"),
					resource.TestCheckResourceAttr(resourceName, "use_sticky_session", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_no"),
				),
			},
		},
	})
}

func testAccCheckLbTargetGroupExists(n string, t *vloadbalancer.TargetGroup, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Target Group ID is set: %s", n)
		}

		config := provider.Meta().(*ProviderConfig)
		resp, err := config.Client.vloadbalancer.V2Api.GetTargetGroupList(&vloadbalancer.GetTargetGroupListRequest{
			RegionCode:        &config.RegionCode,
			TargetGroupNoList: []*string{ncloud.String(rs.Primary.ID)},
		})
		if err != nil {
			return err
		}

		if len(resp.TargetGroupList) < 1 {
			return fmt.Errorf("Not found Target Group : %s", rs.Primary.ID)
		}

		*t = *resp.TargetGroupList[0]
		return nil
	}
}

func testAccCheckLbTargetGroupDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_lb_target_group" {
			continue
		}

		resp, err := config.Client.vloadbalancer.V2Api.GetTargetGroupList(&vloadbalancer.GetTargetGroupListRequest{
			RegionCode:        &config.RegionCode,
			TargetGroupNoList: []*string{ncloud.String(rs.Primary.ID)},
		})
		if err != nil {
			return err
		}

		if len(resp.TargetGroupList) > 0 {
			return fmt.Errorf("Target Group(%s) still exists", ncloud.StringValue(resp.TargetGroupList[0].TargetGroupNo))
		}
	}
	return nil
}

func testAccResourceNcloudLbTargetGroupConfig(name string) string {
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
	usage_type         = "LOADB"
}

resource "ncloud_lb_target_group" "test" {
  vpc_no   = ncloud_vpc.test.vpc_no
  protocol = "HTTP"
  target_type = "VSVR"
  port        = 8080
  name        = "%s"
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
`, name)
}
