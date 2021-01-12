package ncloud

import (
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
