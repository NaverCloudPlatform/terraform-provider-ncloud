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

func TestAccResourceNcloudAutoScalingGroup_classic_basic(t *testing.T) {
	var autoScalingGroup autoscaling.AutoScalingGroup
	resourceName := "ncloud_auto_scaling_group.auto"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckAutoScalingGroupDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingGroupClassicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, GetTestProvider(false)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"wait_for_capacity_timeout",
					"zone_no_list",
				},
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingGroup_vpc_basic(t *testing.T) {
	var autoScalingGroup autoscaling.AutoScalingGroup
	resourceName := "ncloud_auto_scaling_group.auto"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckAutoScalingGroupDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingGroupVpcConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, GetTestProvider(true)),
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

func TestAccResourceNcloudAutoScalingGroup_classic_disappears(t *testing.T) {
	var autoScalingGroup autoscaling.AutoScalingGroup
	resourceName := "ncloud_auto_scaling_group.auto"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckAutoScalingGroupDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingGroupClassicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, GetTestProvider(false)),
					TestAccCheckResourceDisappears(GetTestProvider(false), autoscaling.ResourceNcloudAutoScalingGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingGroup_vpc_disappears(t *testing.T) {
	var autoScalingGroup autoscaling.AutoScalingGroup
	resourceName := "ncloud_auto_scaling_group.auto"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckAutoScalingGroupDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingGroupVpcConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingGroupExists(resourceName, &autoScalingGroup, GetTestProvider(true)),
					TestAccCheckResourceDisappears(GetTestProvider(true), autoscaling.ResourceNcloudAutoScalingGroup(), resourceName),
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
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AutoScalingGroup ID is set: %s", n)
		}

		config := provider.Meta().(*conn.ProviderConfig)
		autoScalingGroup, err := autoscaling.GetAutoScalingGroup(config, rs.Primary.ID)
		if err != nil {
			return err
		}
		if autoScalingGroup == nil {
			return fmt.Errorf("Not found AutoScalingGroup : %s", rs.Primary.ID)
		}
		*a = *autoScalingGroup
		return nil
	}
}

func testAccAutoScalingGroupClassicConfig() string {
	return fmt.Sprintf(`
resource "ncloud_launch_configuration" "lc" {
	server_image_product_code = "SPSW0LINUX000046"
}

resource "ncloud_auto_scaling_group" "auto" {
	launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
	min_size = 2
	max_size = 2
	zone_no_list = ["2"]
}
`)
}

func testAccAutoScalingGroupVpcConfig() string {
	return fmt.Sprintf(`
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
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
}

resource "ncloud_auto_scaling_group" "auto" {
	access_control_group_no_list = [ncloud_vpc.test.default_access_control_group_no]
	subnet_no = ncloud_subnet.test.subnet_no
	launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
	min_size = 1
	max_size = 1
}
`)
}
