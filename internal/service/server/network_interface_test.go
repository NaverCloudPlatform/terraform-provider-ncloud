package server_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccresourceNcloudNetworkInterface_basic(t *testing.T) {
	var networkInterface vserver.NetworkInterface
	name := fmt.Sprintf("tf-nic-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_interface.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkInterfaceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					resource.TestMatchResourceAttr(resourceName, "network_interface_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "for acc test"),
					resource.TestCheckResourceAttr(resourceName, "private_ip", "10.4.0.6"),
					resource.TestCheckResourceAttr(resourceName, "server_instance_no", ""),
					resource.TestCheckResourceAttr(resourceName, "status", "NOTUSED"),
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
	name := fmt.Sprintf("tf-nic-update-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkInterfaceUpdate(name, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					resource.TestCheckResourceAttr(resourceName, "server_instance_no", ""),
				),
			},
			{
				Config: testAccResourceNcloudNetworkInterfaceUpdate(name, "${ncloud_server.server.id}"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					resource.TestMatchResourceAttr(resourceName, "server_instance_no", regexp.MustCompile(`^\d+$`)),
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkInterfaceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceExists(resourceName, &networkInterface),
					testAccCheckNetworkInterfaceDisappears(&networkInterface),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNcloudNetworkInterfaceConfig(name string) string {
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

func testAccResourceNcloudNetworkInterfaceUpdate(name, instanceNo string) string {
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
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_network_interface" "foo" {
	name                  = "%[1]s"
	description           = "for acc test"
	subnet_no             = ncloud_subnet.test.id
	private_ip            = "10.4.0.6"
	access_control_groups = [ncloud_vpc.test.default_access_control_group_no]
	server_instance_no    = "%[2]s"
}

resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
	login_key_name = ncloud_login_key.loginkey.key_name
}

`, name, instanceNo)
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

		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		instance, err := server.GetNetworkInterface(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*NetworkInterface = *instance

		return nil
	}
}

func testAccCheckNetworkInterfaceDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_network_interface" {
			continue
		}

		instance, err := server.GetNetworkInterface(config, rs.Primary.ID)

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
		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		return server.DeleteNetworkInterface(config, *instance.NetworkInterfaceNo)
	}
}
