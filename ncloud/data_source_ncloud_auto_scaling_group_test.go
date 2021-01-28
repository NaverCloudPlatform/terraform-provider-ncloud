package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudAutoScalingGroup_classic_basic(t *testing.T) {
	dataName := "data.ncloud_auto_scaling_group.auto"
	resourceName := "ncloud_auto_scaling_group.auto"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingGroupClassicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "auto_scaling_group_no", resourceName, "auto_scaling_group_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "launch_configuration_no", resourceName, "launch_configuration_no"),
					resource.TestCheckResourceAttrPair(dataName, "desired_capacity", resourceName, "desired_capacity"),
					resource.TestCheckResourceAttrPair(dataName, "min_size", resourceName, "min_size"),
					resource.TestCheckResourceAttrPair(dataName, "max_size", resourceName, "max_size"),
					resource.TestCheckResourceAttrPair(dataName, "default_cooldown", resourceName, "default_cooldown"),
					resource.TestCheckResourceAttrPair(dataName, "health_check_grace_period", resourceName, "health_check_grace_period"),
					resource.TestCheckResourceAttrPair(dataName, "health_check_type_code", resourceName, "health_check_type_code"),
					resource.TestCheckResourceAttrPair(dataName, "zone_no_list", resourceName, "zone_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "access_control_group_no_list", resourceName, "access_control_group_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "target_group_list", resourceName, "target_group_list"),
					resource.TestCheckResourceAttrPair(dataName, "server_name_prefix", resourceName, "server_name_prefix"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no_list", resourceName, "server_instance_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "wait_for_capacity_timeout", resourceName, "wait_for_capacity_timeout"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceNcloudAutoScalingGroup_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_auto_scaling_group.auto"
	resourceName := "ncloud_auto_scaling_group.auto"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingGroupVpcConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "auto_scaling_group_no", resourceName, "auto_scaling_group_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "launch_configuration_no", resourceName, "launch_configuration_no"),
					resource.TestCheckResourceAttrPair(dataName, "desired_capacity", resourceName, "desired_capacity"),
					resource.TestCheckResourceAttrPair(dataName, "min_size", resourceName, "min_size"),
					resource.TestCheckResourceAttrPair(dataName, "max_size", resourceName, "max_size"),
					resource.TestCheckResourceAttrPair(dataName, "default_cooldown", resourceName, "default_cooldown"),
					resource.TestCheckResourceAttrPair(dataName, "health_check_grace_period", resourceName, "health_check_grace_period"),
					resource.TestCheckResourceAttrPair(dataName, "health_check_type_code", resourceName, "health_check_type_code"),
					resource.TestCheckResourceAttrPair(dataName, "zone_no_list", resourceName, "zone_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "access_control_group_no_list", resourceName, "access_control_group_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "target_group_list", resourceName, "target_group_list"),
					resource.TestCheckResourceAttrPair(dataName, "server_name_prefix", resourceName, "server_name_prefix"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no_list", resourceName, "server_instance_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "wait_for_capacity_timeout", resourceName, "wait_for_capacity_timeout"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceNcloudAutoScalingGroupClassicConfig() string {
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

	data "ncloud_auto_scaling_group" "auto" {
		auto_scaling_group_no = ncloud_auto_scaling_group.auto.auto_scaling_group_no
	}
`)
}

func testAccDataSourceNcloudAutoScalingGroupVpcConfig() string {
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

data "ncloud_auto_scaling_group" "auto" {
	auto_scaling_group_no = ncloud_auto_scaling_group.auto.auto_scaling_group_no
}
`)
}
