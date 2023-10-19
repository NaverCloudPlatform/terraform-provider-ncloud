package server_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_access_control_groups.test"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAccessControlGroupsDefault(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupsDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_access_control_groups.default"),
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
