package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMongoDbUsers_vpc_basic(t *testing.T) {
	testName := fmt.Sprintf("tf-mduser-%s", acctest.RandString(3))
	dataName := "data.ncloud_mongodb_users.all"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongodbUsersConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "mongodb_user_list.0.name", "test1"),
					resource.TestCheckResourceAttr(dataName, "mongodb_user_list.1.name", "test2"),
				),
			},
		},
	})
}

func testAccDataSourceMongodbUsersConfig(testName string) string {
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
			name = "test1",
			password = "t123456789!",
			database_name = "testdb1",
			authority = "READ"
		},
		{
			name = "test2",
			password = "t123456789!",
			database_name = "testdb2",
			authority = "READ_WRITE"
		}
	]
}

data "ncloud_mongodb_users" "all" {
	mongodb_instance_no = ncloud_mongodb.mongodb.id
	filter {
		name = "name"
		values = [ncloud_mongodb_users.mongodb_users.mongodb_user_list.0.name, ncloud_mongodb_users.mongodb_users.mongodb_user_list.1.name]
	}
}
`, testName)
}
