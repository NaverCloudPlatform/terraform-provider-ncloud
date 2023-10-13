package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccResourceNcloudLoginKey_classic_basic(t *testing.T) {
	testAccResourceNcloudLoginKeyBasic(t, false)
}

func TestAccResourceNcloudLoginKey_vpc_basic(t *testing.T) {
	testAccResourceNcloudLoginKeyBasic(t, true)
}

func testAccResourceNcloudLoginKeyBasic(t *testing.T, isVpc bool) {
	var loginKey *server.LoginKey
	prefix := GetTestPrefix()
	testKeyName := prefix + "-key"

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if loginKey != nil {
				return fmt.Errorf("loginkey must not be nil")
			}
			return nil
		}
	}
	provider := GetTestProvider(isVpc)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: getProvidersBasedOnVpc(isVpc),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLoginKeyDestroyWithProvider(state, provider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLoginKeyConfig(testKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoginKeyExistsWithProvider("ncloud_login_key.loginkey", loginKey, provider),
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

func getProvidersBasedOnVpc(isVpc bool) map[string]func() (tfprotov6.ProviderServer, error) {
	if isVpc {
		return ProtoV6ProviderFactories
	} else {
		return ClassicProtoV6ProviderFactories
	}
}

func testAccCheckLoginKeyExistsWithProvider(n string, l *server.LoginKey, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		loginKey, err := server.GetLoginKey(config, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if loginKey != nil {
			l = loginKey
			return nil
		}

		return fmt.Errorf("loginKey is not found")
	}
}

func testAccCheckLoginKeyDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_login_key" {
			continue
		}
		loginKey, err := server.GetLoginKey(config, rs.Primary.ID)

		if loginKey == nil {
			continue
		}
		if err != nil {
			return err
		}
		if loginKey != nil && *loginKey.Fingerprint != "" {
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
