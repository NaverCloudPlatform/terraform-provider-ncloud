package ncloud

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"strconv"
	"strings"
	"time"

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
					return nil, fmt.Errorf("unexpected format of ID (%q), expected NETWORK_ACL_NO:RULE_TYPE:PRIORITY", d.Id())
				}
				networkACLNo := idParts[0]
				networkRuleType := idParts[1]
				priority, err := strconv.Atoi(idParts[2])
				if err != nil {
					return nil, err
				}

				d.Set("network_acl_no", networkACLNo)
				d.Set("priority", priority)
				d.Set("rule_type", networkRuleType)
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
			"rule_type": {
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
	config := meta.(*ProviderConfig)

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

	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		if d.Get("rule_type").(string) == "INBND" {
			reqParams = &vpc.AddNetworkAclInboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
				NetworkAclRuleList: []*vpc.AddNetworkAclRuleParameter{networkACLRule},
			}

			logCommonRequest("resource_ncloud_network_acl_rule > AddNetworkAclInboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.AddNetworkAclInboundRule(reqParams.(*vpc.AddNetworkAclInboundRuleRequest))
		} else {
			reqParams = &vpc.AddNetworkAclOutboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
				NetworkAclRuleList: []*vpc.AddNetworkAclRuleParameter{networkACLRule},
			}

			logCommonRequest("resource_ncloud_network_acl_rule > AddNetworkAclOutboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.AddNetworkAclOutboundRule(reqParams.(*vpc.AddNetworkAclOutboundRuleRequest))
		}

		if err == nil {
			return resource.NonRetryableError(err)
		}

		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == ApiErrorNetworkAclCantAccessaApropriate {
			logErrorResponse("retry AddNetworkAclRule", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("resource_ncloud_network_acl_rule > AddNetworkAclRule", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_network_acl_rule > AddNetworkAclRule", resp)

	d.SetId(networkACLIdRuleHash(d.Get("network_acl_no").(string), d.Get("rule_type").(string), d.Get("priority").(int)))

	log.Printf("[INFO] Network ACL Rule ID: %s", d.Id())

	if err = waitForNcloudNetworkACLRunning(config, d.Get("network_acl_no").(string)); err != nil {
		return err
	}

	return resourceNcloudNetworkACLRuleRead(d, meta)
}

func resourceNcloudNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getNetworkACLRuleInstance(d, config)
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
	d.Set("rule_type", instance.NetworkAclRuleType.Code)
	d.Set("description", instance.NetworkAclRuleDescription)

	return nil
}

func resourceNcloudNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudNetworkACLRuleRead(d, meta)
}

func resourceNcloudNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	networkACLRule := &vpc.RemoveNetworkAclRuleParameter{
		IpBlock:          ncloud.String(d.Get("ip_block").(string)),
		RuleActionCode:   ncloud.String(d.Get("rule_action").(string)),
		PortRange:        ncloud.String(d.Get("port_range").(string)),
		Priority:         ncloud.Int32(int32(d.Get("priority").(int))),
		ProtocolTypeCode: ncloud.String(d.Get("protocol").(string)),
	}

	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error

		if d.Get("rule_type").(string) == "INBND" {
			reqParams := &vpc.RemoveNetworkAclInboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
				NetworkAclRuleList: []*vpc.RemoveNetworkAclRuleParameter{networkACLRule},
			}

			logCommonRequest("resource_ncloud_network_acl_rule > RemoveNetworkAclInboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.RemoveNetworkAclInboundRule(reqParams)
		} else {
			reqParams := &vpc.RemoveNetworkAclOutboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Get("network_acl_no").(string)),
				NetworkAclRuleList: []*vpc.RemoveNetworkAclRuleParameter{networkACLRule},
			}

			logCommonRequest("resource_ncloud_network_acl_rule > RemoveNetworkAclOutboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.RemoveNetworkAclOutboundRule(reqParams)
		}

		if err == nil {
			return resource.NonRetryableError(err)
		}

		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == ApiErrorNetworkAclCantAccessaApropriate {
			logErrorResponse("retry RemoveNetworkAclRule", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("resource_ncloud_network_acl_rule > RemoveNetworkAclRule", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_network_acl_rule > RemoveNetworkAclRule", resp)

	if err = waitForNcloudNetworkACLRunning(config, d.Get("network_acl_no").(string)); err != nil {
		return err
	}

	return nil
}

func getNetworkACLRuleInstance(d *schema.ResourceData, config *ProviderConfig) (*vpc.NetworkAclRule, error) {
	reqParams := &vpc.GetNetworkAclRuleListRequest{
		RegionCode:             &config.RegionCode,
		NetworkAclNo:           ncloud.String(d.Get("network_acl_no").(string)),
		NetworkAclRuleTypeCode: ncloud.String(d.Get("rule_type").(string)),
	}

	logCommonRequest("resource_ncloud_network_acl_rule > GetNetworkAclRuleList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNetworkAclRuleList(reqParams)
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

func waitForNcloudNetworkACLRunning(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkACLInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Network ACL (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func networkACLIdRuleHash(networkACLId string, ruleType string, priority int) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", networkACLId))
	buf.WriteString(fmt.Sprintf("%s-", ruleType))
	buf.WriteString(fmt.Sprintf("%d-", priority))
	return fmt.Sprintf("nacl-%d", hashcode.String(buf.String()))
}
