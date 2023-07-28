package cloudmysql_test

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mysqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/cloudmysql"
)

func TestAccResourceNcloudMysql_vpc_basic(t *testing.T) {
	//var mysqlInstance vmysql.CloudMysqlInstance
	name := fmt.Sprintf("test-vpc-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		CheckDestroy: testAccCheckMysqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlConfig(name),
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckMysqlExistsWithProvider(resourceName, &mysqlInstance),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

func testAccCheckMysqlExistsWithProvider(n string, mysql *vmysql.CloudMysqlInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		mysqlInstance, err := mysqlservice.GetMysqlInstance(config, resource.Primary.ID)
		if err != nil {
			return err
		}

		*mysql = *mysqlInstance

		return fmt.Errorf("server instance not found")
	}
}
func testAccCheckMysqlDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_mysql" {
			continue
		}

		instance, err := mysqlservice.GetMysqlInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("mysql still exists")
		}
	}

	return nil
}

func testAccDataSourceMysqlConfig(testMysqlrName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_mysql" "mysql" {
	subnet_no = ncloud_subnet.test.id
	service_name = "testservice"
	name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	host_ip = "192.168.0.1"
	database_name = "test_db"
}
`, testMysqlrName)
}
