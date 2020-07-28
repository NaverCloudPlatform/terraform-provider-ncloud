package ncloud

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNcloudNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNetworkACLRuleCreate,
		Read:   resourceNcloudNetworkACLRuleRead,
		Update: resourceNcloudNetworkACLRuleUpdate,
		Delete: resourceNcloudNetworkACLRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ":")
				if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected NETWORK_ACL_NO:NETWORK_RULE_TYPE:PRIORITY", d.Id())
				}
				networkACLNo := idParts[0]
				networkRuleType := idParts[1]
				priority, err := strconv.Atoi(idParts[2])
				if err != nil {
					return nil, err
				}

				d.Set("network_acl_no", networkACLNo)
				d.Set("priority", priority)
				d.Set("network_rule_type", networkRuleType)
				d.SetId(networkACLIdRuleHash(networkACLNo, networkRuleType, priority))
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"network_acl_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of Network ACL. Get available values using the `default_network_acl_no` from Resource `ncloud_vpc` or Data source `data.ncloud_network_acls`.",
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "The priority can be entered in the range of 0 to 199 digits.",
				ValidateFunc: validation.IntBetween(0, 199),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
				Description:  "The protocol. TCP(TCP), UDP(UDP), ICMP(ICMP)",
			},
			"port_range": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePortRange,
				Description:  "The range port from to. Not entered if protocol type code is ICMP. (e.g. \"22\" or \"1-65535\")",
			},
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DROP"}, false),
				Description:  "Indicates whether to allow or block the traffic that matches the rule. ALLOW(Allowed), DROP(Blocked).",
			},
			"ip_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 32),
				Description:  "The network range to allow or block, in CIDR notation (e.g. \"100.10.20.0/24\")",
			},
			"network_rule_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"INBND", "OTBND"}, false),
				Description:  "Indicates whether this is an inbound rule. INBND(Inbound), OTBND(Outbound)",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
				Description:  "Description of rule",
			},
		},
	}
}

func resourceNcloudNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	networkACLRule := &vpc.AddNetworkAclRuleParameter{
		IpBlock:          ncloud.String(d.Get("ip_block").(string)),
		RuleActionCode:   ncloud.String(d.Get("rule_action").(string)),
		Priority:         ncloud.Int32(int32(d.Get("priority").(int))),
		ProtocolTypeCode: ncloud.String(d.Get("protocol").(string)),
	}

	if v, ok := d.GetOk("port_range"); ok {
		networkACLRule.PortRange = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		networkACLRule.NetworkAclRuleDescription = ncloud.String(v.(string))
	}

	if d.Get("network_rule_type").(string) == "INBND" {
		reqParams := &vpc.AddNetworkAclInboundRuleRequest{
			RegionCode:         regionCode,
			NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
			NetworkAclRuleList: []*vpc.AddNetworkAclRuleParameter{networkACLRule},
		}

		logCommonRequest("resource_ncloud_network_acl_rule > AddNetworkAclInboundRule", reqParams)
		resp, err := client.vpc.V2Api.AddNetworkAclInboundRule(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_network_acl_rule > AddNetworkAclInboundRule", err, reqParams)
			return err
		}

		logResponse("resource_ncloud_network_acl_rule > AddNetworkAclInboundRule", resp)
	} else {
		reqParams := &vpc.AddNetworkAclOutboundRuleRequest{
			RegionCode:         regionCode,
			NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
			NetworkAclRuleList: []*vpc.AddNetworkAclRuleParameter{networkACLRule},
		}

		logCommonRequest("resource_ncloud_network_acl_rule > AddNetworkAclOutboundRule", reqParams)
		resp, err := client.vpc.V2Api.AddNetworkAclOutboundRule(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_network_acl_rule > AddNetworkAclOutboundRule", err, reqParams)
			return err
		}

		logResponse("resource_ncloud_network_acl_rule > AddNetworkAclOutboundRule", resp)
	}

	d.SetId(networkACLIdRuleHash(d.Get("network_acl_no").(string), d.Get("network_rule_type").(string), d.Get("priority").(int)))

	log.Printf("[INFO] Network ACL Rule ID: %s", d.Id())

	return resourceNcloudNetworkACLRuleRead(d, meta)
}

func resourceNcloudNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	instance, err := getNetworkACLRuleInstance(client, d)
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(networkACLIdRuleHash(*instance.NetworkAclNo, *instance.NetworkAclRuleType.Code, int(*instance.Priority)))
	d.Set("network_acl_no", instance.NetworkAclNo)
	d.Set("priority", instance.Priority)
	d.Set("protocol", instance.ProtocolType.Code)
	d.Set("port_range", instance.PortRange)
	d.Set("rule_action", instance.RuleAction.Code)
	d.Set("ip_block", instance.IpBlock)
	d.Set("network_rule_type", instance.NetworkAclRuleType.Code)
	d.Set("description", instance.NetworkAclRuleDescription)

	return nil
}

func resourceNcloudNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudNetworkACLRuleRead(d, meta)
}

func resourceNcloudNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	networkACLRule := &vpc.RemoveNetworkAclRuleParameter{
		IpBlock:          ncloud.String(d.Get("ip_block").(string)),
		RuleActionCode:   ncloud.String(d.Get("rule_action").(string)),
		PortRange:        ncloud.String(d.Get("port_range").(string)),
		Priority:         ncloud.Int32(int32(d.Get("priority").(int))),
		ProtocolTypeCode: ncloud.String(d.Get("protocol").(string)),
	}

	if d.Get("network_rule_type").(string) == "INBND" {
		reqParams := &vpc.RemoveNetworkAclInboundRuleRequest{
			RegionCode:         regionCode,
			NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
			NetworkAclRuleList: []*vpc.RemoveNetworkAclRuleParameter{networkACLRule},
		}

		logCommonRequest("resource_ncloud_network_acl_rule > RemoveNetworkAclInboundRule", reqParams)
		resp, err := client.vpc.V2Api.RemoveNetworkAclInboundRule(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_network_acl_rule > RemoveNetworkAclInboundRule", err, reqParams)
			return err
		}

		logResponse("resource_ncloud_network_acl_rule > RemoveNetworkAclInboundRule", resp)
	} else {
		reqParams := &vpc.RemoveNetworkAclOutboundRuleRequest{
			RegionCode:         regionCode,
			NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
			NetworkAclRuleList: []*vpc.RemoveNetworkAclRuleParameter{networkACLRule},
		}

		logCommonRequest("resource_ncloud_network_acl_rule > RemoveNetworkAclOutboundRule", reqParams)
		resp, err := client.vpc.V2Api.RemoveNetworkAclOutboundRule(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_network_acl_rule > RemoveNetworkAclOutboundRule", err, reqParams)
			return err
		}

		logResponse("resource_ncloud_network_acl_rule > RemoveNetworkAclOutboundRule", resp)
	}

	return nil
}

func getNetworkACLRuleInstance(client *NcloudAPIClient, d *schema.ResourceData) (*vpc.NetworkAclRule, error) {
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return nil, err
	}

	reqParams := &vpc.GetNetworkAclRuleListRequest{
		RegionCode:             regionCode,
		NetworkAclNo:           ncloud.String(d.Get("network_acl_no").(string)),
		NetworkAclRuleTypeCode: ncloud.String(d.Get("network_rule_type").(string)),
	}

	logCommonRequest("resource_ncloud_network_acl_rule > GetNetworkAclRuleList", reqParams)
	resp, err := client.vpc.V2Api.GetNetworkAclRuleList(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_network_acl_rule > GetNetworkAclRuleList", err, reqParams)
		return nil, err
	}
	logResponse("resource_ncloud_network_acl_rule > GetNetworkAclRuleList", resp)

	if resp.NetworkAclRuleList != nil {
		for _, i := range resp.NetworkAclRuleList {
			if *i.Priority == int32(d.Get("priority").(int)) {
				return i, nil
			}
		}
		return nil, nil
	}

	return nil, fmt.Errorf("Not found acl rule, got: %#v", resp.NetworkAclRuleList)
}

func networkACLIdRuleHash(networkACLId string, ruleType string, priority int) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", networkACLId))
	buf.WriteString(fmt.Sprintf("%s-", ruleType))
	buf.WriteString(fmt.Sprintf("%d-", priority))
	return fmt.Sprintf("nacl-%d", hashcode.String(buf.String()))
}
