package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlDatabases_vpc_basic(t *testing.T) {
	/*
		TODO - it's	for atomicity of regression testing. remove when error has solved.
	*/
	t.Skip()

	dataName := "data.ncloud_mysql_databases.all"
	testName := fmt.Sprintf("tf-mysqldb-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMysqlDatabasesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlDatabasesConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "mysql_database_list.0.name", "testdb1"),
					resource.TestCheckResourceAttr(dataName, "mysql_database_list.1.name", "testdb2"),
				),
			},
		},
	})
}

func testAccDataSourceMysqlDatabasesConfig(testName string) string {
	return fmt.Sprintf(`
data "ncloud_vpc" "test_vpc" {
	id = "75658"
}

data "ncloud_subnet" "test_subnet" {
	id = "172709"
}

resource "ncloud_mysql" "mysql" {
	subnet_no = data.ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	host_ip = "192.168.0.1"
	database_name = "test_db"
}

resource "ncloud_mysql_databases" "mysql_db" {
	mysql_instance_no = ncloud_mysql.mysql.id
	mysql_database_list = [
		{
			name = "testdb1"
		},
		{
			name = "testdb2"
		}
	]
}

data "ncloud_mysql_databases" "all" {
	mysql_instance_no = ncloud_mysql.mysql.id
	filter {
		name = "name"
		values = [ncloud_mysql_databases.mysql_db.mysql_database_list.0.name, ncloud_mysql_databases.mysql_db.mysql_database_list.1.name]
	}
}
`, testName)
}
