package ncloud

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNcloudPublicIPInstance_basic(t *testing.T) {
	var publicIPInstance sdk.PublicIPInstance
	testServerInstanceName := getTestServerName()
	testPublicIPDescription := "acceptanceTest"
	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if publicIPInstance.ServerInstance.ServerName != testServerInstanceName {
				return fmt.Errorf("not found: %s", testServerInstanceName)
			}
			if publicIPInstance.PublicIPDescription != testPublicIPDescription {
				return fmt.Errorf("invalid public ip description: %s ", publicIPInstance.PublicIPDescription)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_public_ip.public_ip",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPublicIPInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIPInstanceConfig(testServerInstanceName, testPublicIPDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicIPInstanceExists(
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
		},
	})
}

func testAccCheckPublicIPInstanceExists(n string, i *sdk.PublicIPInstance) resource.TestCheckFunc {
	return testAccCheckPublicIPInstanceExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckPublicIPInstanceExistsWithProvider(n string, i *sdk.PublicIPInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		conn := provider.Meta().(*NcloudSdk).conn
		instance, err := getPublicIPInstance(conn, rs.Primary.ID)
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

func testAccCheckPublicIPInstanceDestroy(s *terraform.State) error {
	return testAccCheckPublicIPInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckPublicIPInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*NcloudSdk).conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_public_ip" {
			continue
		}

		instance, err := getPublicIPInstance(conn, rs.Primary.ID)

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

func testAccPublicIPInstanceConfig(testServerInstanceName string, testPublicIPDescription string) string {
	return fmt.Sprintf(`
resource "ncloud_instance" "test" {
	"server_name" = "%s"
	"server_image_product_code" = "SPSW0LINUX000032"
	"server_product_code" = "SPSVRSTAND000004"
}

resource "ncloud_public_ip" "public_ip" {
	"server_instance_no" = "${ncloud_instance.test.id}"
	"public_ip_description" = "%s"
	"region_no" = "1"
	"zone_no" = "2"
}
`, testServerInstanceName, testPublicIPDescription)
}
