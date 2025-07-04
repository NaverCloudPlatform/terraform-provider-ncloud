package autoscaling_test

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/autoscaling"
)

func TestAccResourceNcloudAutoScalingGroup_vpc_basic(t *testing.T) {
	var autoScalingGroup autoscaling.AutoScalingGroup
	resourceName := "ncloud_auto_scaling_group.auto"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckAutoScalingGroupDestroy(state, TestAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingGroupVpcConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, TestAccProvider),
					resource.TestCheckResourceAttr(resourceName, "default_cooldown", "300"),
					resource.TestCheckResourceAttr(resourceName, "health_check_grace_period", "300"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"wait_for_capacity_timeout",
					"access_control_group_no_list",
					"subnet_no",
				},
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingGroup_vpc_zero_value(t *testing.T) {
	var autoScalingGroup autoscaling.AutoScalingGroup
	resourceName := "ncloud_auto_scaling_group.auto"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckAutoScalingGroupDestroy(state, TestAccProvider)
		},
		Steps: []resource.TestStep{
			{
				//default
				Config: testAccAutoScalingGroupVpcConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, TestAccProvider),
					resource.TestCheckResourceAttr(resourceName, "default_cooldown", "300"),
					resource.TestCheckResourceAttr(resourceName, "health_check_grace_period", "300"),
				),
			},
			{
				//zero-value
				Config: testAccAutoScalingGroupVpcConfigWhenSetZero(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, TestAccProvider),
					resource.TestCheckResourceAttr(resourceName, "default_cooldown", "0"),
					resource.TestCheckResourceAttr(resourceName, "health_check_grace_period", "0"),
				),
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingGroup_vpc_disappears(t *testing.T) {
	var autoScalingGroup autoscaling.AutoScalingGroup
	resourceName := "ncloud_auto_scaling_group.auto"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckAutoScalingGroupDestroy(state, TestAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingGroupVpcConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, TestAccProvider),
					TestAccCheckResourceDisappears(TestAccProvider, autoscaling.ResourceNcloudAutoScalingGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAutoScalingGroupDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_auto_scaling_group" {
			continue
		}
		autoScalingGroup, err := autoscaling.GetAutoScalingGroup(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		if autoScalingGroup != nil {
			return fmt.Errorf("AutoScalingGroup(%s) still exists", ncloud.StringValue(autoScalingGroup.AutoScalingGroupNo))
		}
	}
	return nil
}

func testAccCheckAutoScalingGroupExists(n string, a *autoscaling.AutoScalingGroup, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no nutoScalingGroup ID is set: %s", n)
		}

		config := provider.Meta().(*conn.ProviderConfig)
		autoScalingGroup, err := autoscaling.GetAutoScalingGroup(config, rs.Primary.ID)
		if err != nil {
			return err
		}
		if autoScalingGroup == nil {
			return fmt.Errorf("not found AutoScalingGroup : %s", rs.Primary.ID)
		}
		*a = *autoScalingGroup
		return nil
	}
}

func testAccAutoScalingGroupVpcConfig() string {
	return `
resource "ncloud_vpc" "test" {
	ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_launch_configuration" "lc" {
	server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
}

resource "ncloud_auto_scaling_group" "auto" {
	access_control_group_no_list = [ncloud_vpc.test.default_access_control_group_no]
	subnet_no = ncloud_subnet.test.subnet_no
	launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
	min_size = 1
	max_size = 1
}
`
}

func testAccAutoScalingGroupVpcConfigWhenSetZero() string {
	return `
resource "ncloud_vpc" "test" {
	ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_launch_configuration" "lc" {
	server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
}

resource "ncloud_auto_scaling_group" "auto" {
	access_control_group_no_list = [ncloud_vpc.test.default_access_control_group_no]
	subnet_no = ncloud_subnet.test.subnet_no
	launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
	min_size = 1
	max_size = 1
	default_cooldown = 0
	health_check_grace_period = 0
}
`
}
