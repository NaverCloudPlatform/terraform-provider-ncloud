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

func TestAccResourceNcloudLb_vpc_basic(t *testing.T) {
	var lb vloadbalancer.LoadBalancerInstance
	lbName := fmt.Sprintf("terraform-testacc-lb-%s", acctest.RandString(5))
	resourceName := "ncloud_lb.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLbDestroy(state, testAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudLbConfig(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLbExists(resourceName, &lb, testAccProvider),
					resource.TestCheckResourceAttr(resourceName, "name", lbName),
					resource.TestCheckResourceAttr(resourceName, "description", "tf test description"),
					resource.TestCheckResourceAttr(resourceName, "network_type", "PRIVATE"),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "type", "APPLICATION"),
					resource.TestCheckResourceAttr(resourceName, "throughput_type", "SMALL"),
					resource.TestCheckResourceAttr(resourceName, "subnet_no_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "operation", "NULL"),
					resource.TestCheckResourceAttr(resourceName, "status_code", "USED"),
					resource.TestCheckResourceAttr(resourceName, "status_name", "Running"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_no"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"listener", "description"},
			},
		},
	})
}

func testAccCheckLbExists(n string, lb *vloadbalancer.LoadBalancerInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LB ID is set: %s", n)
		}

		config := provider.Meta().(*ProviderConfig)
		resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(&vloadbalancer.GetLoadBalancerInstanceDetailRequest{
			RegionCode:             &config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		if len(resp.LoadBalancerInstanceList) < 1 {
			return fmt.Errorf("Not found LB : %s", rs.Primary.ID)
		}

		*lb = *resp.LoadBalancerInstanceList[0]
		return nil
	}
}

func testAccCheckLbDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_lb" {
			continue
		}

		resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(&vloadbalancer.GetLoadBalancerInstanceDetailRequest{
			RegionCode:             &config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		if len(resp.LoadBalancerInstanceList) > 0 {
			return fmt.Errorf("LB(%s) still exists", ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
		}
	}
	return nil
}

func testAccResourceNcloudLbConfig(name string) string {
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
  name        = "terraform-testacc-tg"
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

resource "ncloud_lb" "test" {
    name = "%s"
    description = "tf test description"
    network_type = "PRIVATE"
    idle_timeout = 30
    type = "APPLICATION"
    throughput_type = "SMALL"
    subnet_no_list = [ ncloud_subnet.test.subnet_no ]
}
`, name)
}
