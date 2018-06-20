package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDataSourceNcloudAccessControlRuleBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceNcloudAccessControlRuleDataSourceID("data.ncloud_access_control_rule.test"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAccessControlRuleDataSourceID(n string) resource.TestCheckFunc {
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

var testAccDataSourceNcloudAccessControlRuleConfig = `
data "ncloud_access_control_rule" "test" {
  "is_default_group" = "true"
  "destination_port" = "22"
}
`
