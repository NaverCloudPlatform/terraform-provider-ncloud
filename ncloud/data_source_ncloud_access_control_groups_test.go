package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					testAccDataSourceNcloudAccessControlGroupsDataSourceID("data.ncloud_access_control_groups.test"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAccessControlGroupsDataSourceID(n string) resource.TestCheckFunc {
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

var testAccDataSourceNcloudAccessControlGroupsConfig = `
data "ncloud_access_control_groups" "test" {}
`
