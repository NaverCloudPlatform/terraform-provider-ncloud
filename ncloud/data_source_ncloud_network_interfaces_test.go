package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNetworkInterfaces_basic(t *testing.T) {
	dataName := "data.ncloud_network_interfaces.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkInterfacesConfig(),
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNetworkInterfaces_privateIp(t *testing.T) {
	name := fmt.Sprintf("tf-ds-nic-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_interface.foo"
	dataName := "data.ncloud_network_interfaces.by_private_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkInterfacesConfigPrivateIp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.private_ip", resourceName, "private_ip"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.access_control_groups", resourceName, "access_control_groups"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.is_default", resourceName, "is_default"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNetworkInterfaces_filter(t *testing.T) {
	name := fmt.Sprintf("tf-nic-filter-%s", acctest.RandString(5))
	resourceName := "ncloud_network_interface.foo"
	dataName := "data.ncloud_network_interfaces.by_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkInterfacesConfigFilter(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_interfaces.by_filter"),
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.private_ip", resourceName, "private_ip"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.access_control_groups", resourceName, "access_control_groups"),
					resource.TestCheckResourceAttrPair(dataName, "network_interfaces.0.is_default", resourceName, "is_default"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNetworkInterfacesConfig() string {
	return `
data "ncloud_network_interfaces" "test" {}
`
}

func testAccDataSourceNcloudNetworkInterfacesConfigPrivateIp(name string) string {
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

data "ncloud_network_interfaces" "by_private_ip" {
	private_ip = ncloud_network_interface.foo.private_ip
}
`, name)
}

func testAccDataSourceNcloudNetworkInterfacesConfigFilter(name string) string {
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

data "ncloud_network_interfaces" "by_filter" {
	filter {
		name   = "id"
		values = [ncloud_network_interface.foo.id]
	}
}
`, name)
}
