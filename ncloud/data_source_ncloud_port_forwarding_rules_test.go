package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudPortForwardingRulesBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPortForwardingRulesConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_port_forwarding_rules.rules"),
				//),
			},
		},
	})
}

var testAccDataSourceNcloudPortForwardingRulesConfig = `
data "ncloud_port_forwarding_rules" "rules" {
  zone = "KR-1"
}
`
