package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudLbListener_basic(t *testing.T) {
	lbName := fmt.Sprintf("terraform-testacc-lb-%s", acctest.RandString(5))
	dataName := "data.ncloud_lb_listener.test"
	resourceName := "ncloud_lb_listener.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudLbListenerConfig(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataName, "protocol", resourceName, "protocol"),
					resource.TestCheckResourceAttrPair(dataName, "rule_no_list", resourceName, "rule_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "load_balancer_no", resourceName, "load_balancer_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudLbListenerConfig(name string) string {
	return testAccResourceNcloudLbListenerConfig(name) + fmt.Sprintf(`
data "ncloud_lb_listener" "test" {
	id = ncloud_lb_listener.test.listener_no
	load_balancer_no = ncloud_lb.test.load_balancer_no
}
`)
}
