package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudLoginKeyBasic(t *testing.T) {
	var fingerprint *string
	prefix := getTestPrefix()
	testKeyName := prefix + "-key"

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if fingerprint != nil {
				return fmt.Errorf("fingerprint must not be nil")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoginKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoginKeyConfig(testKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoginKeyExists("ncloud_login_key.loginkey", fingerprint),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_login_key.loginkey",
						"key_name",
						testKeyName),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            "ncloud_login_key.loginkey",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_name", "private_key"},
			},
		},
	})
}

func testAccCheckLoginKeyExists(n string, i *string) resource.TestCheckFunc {
	return testAccCheckLoginKeyExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckLoginKeyExistsWithProvider(n string, i *string, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		config := provider.Meta().(*ProviderConfig)
		fingerPrint, err := getFingerPrint(config, &rs.Primary.ID)
		if err != nil {
			return nil
		}

		if fingerPrint != nil {
			i = fingerPrint
			return nil
		}

		return fmt.Errorf("fingerprint is not found")
	}
}

func testAccCheckLoginKeyDestroy(s *terraform.State) error {
	return testAccCheckLoginKeyDestroyWithProvider(s, testAccProvider)
}

func testAccCheckLoginKeyDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_login_key" {
			continue
		}
		fingerPrint, err := getFingerPrint(config, &rs.Primary.ID)

		if fingerPrint == nil {
			continue
		}
		if err != nil {
			return err
		}
		if fingerPrint != nil && *fingerPrint != "" {
			return fmt.Errorf("found not deleted login key: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccLoginKeyConfig(keyName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%s"
}
`, keyName)
}
