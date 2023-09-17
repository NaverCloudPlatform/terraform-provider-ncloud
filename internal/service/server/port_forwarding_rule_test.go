package server_test

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	serverservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccResourceNcloudPortForwardingRuleBasic(t *testing.T) {
	var portForwarding server.PortForwardingRule

	externalPort := int(generateExternalPort(1024, 65534)) // acctest.RandIntRange(1024, 65534+1024)
	log.Printf("[DEBUG] externalPort: %d", externalPort)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPortForwardingRuleDestroy,
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

func generateExternalPort(min, max int32) int32 {
	rand.Seed(time.Now().Unix())
	return rand.Int31n(max-min) + min
}

// TODO: ignore test: may be empty created data
func ignore_TestAccResourceNcloudPortForwardingRuleExistingServer(t *testing.T) {
	var portForwarding server.PortForwardingRule

	externalPort := acctest.RandIntRange(1024, 65534+1024)
	log.Printf("[DEBUG] externalPort: %d", externalPort)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPortForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardingRuleExistingServerConfig(externalPort),
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return SkipNoResultsTest, nil
				},
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
	return testAccCheckPortForwardingRuleExistsWithProvider(n, i, func() *schema.Provider { return GetTestProvider(false) })
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
		client := provider.Meta().(*conn.ProviderConfig).Client
		_, zoneNo, portForwardingExternalPort := serverservice.ParsePortForwardingRuleId(rs.Primary.ID)
		portForwardingRule, err := serverservice.GetPortForwardingRule(client, zoneNo, portForwardingExternalPort)
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
	return testAccCheckPortForwardingRuleDestroyWithProvider(s, GetTestProvider(false))
}

func testAccCheckPortForwardingRuleDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*conn.ProviderConfig).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_port_forwarding_rule" {
			continue
		}
		_, zoneNo, portForwardingExternalPort := serverservice.ParsePortForwardingRuleId(rs.Primary.ID)
		rule, err := serverservice.GetPortForwardingRule(client, zoneNo, portForwardingExternalPort)
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
	prefix := GetTestPrefix()
	testServerName := prefix + "-vm"
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%s-key"
}

resource "ncloud_server" "server" {
	name = "%s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_port_forwarding_rule" "test" {
	server_instance_no = "${ncloud_server.server.id}"
	port_forwarding_external_port = "%d"
	port_forwarding_internal_port = "22"
}`, testServerName, testServerName, externalPort)
}

func testAccPortForwardingRuleExistingServerConfig(externalPort int) string {
	return fmt.Sprintf(`
resource "ncloud_port_forwarding_rule" "test" {
	server_instance_no = "966669"
	port_forwarding_external_port = "%d"
	port_forwarding_internal_port = "22"
}`, externalPort)
}
