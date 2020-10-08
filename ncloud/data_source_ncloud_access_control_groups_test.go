package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudAccessControlGroups_classic_basic(t *testing.T) {
	testAccDataSourceNcloudAccessControlGroupsBasic(t, false)
}

func TestAccDataSourceNcloudAccessControlGroups_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudAccessControlGroupsBasic(t, true)
}

func TestAccDataSourceNcloudAccessControlGroups_classic_default(t *testing.T) {
	testAccDataSourceNcloudAccessControlGroupsDefault(t, false)
}

func TestAccDataSourceNcloudAccessControlGroups_vpc_default(t *testing.T) {
	testAccDataSourceNcloudAccessControlGroupsDefault(t, true)
}

func testAccDataSourceNcloudAccessControlGroupsBasic(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_groups.test"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAccessControlGroupsDefault(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupsDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_groups.default"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudAccessControlGroupsConfig = `
data "ncloud_access_control_groups" "test" {}
`

var testAccDataSourceNcloudAccessControlGroupsDefaultConfig = `
data "ncloud_access_control_groups" "default" {
  is_default_group = "true"
}
`
