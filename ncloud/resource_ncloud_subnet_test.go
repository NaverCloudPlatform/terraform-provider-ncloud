package ncloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudSubnet_Basic(t *testing.T) {
	var subnet vpc.Subnet
	name := "testacc-subnet-basic"
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSubnetConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "subnet", cidr),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "status", "RUN"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudSubnet_UpdateName(t *testing.T) {
	var subnet vpc.Subnet
	name := "testacc-subnet-basic"
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSubnetConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
			{
				Config: testAccResourceNcloudSubnetConfig("testacc-subnet-update", cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
				ExpectError: regexp.MustCompile("Change 'name' is not support, Please set `name` as a old value"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudSubnet_UpdateNetworkACL(t *testing.T) {
	var subnet vpc.Subnet
	name := "testacc-subnet-update"
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSubnetConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
			{
				Config: testAccResourceNcloudSubnetConfigUpdateNetworkACL(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
		},
	})
}

func TestAccResourceNcloudSubnet_InvalidCIDR(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudSubnetConfigInvalidCIDR("10.3.2.0/24"),
				ExpectError: regexp.MustCompile("The subnet must belong to the IPv4 CIDR of the specified VPC."),
			},
		},
	})
}

func testAccResourceNcloudSubnetConfig(name, cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "foo" {
	name               = "testacc-subnet-basic"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "bar" {
	vpc_no             = ncloud_vpc.foo.vpc_no
	name               = "%s"
	subnet             = "%s"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.foo.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}
`, name, cidr)
}

func testAccResourceNcloudSubnetConfigUpdateNetworkACL(name, cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "foo" {
	name               = "testacc-subnet-basic"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.foo.vpc_no
	name        = "testacc-subnet-nacl-update"
	description = "for test acc"
}

resource "ncloud_subnet" "bar" {
	vpc_no             = ncloud_vpc.foo.vpc_no
	name               = "%s"
	subnet             = "%s"
	zone               = "KR-1"
	network_acl_no     = ncloud_network_acl.nacl.network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}
`, name, cidr)
}

func testAccResourceNcloudSubnetConfigInvalidCIDR(cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "foo" {
	name               = "testacc-subnet-invalid-cidr"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "bar" {
	vpc_no             = ncloud_vpc.foo.vpc_no
	name               = "testacc-subnet-invalid-cidr"
	subnet             = "%s"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.foo.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}
`, cidr)
}

func testAccCheckSubnetExists(n string, subnet *vpc.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No subnet no is set")
		}

		client := testAccProvider.Meta().(*NcloudAPIClient)
		instance, err := getSubnetInstance(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		*subnet = *instance

		return nil
	}
}
