package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudAccessControlGroupBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_group.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAccessControlGroupMostRecent(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlGroupMostRecentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_group.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudAccessControlGroupConfig = `
data "ncloud_access_control_group" "test" {
	//"access_control_group_name" = "Default"
	"access_control_group_name" = "winrm-acg"
}
`

var testAccDataSourceNcloudAccessControlGroupMostRecentConfig = `
data "ncloud_access_control_group" "test" {
	"most_recent" = "true"
}
`
