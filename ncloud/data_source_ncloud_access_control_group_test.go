package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

// ignore:  no results. please change search criteria and try again
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
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_group.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudAccessControlGroupConfig = `
data "ncloud_access_control_group" "test" {
	//"access_control_group_name" = "ncloud-default-acg"
	//"access_control_group_name" = "winrm-acg"
	"is_default_group" = true
	"most_recent" = "true"
}
`

var testAccDataSourceNcloudAccessControlGroupMostRecentConfig = `
data "ncloud_access_control_group" "test" {
	"most_recent" = "true"
}
`
