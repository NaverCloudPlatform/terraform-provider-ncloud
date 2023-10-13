package mysql_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mysqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/mysql"
	"regexp"
	"strings"
	"testing"
)

func TestAccResourceNcloudMysql_vpc_basic(t *testing.T) {
	var mysqlInstance vmysql.CloudMysqlInstance
	testMysqlName := fmt.Sprintf("tf-mysql-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql.mysql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMysqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlVpcConfig(testMysqlName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlExistsWithProvider(resourceName, &mysqlInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "service_name", testMysqlName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "testprefix"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "testusername"),
					resource.TestCheckResourceAttr(resourceName, "user_password", "t123456789!a"),
					resource.TestCheckResourceAttr(resourceName, "host_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "database_name", "test_db"),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_zone", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_storage_encryption", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup_file_retention_period", "1"),
				),
			},
		},
	})
}

func TestAccResourceNcloudMysql_vpc_isHa(t *testing.T) {
	var mysqlInstance vmysql.CloudMysqlInstance
	testMysqlName := fmt.Sprintf("tf-mysql-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql.mysql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMysqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlVpcConfigIsHa(testMysqlName, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlExistsWithProvider(resourceName, &mysqlInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_zone", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_storage_encryption", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_backup", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudMysql_vpc_isHa_options(t *testing.T) {
	var mysqlInstance vmysql.CloudMysqlInstance
	testMysqlName := fmt.Sprintf("tf-mysql-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql.mysql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMysqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlVpcConfigMultiZone(testMysqlName, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlExistsWithProvider(resourceName, &mysqlInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_zone", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_storage_encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_automatic_backup", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudMysql_vpc_auto_backup(t *testing.T) {
	var mysqlInstance vmysql.CloudMysqlInstance
	testMysqlName := fmt.Sprintf("tf-mysql-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql.mysql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMysqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlVpcConfigBackupWhenAuto(testMysqlName, false, true, 3, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlExistsWithProvider(resourceName, &mysqlInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup_file_retention_period", "3"),
					resource.TestCheckResourceAttr(resourceName, "is_automatic_backup", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudMysql_vpc_not_auto_backup(t *testing.T) {
	var mysqlInstance vmysql.CloudMysqlInstance
	testMysqlName := fmt.Sprintf("tf-mysql-%s", acctest.RandString(5))
	resourceName := "ncloud_mysql.mysql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMysqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlVpcConfigBackupWhenNotAuto(testMysqlName, false, true, 3, false, "11:15"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlExistsWithProvider(resourceName, &mysqlInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup_file_retention_period", "3"),
					resource.TestCheckResourceAttr(resourceName, "is_automatic_backup", "false"),
					resource.TestCheckResourceAttr(resourceName, "backup_time", "11:15"),
				),
			},
		},
	})
}

func testAccCheckMysqlExistsWithProvider(n string, mysql *vmysql.CloudMysqlInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		mysqlInstance, err := mysqlservice.GetMysqlInstance(context.Background(), config, resource.Primary.ID)
		if err != nil {
			return err
		}

		if mysqlInstance != nil {
			*mysql = *mysqlInstance
			return nil
		}

		return fmt.Errorf("mysql instance not found")
	}
}
func testAccCheckMysqlDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_mysql" {
			continue
		}
		instance, err := mysqlservice.GetMysqlInstance(context.Background(), config, rs.Primary.ID)
		if err != nil && !checkNoInstanceResponse(err) {
			return err
		}

		if instance != nil {
			return errors.New("mysql still exists")
		}
	}

	return nil
}

func checkNoInstanceResponse(err error) bool {
	return strings.Contains(err.Error(), "5001017")
}

func testAccMysqlVpcConfig(testMysqlName string) string {
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
	name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	host_ip = "192.168.0.1"
	database_name = "test_db"
	}
`, testMysqlName)
}

func testAccMysqlVpcConfigIsHa(testMysqlName string, isHa bool, isMultiZone bool, isStorageEncryption bool) string {
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
	name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	host_ip = "192.168.0.1"
	database_name = "test_db"
	is_ha = %[2]t
	is_multi_zone = %[3]t
	is_storage_encryption = %[4]t
}
`, testMysqlName, isHa, isMultiZone, isStorageEncryption)
}

func testAccMysqlVpcConfigMultiZone(testMysqlName string, isHa bool, isMultiZone bool, isStorageEncryption bool) string {
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

resource "ncloud_subnet" "test_subnet_standby" {
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
	name               = "%[1]s-standby"
	subnet             = "10.5.5.0/28"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}


resource "ncloud_mysql" "mysql" {
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	host_ip = "192.168.0.1"
	database_name = "test_db"

	is_ha = %[2]t
	is_multi_zone = %[3]t
	is_storage_encryption = %[4]t
	standby_master_subnet_no = ncloud_subnet.test_subnet_standby.id
}
`, testMysqlName, isHa, isMultiZone, isStorageEncryption)
}

func testAccMysqlVpcConfigBackupWhenAuto(name string, isHa bool, isBackup bool, backupPeriod int, isAutomaticBackup bool) string {
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
	name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	host_ip = "192.168.0.1"
	database_name = "test_db"

	is_ha = %[2]t
	
	is_backup = %[3]t
	is_automatic_backup = %[4]t
	backup_file_retention_period = %[5]d
}
`, name, isHa, isBackup, isAutomaticBackup, backupPeriod)
}

func testAccMysqlVpcConfigBackupWhenNotAuto(name string, isHa bool, isBackup bool, backupPeriod int, isAutomaticBackup bool, backupTime string) string {
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
	name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!"
	host_ip = "192.168.0.1"
	database_name = "test_db"

	is_ha = %[2]t
	
	is_backup = %[3]t
	is_automatic_backup = %[4]t
	backup_file_retention_period = %[5]d
	backup_time = "%[6]s"
}
`, name, isHa, isBackup, isAutomaticBackup, backupPeriod, backupTime)
}
