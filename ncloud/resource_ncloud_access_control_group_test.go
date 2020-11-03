package ncloud

import (
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudAccessControlGroup_basic(t *testing.T) {
	var AccessControlGroup vserver.AccessControlGroup
	name := fmt.Sprintf("tf-acg-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_access_control_group.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessControlGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudAccessControlGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessControlGroupExists(resourceName, &AccessControlGroup),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "access_control_group_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "for acc test"),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(resourceName, "inbound.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "outbound.#", "1"),
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

func TestAccResourceNcloudAccessControlGroup_rule(t *testing.T) {
	var AccessControlGroup vserver.AccessControlGroup
	name := fmt.Sprintf("tf-acg-rule-%s", acctest.RandString(5))
	resourceNameFoo := "ncloud_access_control_group.foo"
	resourceNameBar := "ncloud_access_control_group.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessControlGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudAccessControlGroupConfigRule(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessControlGroupExists(resourceNameFoo, &AccessControlGroup),
					testAccCheckAccessControlGroupExists(resourceNameBar, &AccessControlGroup),
					resource.TestCheckResourceAttr(resourceNameFoo, "inbound.#", "2"),
					resource.TestCheckResourceAttr(resourceNameFoo, "outbound.#", "1"),
					resource.TestCheckResourceAttr(resourceNameBar, "inbound.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNcloudAccessControlGroup_disappears(t *testing.T) {
	var AccessControlGroup vserver.AccessControlGroup
	name := fmt.Sprintf("tf-nic-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_access_control_group.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessControlGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudAccessControlGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessControlGroupExists(resourceName, &AccessControlGroup),
					testAccCheckAccessControlGroupDisappears(&AccessControlGroup),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNcloudAccessControlGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_access_control_group" "foo" {
	name                  = "%[1]s"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id

	inbound {
		protocol    = "TCP"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
	
	inbound {
		protocol    = "TCP"
		port_range  = "443"
		ip_block    = "0.0.0.0/0"
	}

	outbound {
		protocol    = "TCP"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
}
`, name)
}

func testAccResourceNcloudAccessControlGroupConfigRule(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_access_control_group" "foo" {
	name                  = "%[1]s-foo"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id

	inbound {
		protocol    = "TCP"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
	
	inbound {
		protocol    = "TCP"
		port_range  = "443"
		ip_block    = "0.0.0.0/0"
	}

	outbound {
		protocol    = "TCP"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
}

resource "ncloud_access_control_group" "bar" {
	name                  = "%[1]s-bar"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id

	inbound {
		protocol    = "TCP"
		port_range  = "80"
		source_access_control_group_no = ncloud_access_control_group.foo.id
	}
}
`, name)
}

func testAccCheckAccessControlGroupExists(n string, AccessControlGroup *vserver.AccessControlGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Access Control Group id is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getAccessControlGroup(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*AccessControlGroup = *instance

		return nil
	}
}

func testAccCheckAccessControlGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_access_control_group" {
			continue
		}

		instance, err := getAccessControlGroup(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Access Control Group still exists")
		}
	}

	return nil
}

func testAccCheckAccessControlGroupDisappears(instance *vserver.AccessControlGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)
		return deleteAccessControlGroup(config, *instance.AccessControlGroupNo)
	}
}
