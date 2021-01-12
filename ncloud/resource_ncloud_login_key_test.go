package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudLoginKey_classic_basic(t *testing.T) {
	testAccResourceNcloudLoginKeyBasic(t, false)
}

func TestAccResourceNcloudLoginKey_vpc_basic(t *testing.T) {
	testAccResourceNcloudLoginKeyBasic(t, true)
}

func testAccResourceNcloudLoginKeyBasic(t *testing.T, isVpc bool) {
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
	provider := getTestProvider(isVpc)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLoginKeyDestroyWithProvider(state, provider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLoginKeyConfig(testKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoginKeyExistsWithProvider("ncloud_login_key.loginkey", fingerprint, provider),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_login_key.loginkey",
						"key_name",
						testKeyName),
				),
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

func testAccCheckLoginKeyExistsWithProvider(n string, i *string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

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
