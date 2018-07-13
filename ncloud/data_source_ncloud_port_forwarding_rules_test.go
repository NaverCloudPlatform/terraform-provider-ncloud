package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
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
					testAccCheckDataSourceID("data.ncloud_port_forwarding_rules.rules"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudPortForwardingRules_RequiredZoneParam(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceNcloudPortForwardingRulesRequiredZoneParamConfig,
				ExpectError: regexp.MustCompile("required to select one among two parameters: `zone_no` and `zone_code`"),
			},
		},
	})
}

var testAccDataSourceNcloudPortForwardingRulesConfig = `
data "ncloud_port_forwarding_rules" "rules" {
  "zone_code" = "KR-1"
}
`

var testAccDataSourceNcloudPortForwardingRulesRequiredZoneParamConfig = `
data "ncloud_port_forwarding_rules" "rules" {
}
`
