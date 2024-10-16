package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlUsers_vpc_basic(t *testing.T) {
	testName := fmt.Sprintf("tf-mysqluser-%s", acctest.RandString(5))
	dataName := "data.ncloud_mysql_users.all"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMysqlUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlUsersConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "mysql_user_list.0.name", "test1"),
					resource.TestCheckResourceAttr(dataName, "mysql_user_list.1.name", "test2"),
				),
			},
		},
	})
}

func testAccDataSourceMysqlUsersConfig(testName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test_vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "test_subnet" {
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}

resource "ncloud_mysql" "mysql" {
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	host_ip = "192.168.0.1"
	database_name = "test_db"
}

resource "ncloud_mysql_users" "mysql_users" {
	mysql_instance_no = ncloud_mysql.mysql.id
	mysql_user_list = [
		{
			name = "test1",
			password = "t123456789!",
			host_ip = "%%",
			authority = "READ"
		},
		{
			name = "test2",
			password = "t123456789!",
			host_ip = "%%",
			authority = "DDL"
		}
	]
}

data "ncloud_mysql_users" "all" {
	mysql_instance_no = ncloud_mysql.mysql.id
	filter {
		name = "name"
		values = [ncloud_mysql_users.mysql_users.mysql_user_list.0.name, ncloud_mysql_users.mysql_users.mysql_user_list.1.name]
	}
}
`, testName)
}
