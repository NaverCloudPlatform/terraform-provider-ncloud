package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_port_forwarding_rule.test"),
				//),
			},
		},
	})
}

var testAccDataSourceNcloudPortForwardingRuleConfig = `
data "ncloud_port_forwarding_rule" "test" {
  zone_code = "KR-2"
}
`
