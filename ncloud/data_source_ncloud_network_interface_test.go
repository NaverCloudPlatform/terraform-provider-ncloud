package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudNetworkInterfaceBasic(t *testing.T) {
	name := fmt.Sprintf("tf-ds-nic-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_interface.foo"
	dataName := "data.ncloud_network_interface.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkInterfaceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "network_interface_no", resourceName, "network_interface_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "private_ip", resourceName, "private_ip"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "access_control_groups", resourceName, "access_control_groups"),
					resource.TestCheckResourceAttrPair(dataName, "is_default", resourceName, "is_default"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNetworkInterfaceFilter(t *testing.T) {
	name := fmt.Sprintf("tf-ds-nic-filter-%s", acctest.RandString(5))
	resourceName := "ncloud_network_interface.foo"
	dataName := "data.ncloud_network_interface.by_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkInterfaceConfigFilter(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_interface.by_filter"),
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "network_interface_no", resourceName, "network_interface_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "private_ip", resourceName, "private_ip"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "access_control_groups", resourceName, "access_control_groups"),
					resource.TestCheckResourceAttrPair(dataName, "is_default", resourceName, "is_default"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNetworkInterfaceConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s"
	subnet             = "10.4.0.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_network_interface" "foo" {
	name                 = "%[1]s"
	description           = "for acc test"
	subnet_no             = ncloud_subnet.test.id
	access_control_groups = [ncloud_vpc.test.default_access_control_group_no]
}

data "ncloud_network_interface" "by_id" {
	network_interface_no = ncloud_network_interface.foo.id
}
`, name)
}

func testAccDataSourceNcloudNetworkInterfaceConfigFilter(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s"
	subnet             = "10.4.0.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_network_interface" "foo" {
	name                 = "%[1]s"
	description           = "for acc test"
	subnet_no             = ncloud_subnet.test.id
	access_control_groups = [ncloud_vpc.test.default_access_control_group_no]
}

data "ncloud_network_interface" "by_filter" {
	filter {
		name   = "network_interface_no"
		values = [ncloud_network_interface.foo.id]
	}
}
`, name)
}
