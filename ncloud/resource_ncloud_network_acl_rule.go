package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_network_acl_rule", resourceNcloudNetworkACLRule())
}

func resourceNcloudNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNetworkACLRuleCreate,
		Read:   resourceNcloudNetworkACLRuleRead,
		Update: resourceNcloudNetworkACLRuleUpdate,
		Delete: resourceNcloudNetworkACLRuleDelete,
		Schema: map[string]*schema.Schema{
			"network_acl_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"inbound": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 199)),
						},
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false)),
						},
						"ip_block": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IsCIDRNetwork(0, 32)),
						},
						"deny_allow_group_no": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_action": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DROP"}, false)),
						},
						"port_range": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validatePortRange),
							Default:          "",
						},
						"description": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 1000)),
							Default:          "",
						},
					},
				},
			},
			"outbound": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(0, 199)),
						},
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false)),
						},
						"ip_block": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IsCIDRNetwork(0, 32)),
						},
						"deny_allow_group_no": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_action": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DROP"}, false)),
						},
						"port_range": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validatePortRange),
							Default:          "",
						},
						"description": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 1000)),
							Default:          "",
						},
					},
				},
			},
		},
	}
}

func resourceNcloudNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_network_acl_rule`")
	}

	d.SetId(d.Get("network_acl_no").(string))
	log.Printf("[INFO] Network ACL ID: %s", d.Id())

	return resourceNcloudNetworkACLRuleUpdate(d, meta)
}

func resourceNcloudNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	rules, err := getNetworkACLRuleList(config, d.Id())
	if err != nil {
		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == "1011002" { // You cannot access the appropriate Network ACL
			d.SetId("")
		}
		return err
	}

	if len(rules) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("network_acl_no", d.Id())

	// Create empty set for getNetworkACLRuleList
	iSet := schema.NewSet(schema.HashResource(resourceNcloudNetworkACLRule().Schema["inbound"].Elem.(*schema.Resource)), []interface{}{})
	oSet := schema.NewSet(schema.HashResource(resourceNcloudNetworkACLRule().Schema["outbound"].Elem.(*schema.Resource)), []interface{}{})

	for _, r := range rules {
		m := map[string]interface{}{
			"priority":            int(*r.Priority),
			"protocol":            *r.ProtocolType.Code,
			"port_range":          *r.PortRange,
			"rule_action":         *r.RuleAction.Code,
			"ip_block":            *r.IpBlock,
			"deny_allow_group_no": *r.DenyAllowGroupNo,
			"description":         *r.NetworkAclRuleDescription,
		}

		if *r.NetworkAclRuleType.Code == "INBND" {
			iSet.Add(m)
		} else {
			oSet.Add(m)
		}
	}

	// Only set data intersection between resource and list
	if err := d.Set("inbound", iSet.List()); err != nil {
		log.Printf("[WARN] Error setting inbound rule set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("outbound", oSet.List()); err != nil {
		log.Printf("[WARN] Error setting outbound rule set for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNcloudNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("inbound") {
		if err := updateNetworkACLRule(d, config, "inbound"); err != nil {
			return err
		}
	}

	if d.HasChange("outbound") {
		if err := updateNetworkACLRule(d, config, "outbound"); err != nil {
			return err
		}
	}

	return resourceNcloudNetworkACLRuleRead(d, meta)
}

func resourceNcloudNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	i := d.Get("inbound").(*schema.Set)
	o := d.Get("outbound").(*schema.Set)

	if len(i.List()) > 0 {
		if err := removeNetworkACLRule(d, config, "inbound", expandRemoveNetworkAclRule(i.List())); err != nil {
			return err
		}
	}

	if len(o.List()) > 0 {
		if err := removeNetworkACLRule(d, config, "outbound", expandRemoveNetworkAclRule(o.List())); err != nil {
			return err
		}
	}

	return nil
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

func getNetworkACLRuleList(config *ProviderConfig, id string) ([]*vpc.NetworkAclRule, error) {
	reqParams := &vpc.GetNetworkAclRuleListRequest{
		RegionCode:   &config.RegionCode,
		NetworkAclNo: ncloud.String(id),
	}

	logCommonRequest("GetNetworkAclRuleList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNetworkAclRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetNetworkAclRuleList", err, reqParams)
		return nil, err
	}
	logResponse("GetNetworkAclRuleList", resp)

	return resp.NetworkAclRuleList, nil
}

func updateNetworkACLRule(d *schema.ResourceData, config *ProviderConfig, ruleType string) error {
	o, n := d.GetChange(ruleType)

	if o == nil {
		o = new(schema.Set)
	}
	if n == nil {
		n = new(schema.Set)
	}

	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	add := ns.Difference(os).List()
	remove := os.Difference(ns).List()

	removeNetworkACLRuleList := expandRemoveNetworkAclRule(remove)
	addNetworkACLRuleList := expandAddNetworkAclRule(add)

	if len(removeNetworkACLRuleList) > 0 {
		if err := removeNetworkACLRule(d, config, ruleType, removeNetworkACLRuleList); err != nil {
			return err
		}
	}

	if len(addNetworkACLRuleList) > 0 {
		if err := addNetworkACLRule(d, config, ruleType, addNetworkACLRuleList); err != nil {
			return err
		}
	}

	return nil
}

func addNetworkACLRule(d *schema.ResourceData, config *ProviderConfig, ruleType string, addNetworkRuleList []*vpc.AddNetworkAclRuleParameter) error {
	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		if ruleType == "inbound" {
			reqParams = &vpc.AddNetworkAclInboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Id()),
				NetworkAclRuleList: addNetworkRuleList,
			}

			logCommonRequest("AddNetworkAclInboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.AddNetworkAclInboundRule(reqParams.(*vpc.AddNetworkAclInboundRuleRequest))
		} else {
			reqParams = &vpc.AddNetworkAclOutboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Id()),
				NetworkAclRuleList: addNetworkRuleList,
			}

			logCommonRequest("AddNetworkAclOutboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.AddNetworkAclOutboundRule(reqParams.(*vpc.AddNetworkAclOutboundRuleRequest))
		}

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if containsInStringList(errBody.ReturnCode, []string{ApiErrorNetworkAclCantAccessaApropriate, ApiErrorNetworkAclRuleChangeIngRules}) {
				logErrorResponse("retry AddNetworkAclRule", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		logErrorResponse("AddNetworkAclRule", err, reqParams)
		return err
	}

	logResponse("AddNetworkAclRule", resp)

	if err = waitForNcloudNetworkACLRunning(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func removeNetworkACLRule(d *schema.ResourceData, config *ProviderConfig, ruleType string, removeNetworkRuleList []*vpc.RemoveNetworkAclRuleParameter) error {
	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error

		if ruleType == "inbound" {
			reqParams = &vpc.RemoveNetworkAclInboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Id()),
				NetworkAclRuleList: removeNetworkRuleList,
			}

			logCommonRequest("RemoveNetworkAclInboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.RemoveNetworkAclInboundRule(reqParams.(*vpc.RemoveNetworkAclInboundRuleRequest))
		} else {
			reqParams = &vpc.RemoveNetworkAclOutboundRuleRequest{
				RegionCode:         &config.RegionCode,
				NetworkAclNo:       ncloud.String(d.Id()),
				NetworkAclRuleList: removeNetworkRuleList,
			}

			logCommonRequest("RemoveNetworkAclOutboundRule", reqParams)
			resp, err = config.Client.vpc.V2Api.RemoveNetworkAclOutboundRule(reqParams.(*vpc.RemoveNetworkAclOutboundRuleRequest))
		}

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if containsInStringList(errBody.ReturnCode, []string{ApiErrorNetworkAclCantAccessaApropriate, ApiErrorNetworkAclRuleChangeIngRules}) {
				logErrorResponse("retry RemoveNetworkAclRule", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		logErrorResponse("RemoveNetworkAclRule", err, reqParams)
		return err
	}

	logResponse("RemoveNetworkAclRule", resp)

	if err = waitForNcloudNetworkACLRunning(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandAddNetworkAclRule(rules []interface{}) []*vpc.AddNetworkAclRuleParameter {
	var networkRuleList []*vpc.AddNetworkAclRuleParameter

	for _, vi := range rules {
		m := vi.(map[string]interface{})
		networkACLRule := &vpc.AddNetworkAclRuleParameter{
			IpBlock:                   ncloud.String(m["ip_block"].(string)),
			DenyAllowGroupNo:          ncloud.String(m["deny_allow_group_no"].(string)),
			RuleActionCode:            ncloud.String(m["rule_action"].(string)),
			Priority:                  ncloud.Int32(int32(m["priority"].(int))),
			ProtocolTypeCode:          ncloud.String(m["protocol"].(string)),
			PortRange:                 ncloud.String(m["port_range"].(string)),
			NetworkAclRuleDescription: ncloud.String(m["description"].(string)),
		}

		networkRuleList = append(networkRuleList, networkACLRule)
	}

	return networkRuleList
}

func expandRemoveNetworkAclRule(rules []interface{}) []*vpc.RemoveNetworkAclRuleParameter {
	var networkRuleList []*vpc.RemoveNetworkAclRuleParameter

	for _, vi := range rules {
		m := vi.(map[string]interface{})
		networkACLRule := &vpc.RemoveNetworkAclRuleParameter{
			IpBlock:          ncloud.String(m["ip_block"].(string)),
			DenyAllowGroupNo: ncloud.String(m["deny_allow_group_no"].(string)),
			RuleActionCode:   ncloud.String(m["rule_action"].(string)),
			Priority:         ncloud.Int32(int32(m["priority"].(int))),
			ProtocolTypeCode: ncloud.String(m["protocol"].(string)),
			PortRange:        ncloud.String(m["port_range"].(string)),
		}

		networkRuleList = append(networkRuleList, networkACLRule)
	}

	return networkRuleList
}
