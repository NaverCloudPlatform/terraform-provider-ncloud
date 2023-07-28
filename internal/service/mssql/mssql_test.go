package cloudmssql_test

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmssql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mssqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/mssql"
)

func TestAccResourceNcloudMssql_vpc_basic(t *testing.T) {
	var mssqlInstance vmssql.CloudMssqlInstance
	testMssqlName := fmt.Sprintf("tf-mssql-%s", acctest.RandString(5))
	resourceName := "ncloud_mssql.mssql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckCloudMssqlDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCloudMssqlVpcConfig(testMssqlName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudMssqlExists(resourceName, &mssqlInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "service_name", testMssqlName),
					resource.TestCheckResourceAttr(resourceName, "user_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "user_password", "qwer1234!"),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_zone", "false"),
					resource.TestCheckResourceAttr(resourceName, "backup_file_retention_period", "1"),
					resource.TestCheckResourceAttr(resourceName, "is_automatic_backup", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudMssql_vpc_isMultiZone(t *testing.T) {
	var mssqlInstance vmssql.CloudMssqlInstance
	testMssqlName := fmt.Sprintf("tf-mssql-%s", acctest.RandString(5))
	resourceName := "ncloud_mssql.mssql"
	testMirrorSubnetName := fmt.Sprintf("tf-mssql-%s", acctest.RandString(5))
	resourceNameMirrorSubnet := "ncloud_subnet.test_subnet_mirror"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckCloudMssqlDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCloudMssqlVpcConfigIsMultiZone(testMssqlName, testMirrorSubnetName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudMssqlExists(resourceName, &mssqlInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "service_name", testMssqlName),
					resource.TestCheckResourceAttr(resourceName, "user_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "user_password", "qwer1234!"),
					resource.TestCheckResourceAttr(resourceName, "is_ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_zone", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup_file_retention_period", "1"),
					resource.TestCheckResourceAttr(resourceName, "is_automatic_backup", "true"),
					resource.TestCheckResourceAttr(resourceNameMirrorSubnet, "name", testMirrorSubnetName),
				),
			},
		},
	})
}

func testAccCheckCloudMssqlExists(n string, mssql *vmssql.CloudMssqlInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		mssqlInstance, err := mssqlservice.GetCloudMssqlInstance(config, resource.Primary.ID)
		if err != nil {
			return err
		}

		if mssqlInstance == nil {
			return fmt.Errorf("Not found Mssql : %s", resource.Primary.ID)
		}

		*mssql = *mssqlInstance
		return nil
	}
}

func testAccCheckCloudMssqlDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_mssql" {
			continue
		}

		cloudMssql, err := mssqlservice.GetCloudMssqlInstance(config, rs.Primary.ID)
		if err != nil {
			commonErr, parseErr := GetCommonErrorBody(err)
			if parseErr == nil && commonErr.ReturnCode == "5001269" {
				return nil
			}
			return err
		}

		if cloudMssql != nil {
			return fmt.Errorf("CloudMssql(%s) still exists", ncloud.StringValue(cloudMssql.CloudMssqlInstanceNo))
		}
	}

	return nil
}

func testAccCloudMssqlVpcConfig(testMssqlName string) string {
	return fmt.Sprintf(`
		resource "ncloud_vpc" "test_vpc" {
			name               = "%[1]s"
			ipv4_cidr_block    = "10.0.0.0/16"
		}
		resource "ncloud_subnet" "test_subnet" {
			vpc_no             = ncloud_vpc.test_vpc.vpc_no
			name               = "%[1]s"
			subnet             = "10.0.0.0/24"
			zone               = "KR-2"
			network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
			subnet_type        = "PUBLIC"
		}
		resource "ncloud_mssql" "mssql" {
			vpc_no = ncloud_vpc.test_vpc.vpc_no
			subnet_no = ncloud_subnet.test_subnet.id
			service_name = "%[1]s"
			is_ha = true
			is_multi_zone = false
			is_automatic_backup = true
			user_name = "test"
			user_password = "qwer1234!"
		}
`, testMssqlName)
}

func testAccCloudMssqlVpcConfigIsMultiZone(testMssqlName string, testMirrorSubnetName string, isMultiZone bool) string {
	return fmt.Sprintf(`
		resource "ncloud_vpc" "test_vpc" {
			name               = "%[1]s"
			ipv4_cidr_block    = "10.0.0.0/16"
		}
		resource "ncloud_subnet" "test_subnet" {
			vpc_no             = ncloud_vpc.test_vpc.vpc_no
			name               = "%[1]s"
			subnet             = "10.0.0.0/24"
			zone               = "KR-2"
			network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
			subnet_type        = "PUBLIC"
		}
		resource "ncloud_subnet" "test_subnet_mirror" {
			vpc_no             = ncloud_vpc.test_vpc.vpc_no
			name               = "%[2]s"
			subnet             = "10.0.1.0/24"
			zone               = "KR-1"
			network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
			subnet_type        = "PUBLIC"
		}
		resource "ncloud_mssql" "mssql" {
			vpc_no = ncloud_vpc.test_vpc.vpc_no
			subnet_no = ncloud_subnet.test_subnet.id
			service_name = "%[1]s"
			is_ha = true
			is_multi_zone = %[3]t
			mirror_subnet_no = ncloud_subnet.test_subnet_mirror.id
			is_automatic_backup = true
			user_name = "test"
			user_password = "qwer1234!"
		}
`, testMssqlName, testMirrorSubnetName, true)
}
