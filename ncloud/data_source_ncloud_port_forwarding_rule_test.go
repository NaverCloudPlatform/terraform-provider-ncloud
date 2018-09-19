package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceNcloudPortForwardingRuleBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceNcloudPortForwardingRuleConfig,
				ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_port_forwarding_rule.test"),
				//),
			},
		},
	})
}

var testAccDataSourceNcloudPortForwardingRuleConfig = `
data "ncloud_port_forwarding_rule" "test" {
  "zone_code" = "KR-1"
}
`
