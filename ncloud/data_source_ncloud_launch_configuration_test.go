package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudLaunchConfiguration_classic_basic(t *testing.T) {
	dataName := "data.ncloud_launch_configuration.lc"
	resourceName := "ncloud_launch_configuration.lc"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudLaunchConfigurationClassicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "launch_configuration_no", resourceName, "launch_configuration_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "server_image_product_code", resourceName, "server_image_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "server_product_code", resourceName, "server_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "member_server_image_no", resourceName, "member_server_image_no"),
					resource.TestCheckResourceAttrPair(dataName, "login_key_name", resourceName, "login_key_name"),
					resource.TestCheckResourceAttrPair(dataName, "user_data", resourceName, "user_data"),
					resource.TestCheckResourceAttrPair(dataName, "access_control_group_configuration_no_list", resourceName, "access_control_group_configuration_no_list"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudLaunchConfiguration_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_launch_configuration.lc"
	resourceName := "ncloud_launch_configuration.lc"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudLaunchConfigurationVpcConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "launch_configuration_no", resourceName, "launch_configuration_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "server_image_product_code", resourceName, "server_image_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "server_product_code", resourceName, "server_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "member_server_image_no", resourceName, "member_server_image_no"),
					resource.TestCheckResourceAttrPair(dataName, "login_key_name", resourceName, "login_key_name"),
					resource.TestCheckResourceAttrPair(dataName, "user_data", resourceName, "user_data"),
					resource.TestCheckResourceAttrPair(dataName, "access_control_group_configuration_no_list", resourceName, "access_control_group_configuration_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "is_encrypted_volume", resourceName, "is_encrypted_volume"),
					resource.TestCheckResourceAttrPair(dataName, "init_script_no", resourceName, "init_script_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudLaunchConfigurationClassicConfig() string {
	return fmt.Sprintf(`
resource "ncloud_launch_configuration" "lc" {
	server_image_product_code = "SPSW0LINUX000046"
}

data "ncloud_launch_configuration" "lc" {
	launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
}
`)
}

func testAccDataSourceNcloudLaunchConfigurationVpcConfig() string {
	return fmt.Sprintf(`
resource "ncloud_launch_configuration" "lc" {
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
}

data "ncloud_launch_configuration" "lc" {
	launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
}
`)
}
