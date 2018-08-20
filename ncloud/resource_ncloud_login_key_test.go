package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccResourceNcloudLoginKeyBasic(t *testing.T) {
	var loginKey server.LoginKey
	prefix := getTestPrefix()
	testKeyName := prefix + "-key"

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if *loginKey.KeyName != testKeyName {
				return fmt.Errorf("not found: %s", testKeyName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_login_key.loginkey",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoginKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoginKeyConfig(testKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoginKeyExists("ncloud_login_key.loginkey", &loginKey),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_login_key.loginkey",
						"key_name",
						testKeyName),
				),
			},
		},
	})
}

func testAccCheckLoginKeyExists(n string, i *server.LoginKey) resource.TestCheckFunc {
	return testAccCheckLoginKeyExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckLoginKeyExistsWithProvider(n string, i *server.LoginKey, providerF func() *schema.Provider) resource.TestCheckFunc {
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
		loginKey, err := getLoginKey(client, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if loginKey != nil {
			*i = *loginKey
			return nil
		}

		return fmt.Errorf("login key not found")
	}
}

func testAccCheckLoginKeyDestroy(s *terraform.State) error {
	return testAccCheckLoginKeyDestroyWithProvider(s, testAccProvider)
}

func testAccCheckLoginKeyDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*NcloudAPIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_login_key" {
			continue
		}
		loginKey, err := getLoginKey(client, rs.Primary.ID)

		if loginKey == nil {
			continue
		}
		if err != nil {
			return err
		}
		if loginKey != nil && *loginKey.Fingerprint != "" {
			return fmt.Errorf("found not deleted login key: %s", *loginKey.KeyName)
		}
	}

	return nil
}

func testAccLoginKeyConfig(keyName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	"key_name" = "%s"
}
`, keyName)
}
