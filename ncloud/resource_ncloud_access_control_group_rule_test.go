package ncloud

import (
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudAccessControlGroupRule_basic(t *testing.T) {
	var AccessControlGroupRule vserver.AccessControlGroupRule
	name := fmt.Sprintf("tf-acg-rule-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_access_control_group_rule.test-inbound-tcp-8082"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessControlGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudAccessControlGroupRuleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessControlGroupRuleExists(resourceName, &AccessControlGroupRule),
					resource.TestMatchResourceAttr(resourceName, "access_control_group_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "INBND"),
					resource.TestCheckResourceAttr(resourceName, "ip_block", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "port_range", "8082"),

					testAccCheckAccessControlGroupRuleExists("ncloud_access_control_group_rule.test-inbound-tcp-8083", &AccessControlGroupRule),
					testAccCheckAccessControlGroupRuleExists("ncloud_access_control_group_rule.test-inbound-tcp-9000-10000", &AccessControlGroupRule),
					testAccCheckAccessControlGroupRuleExists("ncloud_access_control_group_rule.test-inbound-icmp", &AccessControlGroupRule),
					testAccCheckAccessControlGroupRuleExists("ncloud_access_control_group_rule.test-inbound-acg-22", &AccessControlGroupRule),
					testAccCheckAccessControlGroupRuleExists("ncloud_access_control_group_rule.test-outbound-tcp-8083", &AccessControlGroupRule),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccNcloudAccessControlGroupImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNcloudAccessControlGroupImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		accessControlGroupNo := rs.Primary.Attributes["access_control_group_no"]
		ruleType := rs.Primary.Attributes["rule_type"]
		protocol := rs.Primary.Attributes["protocol"]
		accessSource := rs.Primary.Attributes["ip_block"]
		if len(rs.Primary.Attributes["source_access_control_group_no"]) > 0 {
			accessSource = rs.Primary.Attributes["source_access_control_group_no"]
		}
		portRange := rs.Primary.Attributes["port_range"]

		id := fmt.Sprintf("%s:%s:%s:%s:%s", accessControlGroupNo, ruleType, protocol, accessSource, portRange)

		log.Printf("[INFO] testAccNcloudAccessControlGroupImportStateIDFunc: %s", id)

		return id, nil
	}
}

func TestAccResourceNcloudAccessControlGroupRule_disappears(t *testing.T) {
	var AccessControlGroupRule vserver.AccessControlGroupRule
	name := fmt.Sprintf("tf-nic-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_access_control_group_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

resource "ncloud_access_control_group" "test" {
	name                  = "%[1]s"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id
}

resource "ncloud_access_control_group" "foo" {
	name                  = "%[1]s-src"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id
}

resource "ncloud_access_control_group_rule" "test-inbound-tcp-8082" {
	access_control_group_no = ncloud_access_control_group.test.id
	description             = "%[1]s"
	rule_type               = "INBND"
	protocol                = "TCP"
	ip_block                = "0.0.0.0/0"
	port_range              = "8082"
}

resource "ncloud_access_control_group_rule" "test-inbound-tcp-8083" {
	access_control_group_no = ncloud_access_control_group.test.id
	description             = "%[1]s"
	rule_type               = "INBND"
	protocol                = "TCP"
	ip_block                = "0.0.0.0/0"
	port_range              = "8083"
}

resource "ncloud_access_control_group_rule" "test-inbound-tcp-9000-10000" {
	access_control_group_no = ncloud_access_control_group.test.id
	description             = "%[1]s"
	rule_type               = "INBND"
	protocol                = "TCP"
	ip_block                = "0.0.0.0/0"
	port_range              = "9000-10000"
}

resource "ncloud_access_control_group_rule" "test-inbound-icmp" {
	access_control_group_no = ncloud_access_control_group.test.id
	description             = "%[1]s"
	rule_type               = "INBND"
	protocol                = "ICMP"
	ip_block                = "0.0.0.0/0"
}

resource "ncloud_access_control_group_rule" "test-inbound-acg-22" {
	access_control_group_no = ncloud_access_control_group.test.id
	description                 = "%[1]s"
	rule_type                   = "INBND"
	protocol                    = "TCP"
	source_access_control_group_no = ncloud_access_control_group.foo.id
	port_range                  = "22"
}

resource "ncloud_access_control_group_rule" "test-outbound-tcp-8083" {
	access_control_group_no = ncloud_access_control_group.test.id
	description             = "%[1]s"
	rule_type               = "OTBND"
	protocol                = "TCP"
	ip_block                = "0.0.0.0/0"
	port_range              = "8083"
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
	description             = "%[1]s"
	rule_type               = "INBND"
	protocol                = "TCP"
	ip_block                = "0.0.0.0/0"
	port_range              = "8082"
}
`, name)
}

func testAccCheckAccessControlGroupRuleExists(n string, AccessControlGroupRule *vserver.AccessControlGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Access Control Group id is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)

		rule := &AccessControlGroupRuleParam{
			AccessControlGroupNo:     rs.Primary.Attributes["access_control_group_no"],
			RuleType:                 rs.Primary.Attributes["rule_type"],
			Protocol:                 rs.Primary.Attributes["protocol"],
			IpBlock:                  rs.Primary.Attributes["ip_block"],
			SourceAccessControlGroup: rs.Primary.Attributes["source_access_control_group_no"],
			PortRange:                rs.Primary.Attributes["port_range"],
		}

		instance, err := getAccessControlGroupRule(config, rule)
		if err != nil {
			return err
		}

		*AccessControlGroupRule = *instance

		return nil
	}
}

func testAccCheckAccessControlGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_access_control_group_rule" {
			continue
		}

		rule := &AccessControlGroupRuleParam{
			AccessControlGroupNo:     rs.Primary.Attributes["access_control_group_no"],
			RuleType:                 rs.Primary.Attributes["rule_type"],
			Protocol:                 rs.Primary.Attributes["protocol"],
			IpBlock:                  rs.Primary.Attributes["ip_block"],
			SourceAccessControlGroup: rs.Primary.Attributes["source_access_control_group_no"],
			PortRange:                rs.Primary.Attributes["port_range"],
		}

		instance, err := getAccessControlGroupRule(config, rule)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Access Control Group still exists")
		}
	}

	return nil
}

func testAccCheckAccessControlGroupRuleDisappears(instance *vserver.AccessControlGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rule := &AccessControlGroupRuleParam{
			AccessControlGroupNo:     *instance.AccessControlGroupNo,
			RuleType:                 *instance.AccessControlGroupRuleType.Code,
			Protocol:                 *instance.ProtocolType.Code,
			IpBlock:                  *instance.IpBlock,
			SourceAccessControlGroup: *instance.AccessControlGroupSequence,
			PortRange:                *instance.PortRange,
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		return deleteAccessControlGroupRule(&schema.ResourceData{}, config, rule)
	}
}
