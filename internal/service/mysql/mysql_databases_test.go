package mysql_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mysqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/mysql"
)

func TestAccResourceNcloudMysqlDatabases_vpc_basic(t *testing.T) {
	testName := fmt.Sprintf("tf-mysqldb-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql_databases.mysql_dbs"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMysqlDatabasesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDatabasesConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mysql_database_list.0.name", "testdb1"),
					resource.TestCheckResourceAttr(resourceName, "mysql_database_list.1.name", "testdb2"),
				),
			},
		},
	})
}

func testAccMysqlDatabasesConfig(testMysqlName string) string {
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

resource "ncloud_mysql_databases" "mysql_dbs" {
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
`, testMysqlName)
}

func testAccCheckMysqlDatabasesDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_mysql_databases" {
			continue
		}
		instance, err := mysqlservice.GetMysqlDatabaseList(context.Background(), config, rs.Primary.ID)
		if err != nil && !mysqlservice.CheckIfAlreadyDeleted(err) {
			return err
		}

		if len(instance) > 1 {
			return errors.New("mysql database still exists")
		}
	}

	return nil
}
