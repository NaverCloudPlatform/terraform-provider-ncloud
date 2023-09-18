package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/loadbalancer"
)

func TestAccResourceNcloudLbListener_vpc_basic(t *testing.T) {
	var listener loadbalancer.LoadBalancerListener
	lbName := fmt.Sprintf("terraform-testacc-lb-%s", acctest.RandString(5))
	resourceName := "ncloud_lb_listener.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLbListenerDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudLbListenerConfig(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLbListenerExists(resourceName, &listener, GetTestProvider(true)),
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

func testAccCheckLbListenerExists(n string, l *loadbalancer.LoadBalancerListener, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LB Listener ID is set: %s", n)
		}

		config := provider.Meta().(*conn.ProviderConfig)
		listener, err := loadbalancer.GetVpcLoadBalancerListener(config, rs.Primary.ID, rs.Primary.Attributes["load_balancer_no"])
		if err != nil {
			return err
		}

		if listener == nil {
			return fmt.Errorf("Not found LB Listener : %s", rs.Primary.ID)
		}

		l = listener
		return nil
	}
}

func testAccCheckLbListenerDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_lb_listener" {
			continue
		}

		listener, err := loadbalancer.GetVpcLoadBalancerListener(config, rs.Primary.ID, rs.Primary.Attributes["load_balancer_no"])
		if err != nil {
			return err
		}

		if listener != nil {
			return fmt.Errorf("LB Listener(%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccResourceNcloudLbListenerConfig(lbName string) string {
	return testAccResourceNcloudLbConfig(lbName) + `
resource "ncloud_lb_listener" "test" {
    load_balancer_no = ncloud_lb.test.load_balancer_no
    protocol = "HTTP"
    port = 8080
    target_group_no = ncloud_lb_target_group.test.target_group_no
}
`
}
