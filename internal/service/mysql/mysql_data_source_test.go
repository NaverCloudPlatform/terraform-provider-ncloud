package mysql_test

import (
	"fmt"

	randacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"regexp"
	"testing"
)

func TestAccDataSourceNcloudMysql_vpc_basic(t *testing.T) {
	/*
		TODO - it's	for atomicity of regression testing. remove when error has solved.
	*/
	t.Skip()

	dataName := "data.ncloud_mysql.by_id"
	resourceName := "ncloud_mysql.mysql"
	testMysqlName := fmt.Sprintf("tf-mysql-%s", randacctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlConfig(testMysqlName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "service_name", resourceName, "service_name"),
					resource.TestCheckResourceAttrPair(dataName, "is_ha", resourceName, "is_ha"),
					resource.TestCheckResourceAttrPair(dataName, "is_multi_zone", resourceName, "is_multi_zone"),
					resource.TestCheckResourceAttrPair(dataName, "is_backup", resourceName, "is_backup"),
					resource.TestCheckResourceAttrPair(dataName, "backup_file_retention_period", resourceName, "backup_file_retention_period"),
				),
			},
		},
	})
}

func testAccDataSourceMysqlConfig(testMysqlName string) string {
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

resource "ncloud_mysql" "mysql" {
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	host_ip = "192.168.0.1"
	database_name = "test_db"
}

data "ncloud_mysql" "by_id" {
	id = ncloud_mysql.mysql.id
}
`, testMysqlName)
}
