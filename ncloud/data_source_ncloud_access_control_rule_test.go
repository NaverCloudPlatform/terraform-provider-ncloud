package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudAccessControlRuleBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_access_control_rule.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudAccessControlRuleConfig = `
data "ncloud_access_control_rule" "test" {
  is_default_group = "true"
  destination_port = "22"
}
`
