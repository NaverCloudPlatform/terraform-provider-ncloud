package ncloud

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudSubnet_basic(t *testing.T) {
	var subnet vpc.Subnet
	name := fmt.Sprintf("test-subnet-basic-%s", acctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudSubnetConfig(name, cidr),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "subnet", cidr),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
				),
			},
			{
				SkipFunc:          testOnlyVpc,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudSubnet_disappears(t *testing.T) {
	var subnet vpc.Subnet
	name := fmt.Sprintf("test-subnet-disappears-%s", acctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudSubnetConfig(name, cidr),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
					testAccCheckSubnetDisappears(&subnet),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudSubnet_updateName(t *testing.T) {
	var subnet vpc.Subnet
	name := fmt.Sprintf("test-subnet-name-%s", acctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudSubnetConfig(name, cidr),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
			{
				Config:   testAccResourceNcloudSubnetConfig("testacc-subnet-update", cidr),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
				ExpectError: regexp.MustCompile("Change 'name' is not support, Please set `name` as a old value"),
			},
			{
				SkipFunc:          testOnlyVpc,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudSubnet_updateNetworkACL(t *testing.T) {
	var subnet vpc.Subnet
	name := fmt.Sprintf("test-subnet-update-nacl-%s", acctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudSubnetConfig(name, cidr),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
			{
				Config:   testAccResourceNcloudSubnetConfigUpdateNetworkACL(name, cidr),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
		},
	})
}

func TestAccResourceNcloudSubnet_InvalidCIDR(t *testing.T) {
	name := fmt.Sprintf("test-subnet-update-nacl-%s", acctest.RandString(5))
	cidr := "10.3.2.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudSubnetConfigInvalidCIDR(name, cidr),
				SkipFunc:    testOnlyVpc,
				ExpectError: regexp.MustCompile("The subnet must belong to the IPv4 CIDR of the specified VPC."),
			},
		},
	})
}

func testAccResourceNcloudSubnetConfig(name, cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "foo" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "bar" {
	vpc_no             = ncloud_vpc.foo.vpc_no
	name               = "%[1]s"
	subnet             = "%[2]s"
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
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.foo.vpc_no
	name        = "%[1]s"
	description = "for test acc"
}

resource "ncloud_subnet" "bar" {
	vpc_no             = ncloud_vpc.foo.vpc_no
	name               = "%[1]s"
	subnet             = "%[2]s"
	zone               = "KR-1"
	network_acl_no     = ncloud_network_acl.nacl.network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}
`, name, cidr)
}

func testAccResourceNcloudSubnetConfigInvalidCIDR(name, cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "foo" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "bar" {
	vpc_no             = ncloud_vpc.foo.vpc_no
	name               = "%[1]s"
	subnet             = "%s"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.foo.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}
`, name, cidr)
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

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getSubnetInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*subnet = *instance

		return nil
	}
}

func testAccCheckSubnetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_subnet" {
			continue
		}

		instance, err := getSubnetInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Subnet still exists")
		}
	}

	return nil
}

func testAccCheckSubnetDisappears(instance *vpc.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		reqParams := &vpc.DeleteSubnetRequest{
			RegionCode: &config.RegionCode,
			SubnetNo:   instance.SubnetNo,
		}

		_, err := config.Client.vpc.V2Api.DeleteSubnet(reqParams)

		if err := waitForNcloudSubnetDeletion(config, *instance.SubnetNo); err != nil {
			return err
		}

		return err
	}
}
