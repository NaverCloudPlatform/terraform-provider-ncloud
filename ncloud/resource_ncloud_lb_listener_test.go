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

func TestAccResourceNcloudLbListener_vpc_basic(t *testing.T) {
	var listener vloadbalancer.LoadBalancerListener
	lbName := fmt.Sprintf("terraform-testacc-lb-%s", acctest.RandString(5))
	resourceName := "ncloud_lb_listener.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLbListenerDestroy(state, testAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudLbListenerConfig(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLbListenerExists(resourceName, &listener, testAccProvider),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "rule_no_list.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "target_group_no"),
					resource.TestCheckResourceAttrSet(resourceName, "load_balancer_no"),
				),
			},
		},
	})
}

func testAccCheckLbListenerExists(n string, l *vloadbalancer.LoadBalancerListener, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LB Listener ID is set: %s", n)
		}

		config := provider.Meta().(*ProviderConfig)
		resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerListenerList(&vloadbalancer.GetLoadBalancerListenerListRequest{
			RegionCode:             &config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(rs.Primary.Attributes["load_balancer_no"]),
		})

		if err != nil {
			return err
		}

		exist := false
		for _, l := range resp.LoadBalancerListenerList {
			if rs.Primary.ID == *l.LoadBalancerListenerNo {
				exist = true
				break
			}
		}

		if !exist {
			return fmt.Errorf("Not found LB Listener : %s", rs.Primary.ID)
		}

		*l = *resp.LoadBalancerListenerList[0]
		return nil
	}
}

func testAccCheckLbListenerDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_lb_listener" {
			continue
		}

		resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerListenerList(&vloadbalancer.GetLoadBalancerListenerListRequest{
			RegionCode:             &config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(rs.Primary.Attributes["load_balancer_no"]),
		})

		if err != nil {
			return err
		}

		exist := false
		for _, l := range resp.LoadBalancerListenerList {
			if rs.Primary.ID == *l.LoadBalancerListenerNo {
				exist = true
				break
			}
		}

		if exist {
			return fmt.Errorf("LB Listener(%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccResourceNcloudLbListenerConfig(lbName string) string {
	return testAccResourceNcloudLbConfig(lbName) + fmt.Sprintf(`
resource "ncloud_lb_listener" "test" {
    load_balancer_no = ncloud_lb.test.load_balancer_no
    protocol = "HTTP"
    port = 8080
    target_group_no = ncloud_lb_target_group.test.target_group_no
}
`)
}
