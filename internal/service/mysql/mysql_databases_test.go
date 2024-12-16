package mysql_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mysqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/mysql"
)

func TestAccResourceNcloudMysqlDatabases_vpc_basic(t *testing.T) {
	/*
		TODO - it's	for atomicity of regression testing. remove when error has solved.
	*/
	t.Skip()

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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMysqlAssociationImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMysqlDatabasesConfig(testMysqlName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test_vpc" {
	name             = "%[1]s"
	ipv4_cidr_block  = "10.5.0.0/16"
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
	server_name_prefix = "rprefix"
	user_name = "rusername"
	user_password = "t123456789!a"
	host_ip = "192.168.0.1"
	database_name = "admin_r_db"
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
		instance, err := mysqlservice.GetMysqlDatabaseList(context.Background(), config, rs.Primary.ID, []string{"testdb1", "testdb2"})
		if err != nil && !strings.Contains(err.Error(), "5001067") {
			return err
		}

		if len(instance) > 1 {
			return errors.New("mysql database still exists")
		}
	}

	return nil
}

func testAccMysqlAssociationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		id := rs.Primary.Attributes["id"]
		name1 := rs.Primary.Attributes["mysql_database_list.0.name"]
		name2 := rs.Primary.Attributes["mysql_database_list.1.name"]

		return fmt.Sprintf("%s:%s:%s", id, name1, name2), nil
	}
}
