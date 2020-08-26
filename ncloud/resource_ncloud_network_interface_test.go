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

func TestAccresourceNcloudNetworkInterface_basic(t *testing.T) {
	var networkInterface vserver.NetworkInterface
	name := fmt.Sprintf("tf-nic-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_interface.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccresourceNcloudNetworkInterfaceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					resource.TestMatchResourceAttr(resourceName, "network_interface_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "for acc test"),
					resource.TestCheckResourceAttr(resourceName, "private_ip", "10.4.0.6"),
					resource.TestCheckResourceAttr(resourceName, "server_instance_no", ""),
					resource.TestCheckResourceAttr(resourceName, "status", "RUN"),
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "access_control_groups.#", "1"),
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

func TestAccresourceNcloudNetworkInterface_update(t *testing.T) {
	var networkInterface vserver.NetworkInterface
	resourceName := "ncloud_network_interface.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccresourceNcloudNetworkInterfaceUpdate(""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
				),
			},
			{
				Config: testAccresourceNcloudNetworkInterfaceUpdate("1324440"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
				),
			},
			{
				Config: testAccresourceNcloudNetworkInterfaceUpdate(""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
				),
			},
		},
	})
}

func TestAccresourceNcloudNetworkInterface_disappears(t *testing.T) {
	var networkInterface vserver.NetworkInterface
	name := fmt.Sprintf("tf-nic-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_network_interface.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccresourceNcloudNetworkInterfaceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					testAccCheckNetworkInterfaceDisappears(&networkInterface),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccresourceNcloudNetworkInterfaceConfig(name string) string {
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
	name                  = "%[1]s"
	description           = "for acc test"
	subnet_no             = ncloud_subnet.test.id
	private_ip            = "10.4.0.6"
	access_control_groups = [ncloud_vpc.test.default_access_control_group_no]
}
`, name)
}

func testAccresourceNcloudNetworkInterfaceUpdate(instanceNo string) string {
	// TODO: update test case after vpc server developed
	return fmt.Sprintf(`
resource "ncloud_network_interface" "foo" {
	description           = "for acc test"
	subnet_no             = "906"
	access_control_groups = ["1511"]
	server_instance_no    = "%s"
}
`, instanceNo)
}

func testAccCheckNetworkInterfaceExists(n string, NetworkInterface *vserver.NetworkInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Network Interface id is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getNetworkInterface(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*NetworkInterface = *instance

		return nil
	}
}

func testAccCheckNetworkInterfaceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_network_interface" {
			continue
		}

		instance, err := getNetworkInterface(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("network interface still exists")
		}
	}

	return nil
}

func testAccCheckNetworkInterfaceDisappears(instance *vserver.NetworkInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)
		return deleteNetworkInterface(config, *instance.NetworkInterfaceNo)
	}
}
