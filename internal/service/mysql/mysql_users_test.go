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

func TestAccResourceNcloudMysqlUsers_vpc_basic_update(t *testing.T) {
	testName := fmt.Sprintf("tf-mysqluser-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql_users.mysql_users"
	testUserBefore := "test"
	testUserAfter := "testuser"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMysqlUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlUsersConfig(testName, testUserBefore),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.0.name", "test1"),
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.0.host_ip", "%"),
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.0.authority", "READ"),
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.1.name", "test2"),
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.1.host_ip", "%"),
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.1.authority", "DDL"),
				),
			},
			{
				Config: testAccMysqlUsersConfig(testName, testUserAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.0.name", "testuser1"),
					resource.TestCheckResourceAttr(resourceName, "mysql_user_list.1.name", "testuser2"),
				),
			},
		},
	})
}

func testAccMysqlUsersConfig(testMysqlName string, testUser string) string {
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
			name = "%[2]s1",
			password = "t123456789!",
			host_ip = "%%",
			authority = "READ"
		},
		{
			name = "%[2]s2",
			password = "t123456789!",
			host_ip = "%%",
			authority = "DDL"
		}
	]
}
`, testMysqlName, testUser)
}

func testAccCheckMysqlUsersDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_mysql_users" {
			continue
		}

		instance, err := mysqlservice.GetMysqlUserList(context.Background(), config, rs.Primary.ID, []string{"testuser1", "testuser2"})
		if err != nil && !mysqlservice.CheckIfAlreadyDeleted(err) {
			return err
		}

		if len(instance) > 1 {
			return errors.New("mysql users still exists")
		}
	}

	return nil
}
