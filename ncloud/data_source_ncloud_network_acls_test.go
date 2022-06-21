package ncloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNetworkAclsBasic(t *testing.T) {
	dataSourceName := "data.ncloud_network_acls.all"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkAclsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataSourceName),
					resource.TestMatchResourceAttr(dataSourceName, "network_acls.#", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "network_acls.0.id", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "network_acls.0.vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "network_acls.0.network_acl_no", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNetworkAclsName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkAclsConfigName("default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_acls.by_name"),
					testAccCheckDataSourceID("data.ncloud_network_acls.by_filter"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNetworkAclsVpcNo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkAclsConfigVpcNo(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_acls.by_vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNetworkAclsConfig() string {
	return `
resource "ncloud_vpc" "test" {
	name               = "testacc-data-network-acl"
	ipv4_cidr_block    = "10.2.0.0/16"
}

data "ncloud_network_acls" "all" {}
`
}

func testAccDataSourceNcloudNetworkAclsConfigName(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "testacc-data-network-acl"
	ipv4_cidr_block    = "10.2.0.0/16"
}

data "ncloud_network_acls" "by_name" {
	name = "%[1]s"
}

data "ncloud_network_acls" "by_filter" {
	filter {
		name   = "name"
		values = ["%[1]s"]
	}
}
`, name)
}

func testAccDataSourceNcloudNetworkAclsConfigVpcNo() string {
	return `
resource "ncloud_vpc" "test" {
	name               = "testacc-data-network-acl"
	ipv4_cidr_block    = "10.2.0.0/16"
}

data "ncloud_network_acls" "by_vpc_no" {
	vpc_no = ncloud_vpc.test.vpc_no
}
`
}
