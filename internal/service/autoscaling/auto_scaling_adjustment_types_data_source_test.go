package autoscaling_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudAutoScalingAdjustmentTypes_classic_basic(t *testing.T) {
	lcName := fmt.Sprintf("lc-%s", acctest.RandString(5))
	policyName := fmt.Sprintf("policy-%s", acctest.RandString(5))
	dataName := "data.ncloud_auto_scaling_adjustment_types.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingAdjustmentTypesClassicConfig(lcName, policyName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAutoScalingAdjustmentTypes_vpc_basic(t *testing.T) {
	lcName := fmt.Sprintf("lc-%s", acctest.RandString(5))
	policyName := fmt.Sprintf("policy-%s", acctest.RandString(5))
	dataName := "data.ncloud_auto_scaling_adjustment_types.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingAdjustmentTypesVpcConfig(lcName, policyName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAutoScalingAdjustmentTypes_classic_byFilterCode(t *testing.T) {
	dataName := "data.ncloud_auto_scaling_adjustment_types.by_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingAdjustmentTypesByFilterCodeConfig("EXACT"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAutoScalingAdjustmentTypes_vpc_byFilterCode(t *testing.T) {
	dataName := "data.ncloud_auto_scaling_adjustment_types.by_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingAdjustmentTypesByFilterCodeConfig("EXACT"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAutoScalingAdjustmentTypesClassicConfig(lcName string, policyName string) string {
	return fmt.Sprintf(`
	resource "ncloud_launch_configuration" "lc" {
		name = "%s"
		server_image_product_code = "SPSW0LINUX000046"
		server_product_code = "SPSVRSSD00000003"
	}

	resource "ncloud_auto_scaling_group" "asg" {
		launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
		min_size = 1
		max_size = 1
		zone_no_list = ["2"]
		wait_for_capacity_timeout = "0"
	}

	resource "ncloud_auto_scaling_policy" "policy" {
		name = "%s"
		adjustment_type_code = data.ncloud_auto_scaling_adjustment_types.test.types[0].code
		scaling_adjustment = 2
		auto_scaling_group_no = ncloud_auto_scaling_group.asg.auto_scaling_group_no
	}
	
	
	data "ncloud_auto_scaling_adjustment_types" "test" {
	
	}
	`, lcName, policyName)
}

func testAccDataSourceNcloudAutoScalingAdjustmentTypesVpcConfig(lcName string, policyName string) string {
	return fmt.Sprintf(`
	resource "ncloud_launch_configuration" "lc" {
		name = "%s"
		server_image_product_code = "SPSW0LINUX000046"
		server_product_code = "SPSVRSSD00000003"
	}

	resource "ncloud_auto_scaling_group" "asg" {
		launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
		min_size = 1
		max_size = 1
		zone_no_list = ["2"]
		wait_for_capacity_timeout = "0"
	}

	resource "ncloud_auto_scaling_policy" "policy" {
		name = "%s"
		adjustment_type_code = data.ncloud_auto_scaling_adjustment_types.test.types[2].code 
		scaling_adjustment = 2
		auto_scaling_group_no = ncloud_auto_scaling_group.asg.auto_scaling_group_no
	}
	
	
	data "ncloud_auto_scaling_adjustment_types" "test" {
	}
	`, lcName, policyName)
}

func testAccDataSourceNcloudAutoScalingAdjustmentTypesByFilterCodeConfig(code string) string {
	return fmt.Sprintf(`
	data "ncloud_auto_scaling_adjustment_types" "by_filter" {
		filter {
			name   = "code"
			values = ["%s"]
		}
	}
`, code)
}
