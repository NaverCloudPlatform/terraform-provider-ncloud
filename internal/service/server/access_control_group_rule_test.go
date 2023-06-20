package server

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func TestAccResourceNcloudAccessControlGroupRule_basic(t *testing.T) {
	var AccessControlGroupRule []*vserver.AccessControlGroupRule
	name := fmt.Sprintf("tf-acg-rule-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_access_control_group_rule.acg_rule_foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(true),
		CheckDestroy: testAccCheckAccessControlGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudAccessControlGroupRuleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessControlGroupRuleExists(resourceName, &AccessControlGroupRule),
					resource.TestMatchResourceAttr(resourceName, "access_control_group_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "inbound.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "outbound.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceNcloudAccessControlGroupRule_disappears(t *testing.T) {
	var AccessControlGroupRule []*vserver.AccessControlGroupRule
	name := fmt.Sprintf("tf-nic-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_access_control_group_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(true),
		CheckDestroy: testAccCheckAccessControlGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudAccessControlGroupRuleConfigDisappear(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessControlGroupRuleExists(resourceName, &AccessControlGroupRule),
					testAccCheckAccessControlGroupRuleDisappears(&AccessControlGroupRule),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNcloudAccessControlGroupRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_access_control_group" "foo" {
	name                  = "%[1]s"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id
}

resource "ncloud_access_control_group" "bar" {
	name                  = "%[1]s-src"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id
}

resource "ncloud_access_control_group_rule" "acg_rule_foo" {
	access_control_group_no = ncloud_access_control_group.foo.id

	inbound {
		protocol    = "TCP"
		port_range  = "8083"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}
	
	inbound {
		protocol    = "TCP"
		port_range  = "9000-10000"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}

	inbound {
		protocol    = "254"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}

	inbound {
		protocol    = "120"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}

	inbound {
		protocol    = "ICMP"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}

	inbound {
		protocol                       = "TCP"
		source_access_control_group_no = ncloud_access_control_group.bar.id
		port_range                     = "22"
		description                    = "%[1]s"
	}

	outbound {
		protocol    = "TCP"
		port_range  = "8083"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}

	outbound {
		protocol    = "120"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}
}
`, name)
}

func testAccResourceNcloudAccessControlGroupRuleConfigDisappear(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_access_control_group" "test" {
	name                  = "%[1]s"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id
}

resource "ncloud_access_control_group_rule" "test" {
	access_control_group_no = ncloud_access_control_group.test.id
	inbound {
		protocol    = "TCP"
		port_range  = "8082"
		ip_block    = "0.0.0.0/0"
		description = "%[1]s"
	}
}
`, name)
}

func testAccCheckAccessControlGroupRuleExists(n string, AccessControlGroupRule *[]*vserver.AccessControlGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Access Control Group id is set")
		}

		config := GetTestProvider(true).Meta().(*ProviderConfig)

		rules, err := getAccessControlGroupRuleList(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(rules) == 0 {
			return fmt.Errorf("Entry not found: %s", rs.Primary.ID)
		}

		*AccessControlGroupRule = rules

		return nil
	}
}

func testAccCheckAccessControlGroupRuleDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_access_control_group_rule" {
			continue
		}

		instance, err := getAccessControlGroup(config, rs.Primary.Attributes["access_control_group_no"])

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Access Control Group still exists")
		}
	}

	return nil
}

func testAccCheckAccessControlGroupRuleDisappears(instance *[]*vserver.AccessControlGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GetTestProvider(true).Meta().(*ProviderConfig)

		if len(*instance) == 0 {
			return nil
		}

		id := (*instance)[0].AccessControlGroupNo

		accessControlGroup, err := getAccessControlGroup(config, *id)
		if err != nil {
			return err
		}

		if accessControlGroup == nil {
			return fmt.Errorf("no matching Access Control Group: %s", *id)
		}

		var inbound []*vserver.RemoveAccessControlGroupRuleParameter
		var outbound []*vserver.RemoveAccessControlGroupRuleParameter

		for _, r := range *instance {
			rule := &vserver.RemoveAccessControlGroupRuleParameter{
				IpBlock:                    r.IpBlock,
				AccessControlGroupSequence: r.AccessControlGroupSequence,
				PortRange:                  r.PortRange,
				ProtocolTypeCode:           r.ProtocolType.Code,
			}

			if *r.AccessControlGroupRuleType.Code == "INBND" {
				inbound = append(inbound, rule)
			} else {
				outbound = append(outbound, rule)
			}
		}

		if len(inbound) > 0 {
			reqParams := &vserver.RemoveAccessControlGroupInboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       id,
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: inbound,
			}

			_, err := config.Client.Vserver.V2Api.RemoveAccessControlGroupInboundRule(reqParams)
			if err != nil {
				return err
			}
		} else if len(outbound) > 0 {
			reqParams := &vserver.RemoveAccessControlGroupOutboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       id,
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: outbound,
			}

			_, err := config.Client.Vserver.V2Api.RemoveAccessControlGroupOutboundRule(reqParams)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
