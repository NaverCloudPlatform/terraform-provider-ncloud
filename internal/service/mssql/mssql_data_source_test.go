package cloudmssql_test

import (
	"fmt"

	randacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"testing"
)

func TestAccDataSourceNcloudMssql_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_mssql.test"
	resourceName := "ncloud_mssql.mssql"
	testMssqlName := fmt.Sprintf("tf-mssql-%s", randacctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMssqlConfig(testMssqlName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "service_name", resourceName, "service_name"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
				),
			},
		},
	})
}

func testAccDataSourceMssqlConfig(testMssqlName string) string {
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
			vpc_no             = ncloud_vpc.test_vpc.vpc_no
			subnet_no = ncloud_subnet.test_subnet.id
			service_name = "%[1]s"
			is_ha = true
			is_multi_zone = false
			is_automatic_backup = true
			user_name = "test"
			user_password = "qwer1234!"
		}
		data "ncloud_mssql" "test" {
			service_name = ncloud_mssql.mssql.service_name
			vpc_no = ncloud_mssql.mssql.vpc_no
			subnet_no = ncloud_mssql.mssql.subnet_no
		}
	`, testMssqlName)
}
