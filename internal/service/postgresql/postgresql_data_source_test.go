package postgresql_test

import (
	"fmt"
	"regexp"
	"testing"

	randacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudPostgresql_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_postgresql.by_id"
	resourceName := "ncloud_postgresql.postgresql"
	testPostgresqlName := fmt.Sprintf("tf-postgresql-%s", randacctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePostgresqlConfig(testPostgresqlName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "service_name", resourceName, "service_name"),
					resource.TestCheckResourceAttrPair(dataName, "ha", resourceName, "ha"),
					resource.TestCheckResourceAttrPair(dataName, "multi_zone", resourceName, "multi_zone"),
					resource.TestCheckResourceAttrPair(dataName, "backup", resourceName, "backup"),
					resource.TestCheckResourceAttrPair(dataName, "backup_file_retention_period", resourceName, "backup_file_retention_period"),
				),
			},
		},
	})
}

func testAccDataSourcePostgresqlConfig(testPostgresqlName string) string {
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
    vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	client_cidr = "0.0.0.0/0"
	database_name = "test_db"
}
data "ncloud_postgresql" "by_id" {
	id = ncloud_postgresql.postgresql.id
}
`, testPostgresqlName)
}
