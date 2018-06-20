package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					testAccCheckNcloudPortForwardingRuleDataSourceID("data.ncloud_port_forwarding_rule.test"),
				),
			},
		},
	})
}

func testAccCheckNcloudPortForwardingRuleDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudPortForwardingRuleConfig = `
data "ncloud_port_forwarding_rule" "test" {
  // "server_instance_no" = "690731"
  "port_forwarding_external_port" = "4088"
  //"port_forwarding_internal_port" = "22"
}
`
