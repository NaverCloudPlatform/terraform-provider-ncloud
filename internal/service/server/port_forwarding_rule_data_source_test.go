package server

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudPortForwardingRuleBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPortForwardingRuleConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return SkipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	TestAccCheckDataSourceID("data.ncloud_port_forwarding_rule.test"),
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
