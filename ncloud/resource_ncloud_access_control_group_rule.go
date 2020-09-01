package ncloud

import (
	"bytes"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"regexp"
	"strings"
	"time"
)

func resourceNcloudAccessControlGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAccessControlGroupRuleCreate,
		Read:   resourceNcloudAccessControlGroupRuleRead,
		Update: resourceNcloudAccessControlGroupRuleUpdate,
		Delete: resourceNcloudAccessControlGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ":")
				log.Printf("[INFO] idParts: %s", idParts)
				if len(idParts) != 5 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" || idParts[4] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected ACCESS_CONTROL_GROUP_NO:RULE_TYPE:PROTOCOL:ACCESS_SOURCE(IP_BLOCK or ACG_NO):PORT_RANGE", d.Id())
				}

				acgNo := idParts[0]
				ruleType := idParts[1]
				protocol := idParts[2]
				accessSource := idParts[3]
				portRange := idParts[4]

				rule := &AccessControlGroupRuleParam{
					AccessControlGroupNo: acgNo,
					RuleType:             ruleType,
					Protocol:             protocol,
					PortRange:            portRange,
				}

				d.Set("access_control_group_no", rule.AccessControlGroupNo)
				d.Set("rule_type", rule.RuleType)
				d.Set("protocol", rule.Protocol)

				if regexp.MustCompile(`^\d+$`).MatchString(accessSource) {
					d.Set("source_access_control_group_no", accessSource)
					rule.SourceAccessControlGroup = accessSource
				} else {
					d.Set("ip_block", accessSource)
					rule.IpBlock = accessSource
				}

				d.Set("port_range", rule.PortRange)

				d.SetId(accessControlGroupRuleHash(rule))

				return []*schema.ResourceData{d}, nil
			},
		},
		CustomizeDiff: ncloudVpcCommonCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"access_control_group_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"INBND", "OTBND"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
			},
			"ip_block": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.IsCIDRNetwork(0, 32),
				ConflictsWith: []string{"source_access_control_group_no"},
			},
			"source_access_control_group_no": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ip_block"},
			},
			"port_range": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePortRange,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
	}
}

//AccessControlGroupRuleParam struct for ACG rule
type AccessControlGroupRuleParam struct {
	AccessControlGroupNo     string
	RuleType                 string
	Protocol                 string
	IpBlock                  string
	SourceAccessControlGroup string
	PortRange                string
}

func resourceNcloudAccessControlGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := createAccessControlGroupRule(d, config)

	if err != nil {
		return err
	}

	d.SetId(accessControlGroupRuleHash(&AccessControlGroupRuleParam{
		AccessControlGroupNo:     *instance.AccessControlGroupNo,
		RuleType:                 *instance.AccessControlGroupRuleType.Code,
		Protocol:                 *instance.ProtocolType.Code,
		IpBlock:                  *instance.IpBlock,
		SourceAccessControlGroup: *instance.AccessControlGroupSequence,
		PortRange:                *instance.PortRange,
	}))

	log.Printf("[INFO] ACG ID: %s", d.Id())

	return resourceNcloudAccessControlGroupRuleRead(d, meta)
}

func resourceNcloudAccessControlGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	param := convAccessControlGroupRuleParam(d)
	instance, err := getAccessControlGroupRule(config, param)
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(accessControlGroupRuleHash(&AccessControlGroupRuleParam{
		AccessControlGroupNo:     *instance.AccessControlGroupNo,
		RuleType:                 *instance.AccessControlGroupRuleType.Code,
		Protocol:                 *instance.ProtocolType.Code,
		IpBlock:                  *instance.IpBlock,
		SourceAccessControlGroup: *instance.AccessControlGroupSequence,
		PortRange:                *instance.PortRange,
	}))

	d.Set("access_control_group_no", instance.AccessControlGroupNo)
	d.Set("rule_type", instance.AccessControlGroupRuleType.Code)
	d.Set("protocol", instance.ProtocolType.Code)
	d.Set("ip_block", instance.IpBlock)
	d.Set("source_access_control_group_no", instance.AccessControlGroupSequence)
	d.Set("port_range", instance.PortRange)
	d.Set("description", instance.AccessControlGroupRuleDescription)

	return nil
}

func resourceNcloudAccessControlGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudAccessControlGroupRuleRead(d, meta)
}

func resourceNcloudAccessControlGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if config.SupportVPC {
		rule := convAccessControlGroupRuleParam(d)
		if err := deleteAccessControlGroupRule(d, config, rule); err != nil {
			return err
		}
	} else {
		return NotSupportClassic("resource `ncloud_access_control_group_rule`")
	}

	return nil
}

func getAccessControlGroupRule(config *ProviderConfig, rule *AccessControlGroupRuleParam) (*vserver.AccessControlGroupRule, error) {
	if config.SupportVPC {
		return getVpcAccessControlGroupRule(config, rule)
	}

	return nil, NotSupportClassic("resource `ncloud_access_control_group_rule`")
}

func getVpcAccessControlGroupRule(config *ProviderConfig, rule *AccessControlGroupRuleParam) (*vserver.AccessControlGroupRule, error) {
	reqParams := &vserver.GetAccessControlGroupRuleListRequest{
		RegionCode:                     &config.RegionCode,
		AccessControlGroupNo:           &rule.AccessControlGroupNo,
		AccessControlGroupRuleTypeCode: &rule.RuleType,
	}

	logCommonRequest("getVpcAccessControlGroupRule", reqParams)
	resp, err := config.Client.vserver.V2Api.GetAccessControlGroupRuleList(reqParams)
	if err != nil {
		logErrorResponse("getVpcAccessControlGroupRule", err, reqParams)
		return nil, err
	}
	logResponse("getVpcAccessControlGroupRule", resp)

	if resp.AccessControlGroupRuleList != nil {
		for _, i := range resp.AccessControlGroupRuleList {
			if *i.ProtocolType.Code == rule.Protocol &&
				*i.IpBlock == rule.IpBlock &&
				*i.AccessControlGroupSequence == rule.SourceAccessControlGroup &&
				*i.PortRange == rule.PortRange {
				return i, nil
			}
		}
		return nil, nil
	}

	return nil, nil
}

func createAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig) (*vserver.AccessControlGroupRule, error) {
	if config.SupportVPC {
		return createVpcAccessControlGroupRule(d, config)
	}

	return nil, NotSupportClassic("resource `ncloud_access_control_group_rule`")
}

func createVpcAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig) (*vserver.AccessControlGroupRule, error) {
	accessControlGroup, err := getAccessControlGroup(config, d.Get("access_control_group_no").(string))
	if err != nil {
		return nil, err
	}

	if accessControlGroup == nil {
		return nil, fmt.Errorf("no matching Access Control Group: %s", d.Get("access_control_group_no"))
	}

	accessControlGroupRule := &vserver.AddAccessControlGroupRuleParameter{
		AccessControlGroupRuleDescription: ncloud.String(d.Get("description").(string)),
		IpBlock:                           ncloud.String(d.Get("ip_block").(string)),
		AccessControlGroupSequence:        ncloud.String(d.Get("source_access_control_group_no").(string)),
		PortRange:                         ncloud.String(d.Get("port_range").(string)),
		ProtocolTypeCode:                  ncloud.String(d.Get("protocol").(string)),
	}

	var resp interface{}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var reqParams interface{}
		if d.Get("rule_type").(string) == "INBND" {
			reqParams = &vserver.AddAccessControlGroupInboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Get("access_control_group_no").(string)),
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: []*vserver.AddAccessControlGroupRuleParameter{accessControlGroupRule},
			}

			logCommonRequest("createVpcAccessControlGroupRule", reqParams)
			resp, err = config.Client.vserver.V2Api.AddAccessControlGroupInboundRule(reqParams.(*vserver.AddAccessControlGroupInboundRuleRequest))
		} else {
			reqParams = &vserver.AddAccessControlGroupOutboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       ncloud.String(d.Get("access_control_group_no").(string)),
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: []*vserver.AddAccessControlGroupRuleParameter{accessControlGroupRule},
			}

			logCommonRequest("createVpcAccessControlGroupRule", reqParams)
			resp, err = config.Client.vserver.V2Api.AddAccessControlGroupOutboundRule(reqParams.(*vserver.AddAccessControlGroupOutboundRuleRequest))
		}

		if err == nil {
			return resource.NonRetryableError(err)
		}

		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == ApiErrorAcgCantChangeSameTime {
			logErrorResponse("retry createVpcAccessControlGroupRule", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return nil, err
	}

	logResponse("createVpcAccessControlGroupRule", resp)

	var instance *vserver.AccessControlGroupRule
	if d.Get("rule_type").(string) == "INBND" {
		instance = resp.(*vserver.AddAccessControlGroupInboundRuleResponse).AccessControlGroupRuleList[0]
	} else {
		instance = resp.(*vserver.AddAccessControlGroupOutboundRuleResponse).AccessControlGroupRuleList[0]
	}

	if err := waitForVpcAccessControlGroupRunning(config, d.Get("access_control_group_no").(string)); err != nil {
		return nil, err
	}

	return instance, nil
}

func deleteAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig, rule *AccessControlGroupRuleParam) error {
	if config.SupportVPC {
		return deleteVpcAccessControlGroupRule(d, config, rule)
	}

	return NotSupportClassic("resource `ncloud_access_control_group_rule`")
}

func deleteVpcAccessControlGroupRule(d *schema.ResourceData, config *ProviderConfig, rule *AccessControlGroupRuleParam) error {
	accessControlGroup, err := getAccessControlGroup(config, rule.AccessControlGroupNo)
	if err != nil {
		return err
	}

	if accessControlGroup == nil {
		return fmt.Errorf("no matching Access Control Group: %s", rule.AccessControlGroupNo)
	}

	accessControlGroupRule := &vserver.RemoveAccessControlGroupRuleParameter{
		IpBlock:                    &rule.IpBlock,
		AccessControlGroupSequence: &rule.SourceAccessControlGroup,
		PortRange:                  &rule.PortRange,
		ProtocolTypeCode:           &rule.Protocol,
	}

	var reqParams interface{}
	var resp interface{}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		if rule.RuleType == "INBND" {
			reqParams := &vserver.RemoveAccessControlGroupInboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       &rule.AccessControlGroupNo,
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: []*vserver.RemoveAccessControlGroupRuleParameter{accessControlGroupRule},
			}

			logCommonRequest("deleteVpcAccessControlGroupRule", reqParams)
			resp, err = config.Client.vserver.V2Api.RemoveAccessControlGroupInboundRule(reqParams)
		} else {
			reqParams := &vserver.RemoveAccessControlGroupOutboundRuleRequest{
				RegionCode:                 &config.RegionCode,
				AccessControlGroupNo:       &rule.AccessControlGroupNo,
				VpcNo:                      accessControlGroup.VpcNo,
				AccessControlGroupRuleList: []*vserver.RemoveAccessControlGroupRuleParameter{accessControlGroupRule},
			}

			logCommonRequest("deleteVpcAccessControlGroupRule", reqParams)
			resp, err = config.Client.vserver.V2Api.RemoveAccessControlGroupOutboundRule(reqParams)
		}

		if err == nil {
			return resource.NonRetryableError(err)
		}

		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == ApiErrorAcgCantChangeSameTime {
			logErrorResponse("retry deleteVpcAccessControlGroupRule", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("deleteVpcAccessControlGroupRule", err, reqParams)
		return err
	}

	logResponse("deleteVpcAccessControlGroupRule", resp)

	if err := waitForVpcAccessControlGroupRunning(config, rule.AccessControlGroupNo); err != nil {
		return err
	}

	return nil
}

func accessControlGroupRuleHash(rule *AccessControlGroupRuleParam) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s-", rule.AccessControlGroupNo))
	buf.WriteString(fmt.Sprintf("%s-", rule.RuleType))
	buf.WriteString(fmt.Sprintf("%s-", rule.Protocol))
	buf.WriteString(fmt.Sprintf("%s-", rule.IpBlock))
	if len(rule.SourceAccessControlGroup) > 0 {
		buf.WriteString(fmt.Sprintf("%s-", rule.SourceAccessControlGroup))
	}
	buf.WriteString(fmt.Sprintf("%s-", rule.PortRange))
	return fmt.Sprintf("acgr-%d", hashcode.String(buf.String()))
}

func convAccessControlGroupRuleParam(d *schema.ResourceData) *AccessControlGroupRuleParam {
	return &AccessControlGroupRuleParam{
		AccessControlGroupNo:     d.Get("access_control_group_no").(string),
		RuleType:                 d.Get("rule_type").(string),
		Protocol:                 d.Get("protocol").(string),
		IpBlock:                  d.Get("ip_block").(string),
		SourceAccessControlGroup: d.Get("source_access_control_group_no").(string),
		PortRange:                d.Get("port_range").(string),
	}
}
