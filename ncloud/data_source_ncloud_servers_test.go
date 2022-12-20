package ncloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	dataName   = "data.ncloud_servers.by_id"
	dataName2  = "data.ncloud_servers.by_filter"
	serverName = "ncloud_server.test"
)

func TestAccDataSourceNcloudServers_vpc_basic(t *testing.T) {
	testServerName := getTestServerName()
	testServerName2 := getTestServerName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServersVpcConfig(testServerName, testServerName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(dataName, "ids.#", "2"),
					testAccCheckDataSourceID(dataName2),
					resource.TestCheckResourceAttr(dataName2, "ids.#", "1"),
					resource.TestCheckResourceAttrPair(dataName2, "ids.0", serverName, "id"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServers_classic_basic(t *testing.T) {
	testServerName := getTestServerName()
	testServerName2 := getTestServerName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServersClassicConfig(testServerName, testServerName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(dataName, "ids.#", "2"),
					testAccCheckDataSourceID(dataName2),
					resource.TestCheckResourceAttr(dataName2, "ids.#", "1"),
					resource.TestCheckResourceAttrPair(dataName2, "ids.0", serverName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceServersVpcConfig(testServerName, testServerName2 string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
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

resource "ncloud_server" "test" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_server" "test2" {
	subnet_no = ncloud_subnet.test.id
	name = "%[2]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.loginkey.key_name
}

data "ncloud_servers" "by_id" {
	ids = [ncloud_server.test.id, ncloud_server.test2.id]
}

data "ncloud_servers" "by_filter" {
	filter {
		name = "instance_no"
		values = [ncloud_server.test.id]
	}
}
`, testServerName, testServerName2)
}

func testAccDataSourceServersClassicConfig(testServerName, testServerName2 string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "test" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000139"
	server_product_code = "SPSVRSSD00000003"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_server" "test2" {
	name = "%[2]s"
	server_image_product_code = "SPSW0LINUX000139"
	server_product_code = "SPSVRSSD00000003"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

data "ncloud_servers" "by_id" {
	ids = [ncloud_server.test.id, ncloud_server.test2.id]
}

data "ncloud_servers" "by_filter" {
	filter {
		name = "instance_no"
		values = [ncloud_server.test.id]
	}
}
`, testServerName, testServerName2)
}
