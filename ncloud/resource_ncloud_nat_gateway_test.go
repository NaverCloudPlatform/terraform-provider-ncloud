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

func TestAccResourceNcloudNatGateway_basic(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", acctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "nat_gateway_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
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

func TestAccResourceNcloudNatGateway_disappears(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", acctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
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
	name := fmt.Sprintf("test-nat-gateway-%s", acctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfigOnlyRequiredParam(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "nat_gateway_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^[a-z0-9]+$`)),
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

func TestAccResourceNcloudNatGateway_updateName(t *testing.T) {
	var natGateway vpc.NatGatewayInstance
	name := fmt.Sprintf("test-nat-gateway-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("test-nat-gateway-update-%s", acctest.RandString(5))
	resourceName := "ncloud_nat_gateway.nat_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNatGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccResourceNcloudNatGatewayConfigUpdate(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceName, &natGateway),
				),
				ExpectError: regexp.MustCompile("Change 'name' is not support, Please set `name` as a old value"),
			},
		},
	})
}

func testAccResourceNcloudNatGatewayConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  zone        = "KR-1"
  name        = "%[1]s"
  description = "description"
}
`, name)
}

func testAccResourceNcloudNatGatewayConfigUpdate(name, updateName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  zone        = "KR-1"
  name        = "%s"
  description = "description"
}
`, name, updateName)
}

func testAccResourceNcloudNatGatewayConfigOnlyRequiredParam(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
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

		client := testAccProvider.Meta().(*ProviderConfig).Client
		instance, err := getNatGatewayInstance(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		*natGateway = *instance

		return nil
	}
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nat_gateway" {
			continue
		}

		instance, err := getNatGatewayInstance(client, rs.Primary.ID)

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
		client := testAccProvider.Meta().(*ProviderConfig).Client

		reqParams := &vpc.DeleteNatGatewayInstanceRequest{
			NatGatewayInstanceNo: instance.NatGatewayInstanceNo,
		}

		_, err := client.vpc.V2Api.DeleteNatGatewayInstance(reqParams)

		waitForNcloudNatGatewayDeletion(client, *instance.NatGatewayInstanceNo)

		return err
	}
}
