package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudNetworkAclsBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudNetworkAclsConfig(),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_acls.all"),
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
				Config:   testAccDataSourceNcloudNetworkAclsConfigName("default"),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_acls.by_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNetworkAclsStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudNetworkAclsConfigStatus("RUN"),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_acls.by_status"),
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
				Config:   testAccDataSourceNcloudNetworkAclsConfigVpcNo(),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_network_acls.by_vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNetworkAclsConfig() string {
	return `
data "ncloud_network_acls" "all" {}
`
}

func testAccDataSourceNcloudNetworkAclsConfigName(name string) string {
	return fmt.Sprintf(`
data "ncloud_network_acls" "by_name" {
	name = "%s"
}
`, name)
}

func testAccDataSourceNcloudNetworkAclsConfigStatus(status string) string {
	return fmt.Sprintf(`
data "ncloud_network_acls" "by_status" {
	status = "%s"
}
`, status)
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

func testAccDataSourceNcloudNetworkAclsConfigNetworkACLNoList(networkACLNo string) string {
	return fmt.Sprintf(`
data "ncloud_network_acls" "by_network_acl_no" {
	network_acl_no_list = ["%s"]
}
`, networkACLNo)
}
