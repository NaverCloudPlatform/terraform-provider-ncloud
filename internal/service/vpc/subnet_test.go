package vpc_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	vpcservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
)

func TestAccResourceNcloudSubnet_basic(t *testing.T) {
	var subnet vpc.Subnet
	name := fmt.Sprintf("test-subnet-basic-%s", sdkacctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubnetDestroy,
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

func TestAccResourceNcloudSubnet_disappears(t *testing.T) {
	var subnet vpc.Subnet
	name := fmt.Sprintf("test-subnet-disappears-%s", sdkacctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSubnetConfig(name, cidr),
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
	name := fmt.Sprintf("test-subnet-name-%s", sdkacctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSubnetConfig(name, cidr),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
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

func TestAccResourceNcloudSubnet_updateNetworkACL(t *testing.T) {
	/*
		TODO - it's	for atomicity of regression testing. remove when error has solved.
	*/
	t.Skip()

	var subnet vpc.Subnet
	name := fmt.Sprintf("test-subnet-update-nacl-%s", sdkacctest.RandString(5))
	cidr := "10.2.2.0/24"
	resourceName := "ncloud_subnet.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubnetDestroy,
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
	name := fmt.Sprintf("test-subnet-update-nacl-%s", sdkacctest.RandString(5))
	cidr := "10.3.2.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudSubnetConfigInvalidCIDR(name, cidr),
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

		config := TestAccProvider.Meta().(*conn.ProviderConfig)
		instance, err := vpcservice.GetSubnetInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*subnet = *instance

		return nil
	}
}

func testAccCheckSubnetDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_subnet" {
			continue
		}

		instance, err := vpcservice.GetSubnetInstance(config, rs.Primary.ID)

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
		config := TestAccProvider.Meta().(*conn.ProviderConfig)

		reqParams := &vpc.DeleteSubnetRequest{
			RegionCode: &config.RegionCode,
			SubnetNo:   instance.SubnetNo,
		}

		_, err := config.Client.Vpc.V2Api.DeleteSubnet(reqParams)

		if err := vpcservice.WaitForNcloudSubnetDeletion(config, *instance.SubnetNo); err != nil {
			return err
		}

		return err
	}
}
