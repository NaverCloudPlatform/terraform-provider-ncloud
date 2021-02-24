package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudLb_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-testacc-lb-%s", acctest.RandString(5))
	dataName := "data.ncloud_lb.lb_test"
	resourceName := "ncloud_lb.foo"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudLbConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "network_type", resourceName, "network_type"),
					resource.TestCheckResourceAttrPair(dataName, "idle_timeout", resourceName, "idle_timeout"),
					resource.TestCheckResourceAttrPair(dataName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataName, "throughput_type", resourceName, "throughput_type"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list", resourceName, "subnet_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "listener_list", resourceName, "listener_list"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceNcloudLbConfig(name string) string {
	return testAccResourceNcloudLbConfig(name) + fmt.Sprintf(`
data "ncloud_lb" "lb_test" {
	id = ncloud_lb.foo.load_balancer_no
}
`)
}
