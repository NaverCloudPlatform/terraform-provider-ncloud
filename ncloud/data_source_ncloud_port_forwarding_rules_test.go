package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudPortForwardingRulesBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPortForwardingRulesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_port_forwarding_rules.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudPortForwardingRulesConfig = `
data "ncloud_port_forwarding_rules" "test" {}
`
