package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudRootPassword_classic_basic(t *testing.T) {
	resourceName := "data.ncloud_root_password.default"
	name := fmt.Sprintf("tf-passwd-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRootPasswordClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudRootPassword_vpc_basic(t *testing.T) {
	resourceName := "data.ncloud_root_password.default"
	name := fmt.Sprintf("tf-passwd-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRootPasswordVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func testAccDataSourceRootPasswordClassicConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "key" {
  key_name = "%[1]s-key"
}

resource "ncloud_server" "server" {
  name = "%[1]s"
  server_image_product_code = "SPSW0LINUX000046"
  server_product_code = "SPSVRSTAND000004"
  login_key_name = "${ncloud_login_key.key.key_name}"
}

data "ncloud_root_password" "default" {
  server_instance_no = ncloud_server.server.id
  private_key = ncloud_login_key.key.private_key
}
`, name)
}

func testAccDataSourceRootPasswordVpcConfig(testServerName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "key" {
	key_name = "%[1]s-key"
}

resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.key.key_name
}

data "ncloud_root_password" "default" {
  server_instance_no = ncloud_server.server.id
  private_key = ncloud_login_key.key.private_key
}
`, testServerName)
}
