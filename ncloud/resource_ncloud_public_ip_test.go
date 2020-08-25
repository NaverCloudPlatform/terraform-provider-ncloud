package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudPublicIpInstance_basic(t *testing.T) {
	instance := map[string]interface{}{}

	description := fmt.Sprintf("test-public-ip-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	testCheckAttribute := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			config := testAccProvider.Meta().(*ProviderConfig)

			if instance["instance_status"] != GetValueClassicOrVPC(config, "CREAT", "RUN") {
				return fmt.Errorf("invalid public ip status: %s ", instance["instance_status"])
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicIpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpInstanceConfig(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicIpInstanceExists(resourceName, instance),
					testCheckAttribute(),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"server_instance_no"},
			},
		},
	})
}

func TestAccResourceNcloudPublicIpInstance_update(t *testing.T) {
	instance := map[string]interface{}{}
	serverNameFoo := fmt.Sprintf("test-public-ip-basic-%s", acctest.RandString(5))
	serverNameBar := fmt.Sprintf("test-public-ip-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	log.Print(testAccPublicIpInstanceConfigServer(serverNameFoo, serverNameBar, "${ncloud_server.foo.id}"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicIpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpInstanceConfigServer(serverNameFoo, serverNameBar, ""),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config: testAccPublicIpInstanceConfigServer(serverNameFoo, serverNameBar, "${ncloud_server.foo.id}"),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config: testAccPublicIpInstanceConfigServer(serverNameFoo, serverNameBar, "${ncloud_server.bar.id}"),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config: testAccPublicIpInstanceConfigServer(serverNameFoo, serverNameBar, ""),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
		},
	})
}

func testAccCheckPublicIpInstanceExists(n string, i map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)

		instance, err := getPublicIp(config, rs.Primary.ID)

		if err != nil {
			return nil
		}

		if instance != nil {
			for k, v := range instance {
				i[k] = v
			}

			return nil
		}

		return fmt.Errorf("public ip instance not found")
	}
}

func testAccCheckPublicIpInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		config := testAccProvider.Meta().(*ProviderConfig)

		if rs.Type != "ncloud_public_ip" {
			continue
		}

		instance, err := getPublicIp(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("Public IP still exists:\n%#v", instance)
		}

		break
	}

	return nil
}

func testAccPublicIpInstanceConfig(description string) string {
	return fmt.Sprintf(`
resource "ncloud_public_ip" "public_ip" {
	description = "%s"
}
`, description)
}

func testAccPublicIpInstanceConfigServer(serverNameFoo, serverNameBar, serverInstanceNo string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "foo" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_server" "bar" {
	name = "%[2]s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_public_ip" "public_ip" {
	description = "%[1]s"
	server_instance_no = "%[3]s"
}
`, serverNameFoo, serverNameBar, serverInstanceNo)
}
