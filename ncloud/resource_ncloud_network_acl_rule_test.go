package ncloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudNetworkACLRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLRuleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_inbound_100"),
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_inbound_110"),
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_inbound_120"),
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_outbound_100"),
				),
			},
			{
				ResourceName:      "ncloud_network_acl_rule.nacl_rule_inbound_100",
				ImportState:       true,
				ImportStateIdFunc: testAccNcloudNetworkACLRuleImportStateIDFunc("ncloud_network_acl_rule.nacl_rule_inbound_100"),
				ImportStateVerify: true,
			},
			{
				ResourceName:      "ncloud_network_acl_rule.nacl_rule_outbound_100",
				ImportState:       true,
				ImportStateIdFunc: testAccNcloudNetworkACLRuleImportStateIDFunc("ncloud_network_acl_rule.nacl_rule_outbound_100"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceNcloudNetworkACLRuleConfig() string {
	return `
resource "ncloud_vpc" "vpc" {
	name            = "tf-testacc-network-acl-rule"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "nacl"
	description = "tf-testacc-network-acl-rule"
}

resource "ncloud_network_acl_rule" "nacl_rule_inbound_100" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no
	priority          = 100
	protocol          = "TCP"
	rule_action       = "ALLOW"
	port_range        = "22"
	ip_block          = "0.0.0.0/0"
	network_rule_type = "INBND"
	description       = "tf-testacc-network-acl-rule"
}

resource "ncloud_network_acl_rule" "nacl_rule_inbound_110" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no
	priority          = 110
	protocol          = "TCP"
	rule_action       = "ALLOW"
	port_range        = "80"
	ip_block          = "0.0.0.0/0"
	network_rule_type = "INBND"
	description       = "tf-testacc-network-acl-rule"
}

resource "ncloud_network_acl_rule" "nacl_rule_inbound_120" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no
	priority          = 120
	protocol          = "TCP"
	rule_action       = "ALLOW"
	port_range        = "443"
	ip_block          = "0.0.0.0/0"
	network_rule_type = "INBND"
	description       = "tf-testacc-network-acl-rule"
}

resource "ncloud_network_acl_rule" "nacl_rule_outbound_100" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no
	priority          = 100
	protocol          = "TCP"
	rule_action       = "ALLOW"
	port_range        = "1-65535"
	ip_block          = "0.0.0.0/0"
	network_rule_type = "OTBND"
	description       = "tf-testacc-network-acl-rule"
}`
}

func testAccCheckNetworkACLRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL Rule id is set: %s", n)
		}

		client := testAccProvider.Meta().(*NcloudAPIClient)

		priority, err := strconv.ParseInt(rs.Primary.Attributes["priority"], 10, 32)
		if err != nil {
			return err
		}

		networkRuleType := ncloud.String(rs.Primary.Attributes["network_rule_type"])

		reqParams := &vpc.GetNetworkAclRuleListRequest{
			NetworkAclNo:           ncloud.String(rs.Primary.Attributes["network_acl_no"]),
			NetworkAclRuleTypeCode: networkRuleType,
		}

		resp, err := client.vpc.V2Api.GetNetworkAclRuleList(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_network_acl_rule_test > GetNetworkAclRuleList", err, reqParams)
			return err
		}

		for _, i := range resp.NetworkAclRuleList {
			if *i.Priority == int32(priority) && *i.NetworkAclRuleType.Code == *networkRuleType {
				return nil
			}
		}

		return fmt.Errorf("Entry not found: %v", resp.NetworkAclRuleList)
	}
}

func testAccNcloudNetworkACLRuleImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		networkACLNo := rs.Primary.Attributes["network_acl_no"]
		networkRuleType := rs.Primary.Attributes["network_rule_type"]
		priority := rs.Primary.Attributes["priority"]

		return fmt.Sprintf("%s:%s:%s", networkACLNo, networkRuleType, priority), nil
	}
}
