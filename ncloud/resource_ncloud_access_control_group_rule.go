package ncloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_access_control_group_rule", resourceNcloudAccessControlGroupRule())
}

func resourceNcloudAccessControlGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAccessControlGroupRuleCreate,
		Read:   resourceNcloudAccessControlGroupRuleRead,
		Update: resourceNcloudAccessControlGroupRuleUpdate,
		Delete: resourceNcloudAccessControlGroupRuleDelete,
		Schema: map[string]*schema.Schema{
			"access_control_group_no": {
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
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: ToDiagFunc(validation.All(
								validation.StringMatch(regexp.MustCompile(`TCP|UDP|ICMP|\b([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-4])\b`), "only TCP, UDP, ICMP and 1-254 are valid values."),
								validation.StringNotInSlice([]string{"1", "6", "17"}, false),
							)),
						},
						"port_range": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validatePortRange),
							Default:          "",
						},
						"ip_block": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IsCIDRNetwork(0, 32)),
							Default:          "",
						},
						"source_access_control_group_no": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
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
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: ToDiagFunc(validation.All(
								validation.StringMatch(regexp.MustCompile(`TCP|UDP|ICMP|\b([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-4])\b`), "only TCP, UDP, ICMP and 1-254 are valid values."),
								validation.StringNotInSlice([]string{"1", "6", "17"}, false),
							)),
						},
						"port_range": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validatePortRange),
							Default:          "",
						},
						"ip_block": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IsCIDRNetwork(0, 32)),
							Default:          "",
						},
						"source_access_control_group_no": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
	}
}

func resourceNcloudAccessControlGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_access_control_group_rule`")
	}

	d.SetId(d.Get("access_control_group_no").(string))
	log.Printf("[INFO] ACG ID: %s", d.Id())

	return resourceNcloudAccessControlGroupRuleUpdate(d, meta)
}

func resourceNcloudAccessControlGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	rules, err := getAccessControlGroupRuleList(config, d.Id())

	if err != nil {
		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == "1007000" { // Acg was not found
			d.SetId("")
		}
		return err
	}

	if len(rules) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("access_control_group_no", d.Id())

	// Create empty set for getAccessControlGroupRuleList
	iSet := schema.NewSet(schema.HashResource(resourceNcloudAccessControlGroupRule().Schema["inbound"].Elem.(*schema.Resource)), []interface{}{})
	oSet := schema.NewSet(schema.HashResource(resourceNcloudAccessControlGroupRule().Schema["outbound"].Elem.(*schema.Resource)), []interface{}{})

	for _, r := range rules {
		var protocol string
		if allowedProtocolCodes[*r.ProtocolType.Code] {
			protocol = *r.ProtocolType.Code
		} else {
			protocol = strconv.Itoa(int(*r.ProtocolType.Number))
		}

		m := map[string]interface{}{
			"protocol":                       protocol,
			"port_range":                     *r.PortRange,
			"ip_block":                       *r.IpBlock,
			"source_access_control_group_no": *r.AccessControlGroupSequence,
			"description":                    *r.AccessControlGroupRuleDescription,
		}

		if *r.AccessControlGroupRuleType.Code == "INBND" {
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

func resourceNcloudAccessControlGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("inbound") {
		if err := updateAccessControlGroupRule(d, config, "inbound"); err != nil {
			return err
		}
	}

	if d.HasChange("outbound") {
		if err := updateAccessControlGroupRule(d, config, "outbound"); err != nil {
			return err
		}
	}

	return resourceNcloudAccessControlGroupRuleRead(d, meta)
}

func resourceNcloudAccessControlGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	accessControlGroup, err := getAccessControlGroup(config, d.Id())
	if err != nil {
		return err
	}

	if accessControlGroup == nil {
		return fmt.Errorf("no matching Access Control Group: %s", d.Id())
	}

	i := d.Get("inbound").(*schema.Set)
	o := d.Get("outbound").(*schema.Set)

	if len(i.List()) > 0 {
		if err := removeAccessControlGroupRule(d, config, "inbound", accessControlGroup, expandRemoveAccessControlGroupRule(i.List())); err != nil {
			return err
		}
	}

	if len(o.List()) > 0 {
		if err := removeAccessControlGroupRule(d, config, "outbound", accessControlGroup, expandRemoveAccessControlGroupRule(o.List())); err != nil {
			return err
		}
	}

	return nil
}

func getAccessControlGroupRuleList(config *ProviderConfig, id string) ([]*vserver.AccessControlGroupRule, error) {
	reqParams := &vserver.GetAccessControlGroupRuleListRequest{
		RegionCode:           &config.RegionCode,
		AccessControlGroupNo: ncloud.String(id),
	}

	logCommonRequest("getAccessControlGroupRuleList", reqParams)
	resp, err := config.Client.vserver.V2Api.GetAccessControlGroupRuleList(reqParams)
	if err != nil {
		logErrorResponse("getAccessControlGroupRuleList", err, reqParams)
		return nil, err
	}
	logResponse("getAccessControlGroupRuleList", resp)

	return resp.AccessControlGroupRuleList, nil
}

func updateAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig, ruleType string) error {
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

	accessControlGroup, err := getAccessControlGroup(config, d.Id())
	if err != nil {
		return err
	}

	if accessControlGroup == nil {
		return fmt.Errorf("no matching Access Control Group: %s", d.Id())
	}

	removeAccessControlGroupRuleList := expandRemoveAccessControlGroupRule(remove)
	addAccessControlGroupRuleList, err := expandAddAccessControlGroupRule(add)
	if err != nil {
		return err
	}

	if len(removeAccessControlGroupRuleList) > 0 {
		if err := removeAccessControlGroupRule(d, config, ruleType, accessControlGroup, removeAccessControlGroupRuleList); err != nil {
			return err
		}
	}

	if len(addAccessControlGroupRuleList) > 0 {
		if err := addAccessControlGroupRule(d, config, ruleType, accessControlGroup, addAccessControlGroupRuleList); err != nil {
			return err
		}
	}

	return nil
}

func addAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig, ruleType string, accessControlGroup *vserver.AccessControlGroup, accessControlGroupRule []*vserver.AddAccessControlGroupRuleParameter) error {
	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		var reqParams interface{}
		if ruleType == "inbound" {
			reqParams = &vserver.AddAccessControlGroupInboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("AddAccessControlGroupInboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.AddAccessControlGroupInboundRule(reqParams.(*vserver.AddAccessControlGroupInboundRuleRequest))
		} else {
			reqParams = &vserver.AddAccessControlGroupOutboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("AddAccessControlGroupOutboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.AddAccessControlGroupOutboundRule(reqParams.(*vserver.AddAccessControlGroupOutboundRuleRequest))
		}

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == ApiErrorAcgCantChangeSameTime {
				logErrorResponse("retry AddAccessControlGroupRule", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		logErrorResponse("AddAccessControlGroupRule", err, reqParams)
		return err
	}

	logResponse("AddAccessControlGroupRule", resp)

	if err = waitForVpcAccessControlGroupRunning(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func removeAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig, ruleType string, accessControlGroup *vserver.AccessControlGroup, accessControlGroupRule []*vserver.RemoveAccessControlGroupRuleParameter) error {
	var reqParams interface{}
	var resp interface{}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error

		var reqParams interface{}
		if ruleType == "inbound" {
			reqParams = &vserver.RemoveAccessControlGroupInboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("RemoveAccessControlGroupInboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.RemoveAccessControlGroupInboundRule(reqParams.(*vserver.RemoveAccessControlGroupInboundRuleRequest))
		} else {
			reqParams = &vserver.RemoveAccessControlGroupOutboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Id()),
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: accessControlGroupRule,
			}

			logCommonRequest("RemoveAccessControlGroupOutboundRule", reqParams)
			resp, err = config.Client.vserver.V2Api.RemoveAccessControlGroupOutboundRule(reqParams.(*vserver.RemoveAccessControlGroupOutboundRuleRequest))
		}

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == ApiErrorAcgCantChangeSameTime {
				logErrorResponse("retry RemoveAccessControlGroupRule", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		logErrorResponse("RemoveAccessControlGroupRule", err, reqParams)
		return err
	}

	logResponse("RemoveAccessControlGroupRule", resp)

	if err = waitForVpcAccessControlGroupRunning(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandAddAccessControlGroupRule(rules []interface{}) ([]*vserver.AddAccessControlGroupRuleParameter, error) {
	var acgRuleList []*vserver.AddAccessControlGroupRuleParameter

	for _, vi := range rules {
		m := vi.(map[string]interface{})

		if len(m["ip_block"].(string)) == 0 && len(m["source_access_control_group_no"].(string)) == 0 {
			return nil, fmt.Errorf("one of either `ip_block` or `source_access_control_group_no` is required")
		}

		if len(m["ip_block"].(string)) > 0 && len(m["source_access_control_group_no"].(string)) > 0 {
			return nil, fmt.Errorf("cannot be specified with `ip_block` and `source_access_control_group_no`")
		}

		acgRule := &vserver.AddAccessControlGroupRuleParameter{
			ProtocolTypeCode:                  ncloud.String(m["protocol"].(string)),
			PortRange:                         ncloud.String(m["port_range"].(string)),
			IpBlock:                           ncloud.String(m["ip_block"].(string)),
			AccessControlGroupSequence:        ncloud.String(m["source_access_control_group_no"].(string)),
			AccessControlGroupRuleDescription: ncloud.String(m["description"].(string)),
		}

		acgRuleList = append(acgRuleList, acgRule)
	}

	return acgRuleList, nil
}

func expandRemoveAccessControlGroupRule(rules []interface{}) []*vserver.RemoveAccessControlGroupRuleParameter {
	var acgRuleList []*vserver.RemoveAccessControlGroupRuleParameter

	for _, vi := range rules {
		m := vi.(map[string]interface{})

		acgRule := &vserver.RemoveAccessControlGroupRuleParameter{
			IpBlock:                    ncloud.String(m["ip_block"].(string)),
			AccessControlGroupSequence: ncloud.String(m["source_access_control_group_no"].(string)),
			ProtocolTypeCode:           ncloud.String(m["protocol"].(string)),
			PortRange:                  ncloud.String(m["port_range"].(string)),
		}

		acgRuleList = append(acgRuleList, acgRule)
	}

	return acgRuleList
}

var allowedProtocolCodes = map[string]bool{
	"TCP":  true,
	"UDP":  true,
	"ICMP": true,
}
