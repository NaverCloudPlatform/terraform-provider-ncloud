package classicloadbalancer

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func TestAccNcloudLoadBalancerBasic(t *testing.T) {
	var loadBalancerInstance loadbalancer.LoadBalancerInstance
	prefix := GetTestPrefix()
	testLoadBalancerName := prefix + "_lb"

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if *loadBalancerInstance.LoadBalancerName != testLoadBalancerName {
				return fmt.Errorf("not found: %s", testLoadBalancerName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(false),
		CheckDestroy: testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig(testLoadBalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists("ncloud_load_balancer.lb", &loadBalancerInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer.lb",
						"name",
						testLoadBalancerName),
					resource.TestCheckResourceAttr(
						"ncloud_load_balancer.lb",
						"algorithm_type",
						"SIPHS"),
				),
			},
			{
				ResourceName:            "ncloud_load_balancer.lb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func TestAccNcloudLoadBalancerChangeConfiguration(t *testing.T) {
	var before loadbalancer.LoadBalancerInstance
	var after loadbalancer.LoadBalancerInstance
	prefix := GetTestPrefix()
	testLoadBalancerName := prefix + "_lb"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(false),
		CheckDestroy: testAccCheckLoadBalancerDestroy,
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
						"description",
						"tftest_lb change port")),
			},
			{
				ResourceName:            "ncloud_load_balancer.lb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccCheckLoadBalancerExists(n string, i *loadbalancer.LoadBalancerInstance) resource.TestCheckFunc {
	return testAccCheckLoadBalancerExistsWithProvider(n, i, func() *schema.Provider { return GetTestProvider(false) })
}

func testAccCheckLoadBalancerExistsWithProvider(n string, i *loadbalancer.LoadBalancerInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		client := provider.Meta().(*ProviderConfig).Client
		LoadBalancerInstance, err := getLoadBalancerInstance(client, rs.Primary.ID)
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
	return testAccCheckLoadBalancerDestroyWithProvider(s, GetTestProvider(false))
}

func testAccCheckLoadBalancerDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*ProviderConfig).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_load_balancer" {
			continue
		}
		loadBalancerInstance, err := getLoadBalancerInstance(client, rs.Primary.ID)
		if loadBalancerInstance == nil {
			return nil
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("failed to delete load balancer: %s", *loadBalancerInstance.LoadBalancerName)
	}

	return nil
}

func testAccLoadBalancerConfig(lbName string) string {
	return fmt.Sprintf(`
		resource "ncloud_load_balancer" "lb" {
			name           = "%s"
			algorithm_type = "SIPHS"
			description    = "tftest_lb description"

		rule_list {
    		protocol_type        = "HTTP"
    		load_balancer_port   = 80
    		server_port          = 80
    		l7_health_check_path = "/monitor/l7check"
  		}

			network_usage_type = "PBLIP"
			region             = "KR"
		}
		`, lbName)
}

func testAccLoadBalancerChangedConfig(lbName string) string {
	return fmt.Sprintf(`
		resource "ncloud_load_balancer" "lb" {
			name           = "%s"
			algorithm_type = "SIPHS"
			description    = "tftest_lb change port"

		rule_list {
    		protocol_type        = "HTTP"
    		load_balancer_port   = 8080
    		server_port          = 8080
    		l7_health_check_path = "/monitor/l7check"
  		}

			network_usage_type = "PBLIP"
			region             = "KR"
		}
		`, lbName)
}
