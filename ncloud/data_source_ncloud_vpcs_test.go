package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudVpcsBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfig(),
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigName("test"),
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.by_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsStatus(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigStatus("RUN"),
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.by_status"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsVpcNo(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigVpcNoList("446"),
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_vpcs.by_vpc_no"),
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

func testAccDataSourceNcloudVpcsConfigVpcNoList(vpcNo string) string {
	return fmt.Sprintf(`
data "ncloud_vpcs" "by_vpc_no" {
	vpc_no_list          = ["%s"]
}
`, vpcNo)
}
