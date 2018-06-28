package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strconv"
	"testing"
)

func TestAccResourceNcloudPortForwardingRuleBasic(t *testing.T) {
	var portForwarding sdk.PortForwardingRule

	externalPort := acctest.RandIntRange(1024, 65534+1024)
	log.Printf("[DEBUG] externalPort: %d", externalPort)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_port_forwarding_rule.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPortForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardingRuleBasicConfig(externalPort),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortForwardingRuleExists("ncloud_port_forwarding_rule.test", &portForwarding),
					resource.TestCheckResourceAttr(
						"ncloud_port_forwarding_rule.test",
						"port_forwarding_external_port",
						strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(
						"ncloud_port_forwarding_rule.test",
						"port_forwarding_internal_port",
						"22"),
				),
			},
		},
	})
}

func testAccCheckPortForwardingRuleExists(n string, i *sdk.PortForwardingRule) resource.TestCheckFunc {
	return testAccCheckPortForwardingRuleExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckPortForwardingRuleExistsWithProvider(n string, i *sdk.PortForwardingRule, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		conn := provider.Meta().(*NcloudSdk).conn
		portForwardingRule, err := getPortForwardingRule(conn, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if portForwardingRule != nil {
			*i = *portForwardingRule
			return nil
		}

		return fmt.Errorf("port forwarding rule not found")
	}
}

func testAccCheckPortForwardingRuleDestroy(s *terraform.State) error {
	return testAccCheckPortForwardingRuleDestroyWithProvider(s, testAccProvider)
}

func testAccCheckPortForwardingRuleDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*NcloudSdk).conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_port_forwarding_rule" {
			continue
		}
		rule, err := getPortForwardingRule(conn, rs.Primary.ID)
		if rule == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if rule != nil {
			return fmt.Errorf("found not deleted resource: %s", rule.PortForwardingExternalPort)
		}
	}

	return nil
}

func testAccPortForwardingRuleBasicConfig(externalPort int) string {
	prefix := getTestPrefix()
	testServerName := prefix + "-vm"
	return fmt.Sprintf(`
				data "ncloud_port_forwarding_rules" "rules" {}

				resource "ncloud_server" "server" {
					"server_name" = "%s"
					"server_image_product_code" = "SPSW0LINUX000032"
					"server_product_code" = "SPSVRSTAND000004"
				}

			   resource "ncloud_port_forwarding_rule" "test" {
				   "port_forwarding_configuration_no" = "${data.ncloud_port_forwarding_rules.rules.id}"
				   "server_instance_no" = "${ncloud_server.server.id}"
				   "port_forwarding_external_port" = "%d"
				   "port_forwarding_internal_port" = "22"
			   }`, testServerName, externalPort)

}
