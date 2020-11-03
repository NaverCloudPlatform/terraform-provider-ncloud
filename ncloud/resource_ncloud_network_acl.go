package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterResource("ncloud_network_acl", resourceNcloudNetworkACL())
}

func resourceNcloudNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNetworkACLCreate,
		Read:   resourceNcloudNetworkACLRead,
		Update: resourceNcloudNetworkACLUpdate,
		Delete: resourceNcloudNetworkACLDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: ncloudVpcCommonCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_acl_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inbound": {
				Type:       schema.TypeSet,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 199),
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
						},
						"port_range": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validatePortRange,
						},
						"rule_action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DROP"}, false),
						},
						"ip_block": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 32),
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 1000),
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
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 199),
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
						},
						"port_range": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validatePortRange,
						},
						"rule_action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DROP"}, false),
						},
						"ip_block": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 32),
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 1000),
						},
					},
				},
			},
		},
	}
}

func resourceNcloudNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_network_acl`")
	}

	reqParams := &vpc.CreateNetworkAclRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(d.Get("vpc_no").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NetworkAclName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		reqParams.NetworkAclDescription = ncloud.String(v.(string))
	}

	logCommonRequest("CreateNetworkAcl", reqParams)
	resp, err := config.Client.vpc.V2Api.CreateNetworkAcl(reqParams)
	if err != nil {
		logErrorResponse("CreateNetworkAcl", err, reqParams)
		return err
	}

	logResponse("CreateNetworkAcl", resp)

	instance := resp.NetworkAclList[0]
	d.SetId(*instance.NetworkAclNo)
	log.Printf("[INFO] Network ACL ID: %s", d.Id())

	if err := waitForNcloudNetworkACLCreation(config, d.Id()); err != nil {
		return err
	}

	return resourceNcloudNetworkACLUpdate(d, meta)
}

func resourceNcloudNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getNetworkACLInstance(config, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.NetworkAclNo)
	d.Set("network_acl_no", instance.NetworkAclNo)
	d.Set("name", instance.NetworkAclName)
	d.Set("description", instance.NetworkAclDescription)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("is_default", instance.IsDefault)

	rules, err := getNetworkACLRuleList(d, config)
	if err != nil {
		return err
	}

	var inbound []map[string]interface{}
	var outbound []map[string]interface{}
	for _, r := range rules {
		m := map[string]interface{}{
			"priority":    *r.Priority,
			"protocol":    *r.ProtocolType.Code,
			"port_range":  *r.PortRange,
			"rule_action": *r.RuleAction.Code,
			"ip_block":    *r.IpBlock,
			"description": *r.NetworkAclRuleDescription,
		}

		if *r.NetworkAclRuleType.Code == "INBND" {
			inbound = append(inbound, m)
		} else {
			outbound = append(outbound, m)
		}
	}

	if err := d.Set("inbound", inbound); err != nil {
		log.Printf("[WARN] Error setting inbound rule set for (%s): %s", d.Id(), err)
	}

	if err := d.Set("outbound", outbound); err != nil {
		log.Printf("[WARN] Error setting outbound rule set for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNcloudNetworkACLUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("description") {
		if err := setNetworkACLDescription(d, config); err != nil {
			return err
		}
	}

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

	return resourceNcloudNetworkACLRead(d, meta)
}

func resourceNcloudNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.DeleteNetworkAclRequest{
		RegionCode:   &config.RegionCode,
		NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
	}

	logCommonRequest("DeleteNetworkAcl", reqParams)
	resp, err := config.Client.vpc.V2Api.DeleteNetworkAcl(reqParams)
	if err != nil {
		logErrorResponse("DeleteNetworkAcl", err, reqParams)
		return err
	}

	logResponse("DeleteNetworkAcl", resp)

	if err := waitForNcloudNetworkACLDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForNcloudNetworkACLCreation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkACLInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclStatus")
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Network ACL (%s) to become available: %s", id, err)
	}

	return nil
}

func waitForNcloudNetworkACLDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
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

func getNetworkACLInstance(config *ProviderConfig, id string) (*vpc.NetworkAcl, error) {
	reqParams := &vpc.GetNetworkAclDetailRequest{
		RegionCode:   &config.RegionCode,
		NetworkAclNo: ncloud.String(id),
	}

	logCommonRequest("GetNetworkAclDetail", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNetworkAclDetail(reqParams)
	if err != nil {
		logErrorResponse("GetNetworkAclDetail", err, reqParams)
		return nil, err
	}
	logResponse("GetNetworkAclDetail", resp)

	if len(resp.NetworkAclList) > 0 {
		instance := resp.NetworkAclList[0]
		return instance, nil
	}

	return nil, nil
}

func setNetworkACLDescription(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vpc.SetNetworkAclDescriptionRequest{
		RegionCode:            &config.RegionCode,
		NetworkAclNo:          ncloud.String(d.Id()),
		NetworkAclDescription: StringPtrOrNil(d.GetOk("description")),
	}

	logCommonRequest("setNetworkAclDescription", reqParams)
	resp, err := config.Client.vpc.V2Api.SetNetworkAclDescription(reqParams)
	if err != nil {
		logErrorResponse("setNetworkAclDescription", err, reqParams)
		return err
	}
	logResponse("setNetworkAclDescription", resp)

	return nil
}

func getNetworkACLRuleList(d *schema.ResourceData, config *ProviderConfig) ([]*vpc.NetworkAclRule, error) {
	reqParams := &vpc.GetNetworkAclRuleListRequest{
		RegionCode:   &config.RegionCode,
		NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
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

func updateNetworkACLRule(d *schema.ResourceData, meta interface{}, ruleType string) error {
	config := meta.(*ProviderConfig)

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

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
			IpBlock:          ncloud.String(m["ip_block"].(string)),
			RuleActionCode:   ncloud.String(m["rule_action"].(string)),
			Priority:         ncloud.Int32(int32(m["priority"].(int))),
			ProtocolTypeCode: ncloud.String(m["protocol"].(string)),
		}

		if v, ok := m["port_range"]; ok {
			networkACLRule.PortRange = ncloud.String(v.(string))
		}

		if v, ok := m["description"]; ok {
			networkACLRule.NetworkAclRuleDescription = ncloud.String(v.(string))
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
			RuleActionCode:   ncloud.String(m["rule_action"].(string)),
			Priority:         ncloud.Int32(int32(m["priority"].(int))),
			ProtocolTypeCode: ncloud.String(m["protocol"].(string)),
		}

		if v, ok := m["port_range"]; ok {
			networkACLRule.PortRange = ncloud.String(v.(string))
		}

		networkRuleList = append(networkRuleList, networkACLRule)
	}

	return networkRuleList
}
