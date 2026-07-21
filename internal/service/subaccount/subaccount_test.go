package subaccount_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	subaccountsdk "github.com/terraform-providers/terraform-provider-ncloud/internal/sdk/subaccount"
)

func TestAccResourceNcloudSubAccount_basic(t *testing.T) {
	loginId := GetTestPrefix() + "tfsub"
	resourceName := "ncloud_subaccount.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubAccountConfig(loginId, "tf test sub account"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "login_id", loginId),
					resource.TestCheckResourceAttr(resourceName, "name", "tf test sub account"),
					resource.TestCheckResourceAttr(resourceName, "can_api_gateway_access", "true"),
					resource.TestCheckResourceAttr(resourceName, "can_console_access", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_account_no"),
					resource.TestCheckResourceAttrSet(resourceName, "nrn"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: testAccSubAccountConfig(loginId, "tf test sub account updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "tf test sub account updated"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"generated_password"},
			},
		},
	})
}

func TestAccResourceNcloudSubAccountAccessKey_basic(t *testing.T) {
	loginId := GetTestPrefix() + "tfkey"
	resourceName := "ncloud_subaccount_access_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubAccountAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubAccountAccessKeyConfig(loginId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubAccountAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrPair(resourceName, "sub_account_id", "ncloud_subaccount.test", "id"),
				),
			},
		},
	})
}

func testAccCheckSubAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := TestAccProvider.Meta().(*conn.ProviderConfig)
		_, err := config.Client.SubAccount.GetSubAccount(context.Background(), rs.Primary.ID)
		return err
	}
}

func testAccCheckSubAccountAccessKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := TestAccProvider.Meta().(*conn.ProviderConfig)
		keys, err := config.Client.SubAccount.ListAccessKeys(context.Background(), rs.Primary.Attributes["sub_account_id"])
		if err != nil {
			return err
		}

		for _, key := range keys {
			if key.AccessKey == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("SubAccount access key (%s) not found via API", rs.Primary.ID)
	}
}

func testAccCheckSubAccountDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_subaccount" {
			continue
		}

		_, err := config.Client.SubAccount.GetSubAccount(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("SubAccount (%s) still exists", rs.Primary.ID)
		}
		if !subaccountsdk.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func testAccCheckSubAccountAccessKeyDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_subaccount_access_key" {
			continue
		}

		keys, err := config.Client.SubAccount.ListAccessKeys(context.Background(), rs.Primary.Attributes["sub_account_id"])
		if err != nil {
			// The parent sub account being gone implies the key is gone too.
			if subaccountsdk.IsNotFound(err) {
				continue
			}
			return err
		}

		for _, key := range keys {
			if key.AccessKey == rs.Primary.ID {
				return fmt.Errorf("SubAccount access key (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return testAccCheckSubAccountDestroy(s)
}

func testAccSubAccountConfig(loginId, name string) string {
	return fmt.Sprintf(`
resource "ncloud_subaccount" "test" {
  login_id = %[1]q
  name     = %[2]q
}
`, loginId, name)
}

func testAccSubAccountAccessKeyConfig(loginId string) string {
	return fmt.Sprintf(`
resource "ncloud_subaccount" "test" {
  login_id = %[1]q
  name     = "tf test access key"
}

resource "ncloud_subaccount_access_key" "test" {
  sub_account_id = ncloud_subaccount.test.id
}
`, loginId)
}
