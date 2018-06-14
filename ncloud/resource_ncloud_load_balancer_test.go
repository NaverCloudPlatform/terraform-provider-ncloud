package ncloud

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNcloudLoadBalancerBasic(t *testing.T) {
	var loadBalancerInstance sdk.LoadBalancerInstance
	prefix := getTestPrefix()
	testLoadBalancerName := prefix + "_lb"

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if loadBalancerInstance.LoadBalancerName != testLoadBalancerName {
				return fmt.Errorf("not found: %s", testLoadBalancerName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_load_balancer.lb",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig(testLoadBalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists("ncloud_load_balancer.lb", &loadBalancerInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer.lb",
						"load_balancer_name",
						testLoadBalancerName),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer.lb",
						"load_balancer_algorithm_type_code",
						"SIPHS"),
				),
			},
		},
	})
}

func TestAccNcloudLoadBalancerChangeConfiguration(t *testing.T) {
	var before sdk.LoadBalancerInstance
	var after sdk.LoadBalancerInstance
	prefix := getTestPrefix()
	testLoadBalancerName := prefix + "_lb"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_load_balancer.lb",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig(testLoadBalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists("ncloud_load_balancer.lb", &before),
				),
			},
			{
				Config: testAccLoadBalancerChangedConfig(testLoadBalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists("ncloud_load_balancer.lb", &after),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer.lb",
						"load_balancer_description",
						"tftest_lb change port")),
			},
		},
	})
}

func testAccCheckLoadBalancerExists(n string, i *sdk.LoadBalancerInstance) resource.TestCheckFunc {
	return testAccCheckLoadBalancerExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckLoadBalancerExistsWithProvider(n string, i *sdk.LoadBalancerInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		conn := provider.Meta().(*NcloudSdk).conn
		LoadBalancerInstance, err := getLoadBalancerInstance(conn, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if LoadBalancerInstance != nil {
			*i = *LoadBalancerInstance
			return nil
		}

		return fmt.Errorf("load balancer instance not found")
	}
}

func testAccCheckLoadBalancerDestroy(s *terraform.State) error {
	return testAccCheckLoadBalancerDestroyWithProvider(s, testAccProvider)
}

func testAccCheckLoadBalancerDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*NcloudSdk).conn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_load_balancer" {
			continue
		}
		loadBalancerInstance, err := getLoadBalancerInstance(conn, rs.Primary.ID)
		if loadBalancerInstance == nil {
			return nil
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("failed to delete load balancer: %s", loadBalancerInstance.LoadBalancerName)
	}

	return nil
}

func testAccLoadBalancerConfig(lbName string) string {
	return fmt.Sprintf(`
		resource "ncloud_load_balancer" "lb" {
			"load_balancer_name"                = "%s"
			"load_balancer_algorithm_type_code" = "SIPHS"
			"load_balancer_description"         = "tftest_lb description"

			"load_balancer_rule_list" = [
				{
					"protocol_type_code"   = "HTTP"
					"load_balancer_port"   = 80
					"server_port"          = 80
					"l7_health_check_path" = "/monitor/l7check"
				},
			]

			"internet_line_type_code" = "PUBLC"
			"network_usage_type_code" = "PBLIP"
			"region_no"               = "1"
		}
		`, lbName)
}

func testAccLoadBalancerChangedConfig(lbName string) string {
	return fmt.Sprintf(`
		resource "ncloud_load_balancer" "lb" {
			"load_balancer_name"                = "%s"
			"load_balancer_algorithm_type_code" = "SIPHS"
			"load_balancer_description"         = "tftest_lb change port"

			"load_balancer_rule_list" = [
				{
					"protocol_type_code"   = "HTTP"
					"load_balancer_port"   = 8080
					"server_port"          = 8080
					"l7_health_check_path" = "/monitor/l7check"
				},
			]

			"internet_line_type_code" = "PUBLC"
			"network_usage_type_code" = "PBLIP"
			"region_no"               = "1"
		}
		`, lbName)
}
