package postgresql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	postgresqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/postgresql"
)

func TestAccResourceNcloudPostgresqlDatabases_vpc_basic(t *testing.T) {
	testName := fmt.Sprintf("tf-postgresqldb-%s", acctest.RandString(5))
	resourceName := "ncloud_postgresql_databases.postgresql_db"
	dbResourceName := "ncloud_postgresql.postgresql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPostgresqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDatabasesConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "postgresql_database_list.0.name", "testdb1"),
					resource.TestCheckResourceAttr(resourceName, "postgresql_database_list.0.owner", "testowner1"),
					resource.TestCheckResourceAttr(resourceName, "postgresql_database_list.1.name", "testdb2"),
					resource.TestCheckResourceAttr(resourceName, "postgresql_database_list.1.owner", "testowner2"),
				),
				Destroy: false,
			},
			{
				Config: testAccPostgresqlDatabasesRemoveConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccPostgresqlDatabasesNotExist(dbResourceName, []string{"testdb1", "testdb2"}, GetTestProvider(true)),
				),
			},
		},
	})
}

func testAccPostgresqlDatabasesConfig(testPostgresqlName string) string {
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

resource "ncloud_postgresql" "postgresql" {
	vpc_no            = ncloud_vpc.test_vpc.vpc_no
	subnet_no         = ncloud_subnet.test_subnet.id
	service_name      = "%[1]s"
	server_name_prefix = "testprefix"
	user_name         = "testusername"
	user_password     = "t123456789!a"
	client_cidr       = "0.0.0.0/0"
	database_name     = "test_db"
}
resource "ncloud_postgresql_users" "postgresql_users" {
	id = ncloud_postgresql.postgresql.id
	postgresql_user_list = [
		{
			name = "testowner1",
			password = "t123456789!",
			client_cidr = "0.0.0.0/0",
			replication_role = "false"
		},
		{
			name = "testowner2",
			password = "t123456789!",
			client_cidr = "0.0.0.0/0",
			replication_role = "false"
		}
	]
}

resource "ncloud_postgresql_databases" "postgresql_db" {
	id = ncloud_postgresql.postgresql.id
	postgresql_database_list = [
		{
			name = "testdb1",
			owner = ncloud_postgresql_users.postgresql_users.postgresql_user_list.0.name
		},
		{
			name = "testdb2",
			owner = ncloud_postgresql_users.postgresql_users.postgresql_user_list.1.name
		}
	]
}
`, testPostgresqlName)
}

func testAccPostgresqlDatabasesRemoveConfig(testPostgresqlName string) string {
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

resource "ncloud_postgresql" "postgresql" {
	vpc_no            = ncloud_vpc.test_vpc.vpc_no
	subnet_no         = ncloud_subnet.test_subnet.id
	service_name      = "%[1]s"
	server_name_prefix = "testprefix"
	user_name         = "testusername"
	user_password     = "t123456789!a"
	client_cidr       = "0.0.0.0/0"
	database_name     = "test_db"
}

resource "ncloud_postgresql_users" "postgresql_users" {
	id = ncloud_postgresql.postgresql.id
	postgresql_user_list = [
		{
			name = "testowner1",
			password = "t123456789!",
			client_cidr = "0.0.0.0/0",
			replication_role = "false"
		},
		{
			name = "testowner2",
			password = "t123456789!",
			client_cidr = "0.0.0.0/0",
			replication_role = "false"
		}
	]
}
`, testPostgresqlName)
}

func testAccPostgresqlDatabasesNotExist(n string, dbNames []string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no Postgresql instance ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)

		dbs, err := postgresqlservice.GetPostgresqlDatabaseList(context.Background(), config, resource.Primary.ID, dbNames)
		if err != nil {
			return err
		}

		if dbs == nil {
			return nil
		}

		return fmt.Errorf("Postgresql dbs still exist: %v", dbs)
	}
}
