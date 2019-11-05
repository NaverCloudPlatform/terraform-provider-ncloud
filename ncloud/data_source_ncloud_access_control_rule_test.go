package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
