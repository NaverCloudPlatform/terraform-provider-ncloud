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

func TestAccResourceNcloudNatGateway_basic(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", sdkacctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"
	resourcePrivate := "ncloud_nat_gateway.nat_gateway_private"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "nat_gateway_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),

					testAccCheckNatGatewayExists(resourcePrivate, &natGateway),
					resource.TestMatchResourceAttr(resourcePrivate, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourcePrivate, "nat_gateway_no", regexp.MustCompile(`^\d+$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourcePrivate,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudNatGateway_disappears(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", sdkacctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					testAccCheckNatGatewayDisappears(&natGateway),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudNatGateway_onlyRequiredParam(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", sdkacctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfigOnlyRequiredParam(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "nat_gateway_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^[a-z0-9]+$`)),
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
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

func TestAccResourceNcloudNatGateway_updateName(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", sdkacctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func TestAccResourceNcloudNatGateway_description(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", sdkacctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfigDescription(name, "foo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "description", "foo"),
				),
			},
			{
				Config: testAccResourceNcloudNatGatewayConfigDescription(name, "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "description", "bar"),
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

func testAccResourceNcloudNatGatewayConfig(name string) string {
	return testAccResourceNcloudNatGatewayConfigDescription(name, "for acc test")
}

func testAccResourceNcloudNatGatewayConfigDescription(name, description string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_subnet" "subnet_public" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = cidrsubnet(ncloud_vpc.vpc.ipv4_cidr_block, 8, 1)
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
  usage_type     = "NATGW"
}

resource "ncloud_subnet" "subnet_private" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = cidrsubnet(ncloud_vpc.vpc.ipv4_cidr_block, 8, 2)
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  usage_type     = "NATGW"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  subnet_no   = ncloud_subnet.subnet_public.id
  zone        = "KR-1"
  name        = "%[1]s"
  description = "%[2]s"
}

resource "ncloud_nat_gateway" "nat_gateway_private" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  subnet_no   = ncloud_subnet.subnet_private.id
  zone        = "KR-1"
  description = "%[2]s"
}
`, name, description)
}

func testAccResourceNcloudNatGatewayConfigOnlyRequiredParam(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_subnet" "subnet_public" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = cidrsubnet(ncloud_vpc.vpc.ipv4_cidr_block, 8, 1)
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
  usage_type     = "NATGW"
}
	
resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  subnet_no   = ncloud_subnet.subnet_public.id
  zone        = "KR-1"
}
`, name)
}

func testAccCheckNatGatewayExists(n string, natGateway *vpc.NatGatewayInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NAT Gateway id is set")
		}

		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		instance, err := vpcservice.GetNatGatewayInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*natGateway = *instance

		return nil
	}
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nat_gateway" {
			continue
		}

		instance, err := vpcservice.GetNatGatewayInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("NAT Gateway still exists")
		}
	}

	return nil
}

func testAccCheckNatGatewayDisappears(instance *vpc.NatGatewayInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

		reqParams := &vpc.DeleteNatGatewayInstanceRequest{
			RegionCode:           &config.RegionCode,
			NatGatewayInstanceNo: instance.NatGatewayInstanceNo,
		}

		_, err := config.Client.Vpc.V2Api.DeleteNatGatewayInstance(reqParams)

		if err := vpcservice.WaitForNcloudNatGatewayDeletion(config, *instance.NatGatewayInstanceNo); err != nil {
			return err
		}

		return err
	}
}
