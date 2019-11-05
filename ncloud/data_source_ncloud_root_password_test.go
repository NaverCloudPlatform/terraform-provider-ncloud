package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudRootPasswordBasic(t *testing.T) {

	resourceName := "data.ncloud_root_password.default"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRootPasswordBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func testAccDataSourceRootPasswordBasicConfig() string {
	prefix := getTestPrefix()
	return fmt.Sprintf(`
resource "ncloud_login_key" "key" {
  key_name = "%s-key"
}

resource "ncloud_server" "server" {
  name = "%s-vm"
  server_image_product_code = "SPSW0LINUX000032"
  server_product_code = "SPSVRSTAND000004"
  login_key_name = "${ncloud_login_key.key.key_name}"
}

data "ncloud_root_password" "default" {
  server_instance_no = "${ncloud_server.server.id}"
  private_key = "${ncloud_login_key.key.private_key}"
}
`, prefix, prefix)
}
