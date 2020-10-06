package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudSubnetsBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudSubnetsConfig(),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_subnets.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSubnetsName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudSubnetsConfigSubnet("10.2.1.0"),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_subnets.by_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSubnetsStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudSubnetsConfigStatus("RUN"),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_subnets.by_status"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSubnetsVpcNo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudSubnetsConfigVpcNo("502"),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_subnets.by_vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSubnetsConfig() string {
	return `
data "ncloud_subnets" "all" {}
`
}

func testAccDataSourceNcloudSubnetsConfigSubnet(cidr string) string {
	return fmt.Sprintf(`
data "ncloud_subnets" "by_cidr" {
	subnet = "%s"
}
`, cidr)
}

func testAccDataSourceNcloudSubnetsConfigStatus(status string) string {
	return fmt.Sprintf(`
data "ncloud_subnets" "by_status" {
	status = "%s"
}
`, status)
}

func testAccDataSourceNcloudSubnetsConfigVpcNo(vpcNo string) string {
	return fmt.Sprintf(`
data "ncloud_subnets" "by_vpc_no" {
	vpc_no = "%s"
}
`, vpcNo)
}
