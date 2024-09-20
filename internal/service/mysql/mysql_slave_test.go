package mysql_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mysqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/mysql"
)

func TestAccrResourceNcloudMysqlSlave_vpc_basic(t *testing.T) {
	var mysqlServerInstance vmysql.CloudMysqlServerInstance
	testName := fmt.Sprintf("tf-mysqlsv-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql_slave.mysql_slave"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMysqlSlaveDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSlaveConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlSlaveExists(resourceName, &mysqlServerInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttrSet(resourceName, "mysql_instance_no"),
				),
			},
		},
	})
}

func testAccMysqlSlaveConfig(testName string) string {
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

resource "ncloud_mysql_slave" "mysql_slave" {
	mysql_instance_no = ncloud_mysql.mysql.id
}
`, testName)
}

func testAccCheckMysqlSlaveExists(n string, slave *vmysql.CloudMysqlServerInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		mysqlSlave, err := mysqlservice.GetMysqlSlave(context.Background(), config, resource.Primary.Attributes["mysql_instance_no"])
		if err != nil {
			return err
		}

		if mysqlSlave != nil {
			*slave = *mysqlSlave
			return nil
		}

		return fmt.Errorf("mysql slave not found")
	}
}

func testAccCheckMysqlSlaveDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_mysql_slave" {
			continue
		}
		instance, err := mysqlservice.GetMysqlSlave(context.Background(), config, rs.Primary.Attributes["mysql_instance_no"])
		if err != nil && !strings.Contains(err.Error(), "5001017") {
			return nil
		}

		if instance != nil {
			return errors.New("mysql slave still exists")
		}
	}

	return nil
}
