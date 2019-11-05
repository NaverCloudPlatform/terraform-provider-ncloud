package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceNcloudPortForwardingRulesBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
