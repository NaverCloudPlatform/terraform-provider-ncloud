package server_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudAccessControlRuleBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_access_control_rule.test"),
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
