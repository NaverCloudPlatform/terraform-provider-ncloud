package mongodb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mongodbservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/mongodb"
)

func TestAccResourceNcloudMongoDbUsers_vpc_basic(t *testing.T) {
	testName := fmt.Sprintf("tf-monuser-%s", acctest.RandString(3))
	resourceName := "ncloud_mongodb_users.mongodb_users"
	dbResourceName := "ncloud_mongodb.mongodb"
	testUserBefore := "testuser"
	testUserAfter := "user"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDbUsersConfig(testName, testUserBefore),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.0.name", "testuser1"),
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.0.database_name", "testdb1"),
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.0.authority", "READ"),
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.1.name", "testuser2"),
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.1.database_name", "testdb2"),
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.1.authority", "READ_WRITE"),
				),
			},
			{
				Config: testAccMongoDbUsersConfig(testName, testUserAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.0.name", "user1"),
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_list.1.name", "user2"),
				),
				Destroy: false,
			},
			{
				Config: testAccMongoDbUsersRemoveConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDbUsersNotExists(dbResourceName, []string{"testuser1", "testuser2"}, GetTestProvider(true)),
				),
			},
		},
	})
}

func testAccMongoDbUsersConfig(testMongodbName, testUserName string) string {
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

resource "ncloud_mongodb" "mongodb" {
	vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
    server_name_prefix = "ex-svr"
	cluster_type_code = "STAND_ALONE"
	user_name = "testuser"
	user_password = "t123456789!"
}

resource "ncloud_mongodb_users" "mongodb_users" {
	mongodb_instance_no = ncloud_mongodb.mongodb.id
	mongodb_user_list = [
		{
			name = "%[2]s1",
			password = "t123456789!",
			database_name = "testdb1",
			authority = "READ"
		},
		{
			name = "%[2]s2",
			password = "t123456789!",
			database_name = "testdb2",
			authority = "READ_WRITE"
		}
	]
}
`, testMongodbName, testUserName)
}

func testAccMongoDbUsersRemoveConfig(testName string) string {
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

resource "ncloud_mongodb" "mongodb" {
	vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
    server_name_prefix = "ex-svr"
	cluster_type_code = "STAND_ALONE"
	user_name = "testuser"
	user_password = "t123456789!"
}
`, testName)
}

func testAccCheckMongoDbUsersNotExists(n string, userNames []string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no MongoDB instance ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)

		users, err := mongodbservice.GetMongoDbUserList(context.Background(), config, resource.Primary.ID, userNames)
		if err != nil {
			return err
		}

		if users == nil {
			return nil
		}

		return fmt.Errorf("MongoDB users still exist: %v", users)
	}
}
