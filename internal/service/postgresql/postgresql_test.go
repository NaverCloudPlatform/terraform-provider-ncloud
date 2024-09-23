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
					testAccCheckPostgresqlExistsWithProvider(resourceName, &postgresqlInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "service_name", testPostgresqlName),
					resource.TestCheckResourceAttr(resourceName, "server_name_prefix", "testprefix"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "testusername"),
					resource.TestCheckResourceAttr(resourceName, "user_password", "t123456789!a"),
					resource.TestCheckResourceAttr(resourceName, "client_cidr", "0.0.0.0/0"),
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

func TestAccResourceNcloudPostgresql_vpc_isHa(t *testing.T) {
	var postgresqlInstance vpostgresql.CloudPostgresqlInstance
	testPostgresqlName := fmt.Sprintf("tf-posrgresql-%s", acctest.RandString(5))
	resourceName := "ncloud_postgresql.postgresql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPostgresqlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlVpcConfigIsHa(testPostgresqlName, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlExistsWithProvider(resourceName, &postgresqlInstance, GetTestProvider(true)),
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

func TestAccResourceNcloudPostgresql_vpc_isHa_options(t *testing.T) {
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
					testAccCheckPostgresqlExistsWithProvider(resourceName, &postgresqlInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_zone", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_storage_encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_backup", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudPostgresql_error_case(t *testing.T) {
	testPostgresqlName := fmt.Sprintf("tf-postgresql-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPostgresqlDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPostgresqlVpcConfigErrorCaseWhenIsHaSetFalse1(testPostgresqlName),
				ExpectError: regexp.MustCompile("when `is_ha` is false, `is_multi_zone` parameter is not used"),
			},
			{
				Config:      testAccPostgresqlVpcConfigErrorCaseWhenIsHaSetFalse2(testPostgresqlName),
				ExpectError: regexp.MustCompile("when `is_ha` is false, `secondary_subnet_no` is not used"),
			},
			{
				Config:      testAccPostgresqlVpcConfigErrorCaseWhenIsHaSetTrue(testPostgresqlName),
				ExpectError: regexp.MustCompile("when `is_ha` is true, `is_backup` must be true or not be input"),
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
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

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
}
`, testPostgresqlName)
}

func testAccPostgresqlVpcConfigIsHa(testPostgresqlName string, isHa bool, isMultiZone bool, isStorageEncryption bool) string {
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
	user_password = "t123456789!"
	client_cidr = "0.0.0.0/0"
	database_name = "test_db"
	is_ha = %[2]t
	is_multi_zone = %[3]t
	is_storage_encryption = %[4]t
}
`, testPostgresqlName, isHa, isMultiZone, isStorageEncryption)
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
	is_ha = %[2]t
	is_multi_zone = %[3]t
	is_storage_encryption = %[4]t
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

func testAccPostgresqlVpcConfigErrorCaseWhenIsHaSetFalse1(name string) string {
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
	is_ha = false
	is_multi_zone = true
}
`, name)
}

func testAccPostgresqlVpcConfigErrorCaseWhenIsHaSetFalse2(name string) string {
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
	is_ha = false
	secondary_subnet_no = "12346" 
}
`, name)
}

func testAccPostgresqlVpcConfigErrorCaseWhenIsHaSetTrue(name string) string {
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
	is_ha = true
	is_backup = false
}
`, name)
}
