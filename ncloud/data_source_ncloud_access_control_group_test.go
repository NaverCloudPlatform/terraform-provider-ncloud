package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					testAccDataSourceNcloudAccessControlGroupDataSourceID("data.ncloud_access_control_group.test"),
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
					testAccDataSourceNcloudAccessControlGroupDataSourceID("data.ncloud_access_control_group.test"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAccessControlGroupDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudAccessControlGroupConfig = `
data "ncloud_access_control_group" "test" {
	"access_control_group_name" = "Default"
}
`

var testAccDataSourceNcloudAccessControlGroupMostRecentConfig = `
data "ncloud_access_control_group" "test" {
	"most_recent" = "true"
}
`
