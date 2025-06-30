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

func TestAccResourceNcloudMongoDbUsers_vpc_update(t *testing.T) {
	testName := fmt.Sprintf("tf-monuser-%s", acctest.RandString(3))
	resourceName := "ncloud_mongodb_users.mongodb_users"
	dbResourceName := "ncloud_mongodb.mongodb"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDbUsersConfig(testName, "READ", "READ_WRITE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_set.#", "2"),
				),
			},
			{
				Config: testAccMongoDbUsersConfig(testName, "READ_WRITE", "READ"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mongodb_user_set.#", "2"),
				),
				Destroy: false,
			},
			{
				Config: testAccMongoDbUsersRemoveConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDbUsersNotExists(dbResourceName, TestAccProvider),
				),
			},
		},
	})
}

func testAccMongoDbUsersConfig(testMongodbName, testAuthority1, testAuthority2 string) string {
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
	id = ncloud_mongodb.mongodb.id
	mongodb_user_set = [
		{
			name = "testuser1",
			password = "t123456789!",
			database_name = "testdb1",
			authority = "%[2]s"
		},
		{
			name = "testuser2",
			password = "t123456789!",
			database_name = "testdb2",
			authority = "%[3]s"
		}
	]
}
`, testMongodbName, testAuthority1, testAuthority2)
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

func testAccCheckMongoDbUsersNotExists(n string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no MongoDB instance ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)

		users, err := mongodbservice.GetMongoDbUserAllList(context.Background(), config, resource.Primary.ID)
		if err != nil {
			return err
		}

		if users == nil {
			return nil
		}

		return fmt.Errorf("MongoDB users still exist: %v", users)
	}
}
