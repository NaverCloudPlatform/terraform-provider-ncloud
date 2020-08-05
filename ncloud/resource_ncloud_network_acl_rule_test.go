package ncloud

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudNetworkACLRule_basic(t *testing.T) {
	var networkACLRule vpc.NetworkAclRule

	name := fmt.Sprintf("test-network-acl-rule-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLRuleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_inbound_100", &networkACLRule),
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_inbound_110", &networkACLRule),
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_inbound_120", &networkACLRule),
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule_outbound_100", &networkACLRule),
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

func TestAccResourceNcloudNetworkACLRule_disappears(t *testing.T) {
	var networkACLRule vpc.NetworkAclRule

	name := fmt.Sprintf("test-network-acl-rule-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLRuleConfigDisappear(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.test", &networkACLRule),
					testAccCheckNetworkACLRuleDisappears(&networkACLRule),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNcloudNetworkACLRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "%[1]s"
	description = "test acc for network acl"
}

resource "ncloud_network_acl_rule" "nacl_rule_inbound_100" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no
	priority          = 100
	protocol          = "TCP"
	rule_action       = "ALLOW"
	port_range        = "22"
	ip_block          = "0.0.0.0/0"
	network_rule_type = "INBND"
	description       = "%[1]s"
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
}`, name)
}

func testAccResourceNcloudNetworkACLRuleConfigDisappear(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "%[1]s"
	description = "test acc for network acl"
}

resource "ncloud_network_acl_rule" "test" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no
	priority          = 100
	protocol          = "TCP"
	rule_action       = "ALLOW"
	port_range        = "22"
	ip_block          = "0.0.0.0/0"
	network_rule_type = "INBND"
	description       = "%[1]s"
}`, name)
}

func testAccCheckNetworkACLRuleExists(n string, networkACLRule *vpc.NetworkAclRule) resource.TestCheckFunc {
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
				*networkACLRule = *i
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

func testAccCheckNetworkACLRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*NcloudAPIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_network_acl_rule" {
			continue
		}

		instance, err := getNetworkACLInstance(client, rs.Primary.Attributes["network_acl_no"])

		if err != nil {
			return err
		}

		if instance == nil {
			return nil
		}

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
		}

		for _, i := range resp.NetworkAclRuleList {
			if *i.Priority == int32(priority) && *i.NetworkAclRuleType.Code == *networkRuleType {
				return errors.New("Network ACL Rule still exists")
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckNetworkACLRuleDisappears(instance *vpc.NetworkAclRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*NcloudAPIClient)

		networkACLRule := &vpc.RemoveNetworkAclRuleParameter{
			IpBlock:          instance.IpBlock,
			RuleActionCode:   instance.RuleAction.Code,
			PortRange:        instance.PortRange,
			Priority:         instance.Priority,
			ProtocolTypeCode: instance.ProtocolType.Code,
		}

		if *instance.NetworkAclRuleType.Code == "INBND" {
			reqParams := &vpc.RemoveNetworkAclInboundRuleRequest{
				NetworkAclNo:       instance.NetworkAclNo,
				NetworkAclRuleList: []*vpc.RemoveNetworkAclRuleParameter{networkACLRule},
			}

			_, err := client.vpc.V2Api.RemoveNetworkAclInboundRule(reqParams)
			if err != nil {
				return err
			}
		} else {
			reqParams := &vpc.RemoveNetworkAclOutboundRuleRequest{
				NetworkAclNo:       instance.NetworkAclNo,
				NetworkAclRuleList: []*vpc.RemoveNetworkAclRuleParameter{networkACLRule},
			}

			_, err := client.vpc.V2Api.RemoveNetworkAclOutboundRule(reqParams)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
