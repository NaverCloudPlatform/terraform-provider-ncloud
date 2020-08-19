package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudVpcsBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigName("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.by_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigStatus("RUN"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.by_status"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsVpcNo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigVpcNo("446"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.by_vpc_no"),
					testAccCheckDataSourceID("data.ncloud_vpcs.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudVpcsConfig() string {
	return fmt.Sprintf(`
data "ncloud_vpcs" "all" {}
`)
}

func testAccDataSourceNcloudVpcsConfigName(name string) string {
	return fmt.Sprintf(`
data "ncloud_vpcs" "by_name" {
	name               = "%s"
}
`, name)
}

func testAccDataSourceNcloudVpcsConfigStatus(status string) string {
	return fmt.Sprintf(`
data "ncloud_vpcs" "by_status" {
	status               = "%s"
}
`, status)
}

func testAccDataSourceNcloudVpcsConfigVpcNo(vpcNo string) string {
	return fmt.Sprintf(`
data "ncloud_vpcs" "by_vpc_no" {
	vpc_no          = "%[1]s"
}

data "ncloud_vpcs" "by_filter" {
	filter {
		name   = "vpc_no"
		values = ["%[1]s"]
	}
}
`, vpcNo)
}
