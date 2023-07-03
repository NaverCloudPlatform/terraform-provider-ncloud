package server_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudPortForwardingRulesBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPortForwardingRulesConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return SkipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	TestAccCheckDataSourceID("data.ncloud_port_forwarding_rules.rules"),
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
