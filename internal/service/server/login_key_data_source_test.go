package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudLoginKey_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_login_key.all"
	testKeyName := "test-key"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLoginKeyConfig(testKeyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "login_key_list.0.key_name", "test-key1"),
					resource.TestCheckResourceAttr(dataName, "login_key_list.1.key_name", "test-key2"),
				),
			},
		},
	})
}

func testAccDataSourceLoginKeyConfig(keyName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey1" {
	key_name = "%[1]s1"
}
resource "ncloud_login_key" "loginkey2" {
	key_name = "%[1]s2"
}
data "ncloud_login_key" "all" {
	filter {
		name = "key_name"
		values = [ncloud_login_key.loginkey1.key_name, ncloud_login_key.loginkey2.key_name]
	}
}
`, keyName)
}
