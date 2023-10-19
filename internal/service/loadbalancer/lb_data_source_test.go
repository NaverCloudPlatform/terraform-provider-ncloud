package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudLb_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-testacc-lb-%s", acctest.RandString(5))
	dataName := "data.ncloud_lb.test"
	resourceName := "ncloud_lb.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudLbConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "network_type", resourceName, "network_type"),
					resource.TestCheckResourceAttrPair(dataName, "idle_timeout", resourceName, "idle_timeout"),
					resource.TestCheckResourceAttrPair(dataName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataName, "throughput_type", resourceName, "throughput_type"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list", resourceName, "subnet_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudLbConfig(name string) string {
	return testAccResourceNcloudLbConfig(name) + `
data "ncloud_lb" "test" {
	id = ncloud_lb.test.load_balancer_no
}
`
}
