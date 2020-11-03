package ncloud

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudNetworkACLRule_basic(t *testing.T) {
	var networkACLRule []*vpc.NetworkAclRule

	resourceName := "ncloud_network_acl_rule.nacl_rule"
	name := fmt.Sprintf("test-network-acl-rule-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLRuleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceName, &networkACLRule),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "inbound.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "outbound.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudNetworkACLRule_AssociatedSubnet(t *testing.T) {
	var networkACLRule []*vpc.NetworkAclRule

	name := fmt.Sprintf("test-nacl-rule-subnet-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLRuleConfigAssociatedSubnet(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists("ncloud_network_acl_rule.nacl_rule", &networkACLRule),
				),
			},
		},
	})
}

func TestAccResourceNcloudNetworkACLRule_disappears(t *testing.T) {
	var networkACLRule []*vpc.NetworkAclRule

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

resource "ncloud_network_acl_rule" "nacl_rule" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no

	inbound {
		priority    = 1
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
	
	inbound {
		priority    = 2
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "443"
		ip_block    = "0.0.0.0/0"
	}

	outbound {
		priority    = 3
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
}
`, name)
}

func testAccResourceNcloudNetworkACLRuleConfigAssociatedSubnet(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.3.0.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_network_acl.nacl.id
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "%[1]s"
	description = "test acc for network acl"
}

resource "ncloud_network_acl_rule" "nacl_rule" {
	network_acl_no    = ncloud_network_acl.nacl.network_acl_no
	inbound {
		priority    = 1
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
	
	inbound {
		priority    = 2
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "443"
		ip_block    = "0.0.0.0/0"
	}

	outbound {
		priority    = 3
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
	}
	depends_on        = [ncloud_subnet.subnet]
}
`, name)
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
	inbound {
		priority    = 1
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
		description       = "%[1]s"
	}
	
	inbound {
		priority    = 2
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "443"
		ip_block    = "0.0.0.0/0"
		description       = "%[1]s"
	}

	outbound {
		priority    = 3
		protocol    = "TCP"
		rule_action = "ALLOW"
		port_range  = "80"
		ip_block    = "0.0.0.0/0"
		description       = "%[1]s"
	}	
}`, name)
}

func testAccCheckNetworkACLRuleExists(n string, networkACLRule *[]*vpc.NetworkAclRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL Rule id is set: %s", n)
		}

		config := testAccProvider.Meta().(*ProviderConfig)

		rules, err := getNetworkACLRuleList(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if len(rules) == 0 {
			return fmt.Errorf("Entry not found: %s", rs.Primary.ID)
		}

		*networkACLRule = rules

		return nil
	}
}

func testAccCheckNetworkACLRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_network_acl_rule" {
			continue
		}

		rules, err := getNetworkACLRuleList(config, rs.Primary.Attributes["network_acl_no"])

		if err != nil {
			return err
		}

		if len(rules) > 0 {
			return errors.New("Network ACL Rule still exists")
		}
	}

	return nil
}

func testAccCheckNetworkACLRuleDisappears(instance *[]*vpc.NetworkAclRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		var inbound []*vpc.RemoveNetworkAclRuleParameter
		var outbound []*vpc.RemoveNetworkAclRuleParameter

		if len(*instance) == 0 {
			return nil
		}

		for _, r := range *instance {
			networkACLRule := &vpc.RemoveNetworkAclRuleParameter{
				IpBlock:          r.IpBlock,
				RuleActionCode:   r.RuleAction.Code,
				PortRange:        r.PortRange,
				Priority:         r.Priority,
				ProtocolTypeCode: r.ProtocolType.Code,
			}

			if *r.NetworkAclRuleType.Code == "INBND" {
				inbound = append(inbound, networkACLRule)
			} else {
				outbound = append(outbound, networkACLRule)
			}
		}

		if len(inbound) > 0 {
			reqParams := &vpc.RemoveNetworkAclInboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       (*instance)[0].NetworkAclNo,
				NetworkAclRuleList: inbound,
			}

			_, err := config.Client.vpc.V2Api.RemoveNetworkAclInboundRule(reqParams)
			if err != nil {
				return err
			}
		} else if len(outbound) > 0 {
			reqParams := &vpc.RemoveNetworkAclOutboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       (*instance)[0].NetworkAclNo,
				NetworkAclRuleList: outbound,
			}

			_, err := config.Client.vpc.V2Api.RemoveNetworkAclOutboundRule(reqParams)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
