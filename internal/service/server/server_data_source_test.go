package server_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServer_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_server.by_id"
	resourceName := "ncloud_server.server"
	testServerName := GetTestServerName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerVpcConfig(testServerName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "server_image_product_code", resourceName, "server_image_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "server_product_code", resourceName, "server_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "base_block_storage_disk_type", resourceName, "base_block_storage_disk_type"),
					resource.TestCheckResourceAttrPair(dataName, "cpu_count", resourceName, "cpu_count"),
					resource.TestCheckResourceAttrPair(dataName, "memory_size", resourceName, "memory_size"),
					resource.TestCheckResourceAttrPair(dataName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "platform_type", resourceName, "platform_type"),
					resource.TestCheckResourceAttrPair(dataName, "is_protect_server_termination", resourceName, "is_protect_server_termination"),
					resource.TestCheckResourceAttrPair(dataName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),

					// VPC only
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "network_interface.#", resourceName, "network_interface.#"),
					resource.TestCheckResourceAttrPair(dataName, "network_interface.0.network_interface_no", resourceName, "network_interface.0.network_interface_no"),
					TestAccCheckDataSourceID("data.ncloud_server.by_filter"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServer_classic_basic(t *testing.T) {
	dataName := "data.ncloud_server.by_id"
	resourceName := "ncloud_server.server"
	testServerName := GetTestServerName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerClassicConfig(testServerName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "server_image_product_code", resourceName, "server_image_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "server_product_code", resourceName, "server_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "base_block_storage_disk_type", resourceName, "base_block_storage_disk_type"),
					resource.TestCheckResourceAttrPair(dataName, "cpu_count", resourceName, "cpu_count"),
					resource.TestCheckResourceAttrPair(dataName, "memory_size", resourceName, "memory_size"),
					resource.TestCheckResourceAttrPair(dataName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "platform_type", resourceName, "platform_type"),
					resource.TestCheckResourceAttrPair(dataName, "is_protect_server_termination", resourceName, "is_protect_server_termination"),
					resource.TestCheckResourceAttrPair(dataName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),
					TestAccCheckDataSourceID("data.ncloud_server.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceServerVpcConfig(testServerName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

data "ncloud_server_image_numbers" "server_images" {
    filter {
        name = "name"
        values = ["ubuntu-22.04-base"]
  }
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
	server_spec_code = "s2-g3"
	login_key_name = ncloud_login_key.loginkey.key_name
}

data "ncloud_server" "by_id" {
	id = ncloud_server.server.id
}

data "ncloud_server" "by_filter" {
	filter {
		name = "instance_no"
		values = [ncloud_server.server.id]
	}
}
`, testServerName)
}

func testAccDataSourceServerClassicConfig(testServerName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "server" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000046"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

data "ncloud_server" "by_id" {
	id = ncloud_server.server.id
}

data "ncloud_server" "by_filter" {
	filter {
		name = "instance_no"
		values = [ncloud_server.server.id]
	}
}
`, testServerName)
}
