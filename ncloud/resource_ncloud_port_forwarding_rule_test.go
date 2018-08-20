package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strconv"
	"testing"
)

func TestAccResourceNcloudPortForwardingRuleBasic(t *testing.T) {
	var portForwarding server.PortForwardingRule

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

func testAccCheckPortForwardingRuleExists(n string, i *server.PortForwardingRule) resource.TestCheckFunc {
	return testAccCheckPortForwardingRuleExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckPortForwardingRuleExistsWithProvider(n string, i *server.PortForwardingRule, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		client := provider.Meta().(*NcloudAPIClient)
		_, zoneNo, portForwardingExternalPort := parsePortForwardingRuleId(rs.Primary.ID)
		portForwardingRule, err := getPortForwardingRule(client, zoneNo, portForwardingExternalPort)
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
	client := provider.Meta().(*NcloudAPIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_port_forwarding_rule" {
			continue
		}
		_, zoneNo, portForwardingExternalPort := parsePortForwardingRuleId(rs.Primary.ID)
		rule, err := getPortForwardingRule(client, zoneNo, portForwardingExternalPort)
		if rule == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if rule != nil {
			return fmt.Errorf("found not deleted resource: %d", *rule.PortForwardingExternalPort)
		}
	}

	return nil
}

func testAccPortForwardingRuleBasicConfig(externalPort int) string {
	prefix := getTestPrefix()
	testServerName := prefix + "-vm"
	return fmt.Sprintf(`
				resource "ncloud_server" "server" {
					"server_name" = "%s"
					"server_image_product_code" = "SPSW0LINUX000032"
					"server_product_code" = "SPSVRSTAND000004"
				}

			   resource "ncloud_port_forwarding_rule" "test" {
				   "server_instance_no" = "${ncloud_server.server.id}"
				   "port_forwarding_external_port" = "%d"
				   "port_forwarding_internal_port" = "22"
			   }`, testServerName, externalPort)

}
