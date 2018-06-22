package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudPortForwardingRuleBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPortForwardingRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_port_forwarding_rule.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudPortForwardingRuleConfig = `
data "ncloud_port_forwarding_rule" "test" {
  // "server_instance_no" = "690731"
  "port_forwarding_external_port" = "4088"
  //"port_forwarding_internal_port" = "22"
}
`
