package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudAccessControlGroupsBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func TestAccDataSourceNcloudAccessControlGroupsDefault(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupsDefaultConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
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
