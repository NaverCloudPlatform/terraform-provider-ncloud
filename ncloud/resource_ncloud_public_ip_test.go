package ncloud

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceNcloudPublicIpInstanceBasic(t *testing.T) {
	var publicIPInstance server.PublicIpInstance
	testServerInstanceName := getTestServerName()
	testPublicIpDescription := "acceptanceTest"
	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if *publicIPInstance.ServerInstanceAssociatedWithPublicIp.ServerName != testServerInstanceName {
				return fmt.Errorf("not found: %s", testServerInstanceName)
			}
			if *publicIPInstance.PublicIpDescription != testPublicIpDescription {
				return fmt.Errorf("invalid public ip description: %s ", *publicIPInstance.PublicIpDescription)
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
				Config: testAccPublicIpInstanceConfig(testServerInstanceName, testPublicIpDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicIpInstanceExists(
						"ncloud_public_ip.public_ip", &publicIPInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_public_ip.public_ip",
						"region_no",
						"1"),
					resource.TestCheckResourceAttr(
						"ncloud_public_ip.public_ip",
						"zone_no",
						"2"),
				),
			},
			{
				ResourceName:      "ncloud_public_ip.public_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPublicIpInstanceExists(n string, i *server.PublicIpInstance) resource.TestCheckFunc {
	return testAccCheckPublicIpInstanceExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckPublicIpInstanceExistsWithProvider(n string, i *server.PublicIpInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		client := provider.Meta().(*NcloudAPIClient)
		instance, err := getPublicIpInstance(client, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if instance != nil {
			*i = *instance
			return nil
		}

		return fmt.Errorf("public ip instance not found")
	}
}

func testAccCheckPublicIpInstanceDestroy(s *terraform.State) error {
	return testAccCheckPublicIpInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckPublicIpInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*NcloudAPIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_public_ip" {
			continue
		}

		instance, err := getPublicIpInstance(client, rs.Primary.ID)

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

func testAccPublicIpInstanceConfig(testServerInstanceName string, testPublicIpDescription string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	"key_name" = "%s-key"
}

resource "ncloud_server" "test" {
	"name" = "%s"
	"server_image_product_code" = "SPSW0LINUX000032"
	"server_product_code" = "SPSVRSTAND000004"
	"login_key_name" = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_public_ip" "public_ip" {
	"server_instance_no" = "${ncloud_server.test.id}"
	"description" = "%s"
	"region_no" = "1"
	"zone_no" = "2"
}
`, testServerInstanceName, testServerInstanceName, testPublicIpDescription)
}
