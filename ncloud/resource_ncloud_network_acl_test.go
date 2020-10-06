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

func TestAccResourceNcloudNetworkACL_basic(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
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

func TestAccResourceNcloudNetworkACL_disappears(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					testAccCheckNetworkACLDisappears(&networkACL),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudNetworkACL_onlyRequiredParam(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfigOnlyRequired(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
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

func TestAccResourceNcloudNetworkACL_updateName(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfigOnlyRequired(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
				),
			},
			{
				Config: testAccResourceNcloudNetworkACLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
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

func TestAccResourceNcloudNetworkACL_description(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-desc-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfigDescription(name, "foo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestCheckResourceAttr(resourceName, "description", "foo"),
				),
			},
			{
				Config: testAccResourceNcloudNetworkACLConfigDescription(name, "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
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

func testAccResourceNcloudNetworkACLConfig(name string) string {
	return testAccResourceNcloudNetworkACLConfigDescription(name, "for test acc")
}

func testAccResourceNcloudNetworkACLConfigDescription(name, description string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "%[1]s"
	description = "%[2]s"
}
`, name, description)
}

func testAccResourceNcloudNetworkACLConfigOnlyRequired(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
}
`, name)
}

func testAccCheckNetworkACLExists(n string, networkACL *vpc.NetworkAcl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network acl id is set: %s", n)
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getNetworkACLInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*networkACL = *instance

		return nil
	}
}

func testAccCheckNetworkACLDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_network_acl" {
			continue
		}

		instance, err := getNetworkACLInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Network ACL still exists")
		}
	}

	return nil
}

func testAccCheckNetworkACLDisappears(instance *vpc.NetworkAcl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		reqParams := &vpc.DeleteNetworkAclRequest{
			RegionCode:   &config.RegionCode,
			NetworkAclNo: instance.NetworkAclNo,
		}

		_, err := config.Client.vpc.V2Api.DeleteNetworkAcl(reqParams)

		if err := waitForNcloudNetworkACLDeletion(config, *instance.NetworkAclNo); err != nil {
			return err
		}

		return err
	}
}
