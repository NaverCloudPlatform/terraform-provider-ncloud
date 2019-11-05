package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// ignore:  no results. please change search criteria and try again
func TestAccDataSourceNcloudAccessControlGroupBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_group.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAccessControlGroupSelectByName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupSelectByNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_group.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudAccessControlGroupDefaultConfig = `
data "ncloud_access_control_group" "test" {
	is_default_group = true
}
`

var testAccDataSourceNcloudAccessControlGroupSelectByNameConfig = `
data "ncloud_access_control_group" "test" {
	name = "ncloud-default-acg"
}
`
