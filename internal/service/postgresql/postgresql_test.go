package postgresql_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	postgresqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/postgresql"
)

func TestAccResourceNcloudPostgresql_vpc_basic(t *testing.T) {
	var postgresqlInstance vpostgresql.CloudPostgresqlInstance
	testPostgresqlName := fmt.Sprintf("tf-postgresql-%s", acctest.RandString(5))
	resourceName := "ncloud_postgresql.postgresql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPostgresqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlConfig(testPostgresqlName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlExistsWithProvider(resourceName, &postgresqlInstance, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "service_name", testPostgresqlName),
					resource.TestCheckResourceAttr(resourceName, "ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "multi_zone", "false"),
					resource.TestCheckResourceAttr(resourceName, "backup", "true"),
				),
			},
		},
	})
}

// Available only `pub` and 'fin' site.
func TestAccResourceNcloudPostgresql_vpc_multizone(t *testing.T) {
	var postgresqlInstance vpostgresql.CloudPostgresqlInstance
	testPostgresqlName := fmt.Sprintf("tf-postgresql-%s", acctest.RandString(5))
	resourceName := "ncloud_postgresql.postgresql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPostgresqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlVpcConfigMultiZone(testPostgresqlName, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlExistsWithProvider(resourceName, &postgresqlInstance, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "multi_zone", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudPostgresql_vpc_error(t *testing.T) {
	testPostgresqlName := fmt.Sprintf("tf-postgresql-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPostgresqlDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPostgresqlVpcConfigErrorCaseBackupFalse(testPostgresqlName),
				ExpectError: regexp.MustCompile("when `ha` is true, `backup` must be true or not be input"),
			},
			{
				Config:      testAccPostgresqlVpcConfigErrorCaseBackupTimeSet(testPostgresqlName),
				ExpectError: regexp.MustCompile("when `backup` and `automatic_backup` is true, `backup_time` must not be"),
			},
			{
				Config:      testAccPostgresqlVpcConfigErrorCaseBackupTimeUnset(testPostgresqlName),
				ExpectError: regexp.MustCompile("when `backup` is true and `automatic_backup` is false, `backup_time` must be"),
			},
		},
	})
}

func testAccCheckPostgresqlExistsWithProvider(n string, postgresql *vpostgresql.CloudPostgresqlInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		postgresqlInstance, err := postgresqlservice.GetPostgresqlInstance(context.Background(), config, resource.Primary.ID)
		if err != nil {
			return err
		}

		if postgresqlInstance != nil {
			*postgresql = *postgresqlInstance
			return nil
		}

		return fmt.Errorf("postgresql instance not found")
	}
}

func testAccCheckPostgresqlDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_postgresql" {
			continue
		}
		instance, err := postgresqlservice.GetPostgresqlInstance(context.Background(), config, rs.Primary.ID)
		if err != nil && !checkNoInstanceResponse(err) {
			return err
		}

		if instance != nil {
			return errors.New("postgresql still exists")
		}
	}

	return nil
}

func checkNoInstanceResponse(err error) bool {
	return strings.Contains(err.Error(), "5001017")
}

func testAccPostgresqlConfig(testPostgresqlName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test_vpc" {
	name         = "%[1]s"
	ipv4_cidr_block = "10.5.0.0/16"
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
	data_storage_type = "SSD"
	backup = true
	backup_file_retention_period = 2
	backup_time = "02:00"
	backup_file_storage_count = 3
	backup_file_compression = true
	automatic_backup = false
	port = 5432
}
`, testPostgresqlName)
}

func testAccPostgresqlVpcConfigMultiZone(testPostgresqlName string, isHa bool, isMultiZone bool, isStorageEncryption bool) string {
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
resource "ncloud_subnet" "test_secondary_subnet" {
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
	name               = "%[1]s-secondary"
	subnet             = "10.5.5.0/28"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}
resource "ncloud_postgresql" "postgresql" {
	vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	client_cidr = "0.0.0.0/0"
	database_name = "test_db"
	ha = %[2]t
	multi_zone = %[3]t
	storage_encryption = %[4]t
	secondary_subnet_no = ncloud_subnet.test_secondary_subnet.id
}
`, testPostgresqlName, isHa, isMultiZone, isStorageEncryption)
}

func testAccPostgresqlVpcConfigBase(name string) string {
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
`, name)
}

func testAccPostgresqlVpcConfigErrorCaseBackupFalse(name string) string {
	return testAccPostgresqlVpcConfigBase(name) + fmt.Sprintf(`
resource "ncloud_postgresql" "postgresql" {
	vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	client_cidr = "0.0.0.0/0"
	database_name = "test_db"
	backup = false
}
`, name)
}

func testAccPostgresqlVpcConfigErrorCaseBackupTimeSet(name string) string {
	return testAccPostgresqlVpcConfigBase(name) + fmt.Sprintf(`
resource "ncloud_postgresql" "postgresql" {
	vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	client_cidr = "0.0.0.0/0"
	database_name = "test_db"
	ha = true
	backup = true
	automatic_backup = true
	backup_time = "01:30"
}
`, name)
}

func testAccPostgresqlVpcConfigErrorCaseBackupTimeUnset(name string) string {
	return testAccPostgresqlVpcConfigBase(name) + fmt.Sprintf(`
resource "ncloud_postgresql" "postgresql" {
	vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	client_cidr = "0.0.0.0/0"
	database_name = "test_db"
	ha = true
	backup = true 
	automatic_backup = false
}
`, name)
}
