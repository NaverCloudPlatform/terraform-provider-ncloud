package ncloud

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudNetworkACLDenyAllowGroup_basic(t *testing.T) {
	var networkAclDenyAllowGroup vpc.NetworkAclDenyAllowGroup
	name := fmt.Sprintf("tf-nacl-allow-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl_deny_allow_group.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDenyAllowGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLDenyAllowGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLDenyAllowGroupExists(resourceName, &networkAclDenyAllowGroup),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_deny_allow_group_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
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

func TestAccResourceNcloudNetworkACLDenyAllowGroup_disappears(t *testing.T) {
	var networkAclDenyAllowGroup vpc.NetworkAclDenyAllowGroup
	name := fmt.Sprintf("tf-nacl-allow-ds-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl_deny_allow_group.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDenyAllowGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLDenyAllowGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLDenyAllowGroupExists(resourceName, &networkAclDenyAllowGroup),
					testAccCheckNetworkACLDenyAllowGroupDisappears(&networkAclDenyAllowGroup),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudNetworkACLDenyAllowGroup_update(t *testing.T) {
	var networkAclDenyAllowGroup vpc.NetworkAclDenyAllowGroup
	name := fmt.Sprintf("tf-nacl-allow-update-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl_deny_allow_group.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDenyAllowGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLDenyAllowGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLDenyAllowGroupExists(resourceName, &networkAclDenyAllowGroup),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
				),
			},
			{
				Config: testAccResourceNcloudNetworkACLDenyAllowGroupConfigUpdateIpList(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLDenyAllowGroupExists(resourceName, &networkAclDenyAllowGroup),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
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

func TestAccResourceNcloudNetworkACLDenyAllowGroup_description(t *testing.T) {
	var networkAclDenyAllowGroup vpc.NetworkAclDenyAllowGroup
	name := fmt.Sprintf("tf-nacl-allow-desc-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl_deny_allow_group.this"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLDenyAllowGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLDenyAllowGroupConfigDescription(name, "foo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLDenyAllowGroupExists(resourceName, &networkAclDenyAllowGroup),
					resource.TestCheckResourceAttr(resourceName, "description", "foo"),
				),
			},
			{
				Config: testAccResourceNcloudNetworkACLDenyAllowGroupConfigDescription(name, "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLDenyAllowGroupExists(resourceName, &networkAclDenyAllowGroup),
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

func testAccResourceNcloudNetworkACLDenyAllowGroupConfig(name string) string {
	return testAccResourceNcloudNetworkACLDenyAllowGroupConfigDescription(name, "for test acc")
}

func testAccResourceNcloudNetworkACLDenyAllowGroupConfigDescription(name, description string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl_deny_allow_group" "this" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "%[1]s"
	description = "%[2]s"
	ip_list     = ["10.0.0.1", "10.0.0.2"]
}
`, name, description)
}

func testAccResourceNcloudNetworkACLDenyAllowGroupConfigUpdateIpList(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl_deny_allow_group" "this" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "%[1]s"
	description = "for test acc"
	ip_list     = ["10.0.0.1"]
}
`, name)
}

func testAccCheckNetworkACLDenyAllowGroupExists(n string, networkACL *vpc.NetworkAclDenyAllowGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NetworkAclDenyAllowGroup id is set: %s", n)
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getNetworkAclDenyAllowGroupDetail(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*networkACL = *instance

		return nil
	}
}

func testAccCheckNetworkACLDenyAllowGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_network_acl" {
			continue
		}

		instance, err := getNetworkAclDenyAllowGroupDetail(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("NetworkAclDenyAllowGroup still exists")
		}
	}

	return nil
}

func testAccCheckNetworkACLDenyAllowGroupDisappears(instance *vpc.NetworkAclDenyAllowGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		reqParams := &vpc.DeleteNetworkAclDenyAllowGroupRequest{
			RegionCode:                 &config.RegionCode,
			NetworkAclDenyAllowGroupNo: instance.NetworkAclDenyAllowGroupNo,
		}

		_, err := config.Client.vpc.V2Api.DeleteNetworkAclDenyAllowGroup(reqParams)

		if err := waitForNcloudNetworkACLDeletion(config, *instance.NetworkAclDenyAllowGroupNo); err != nil {
			return err
		}

		return err
	}
}
