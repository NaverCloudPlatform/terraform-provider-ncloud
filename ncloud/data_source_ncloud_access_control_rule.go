package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
)

func dataSourceNcloudAccessControlRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAccessControlRuleRead,

		Schema: map[string]*schema.Schema{
			"access_control_group_configuration_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Access control group setting number to search",
			},
			"access_control_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateBoolValue,
			},

			"source_access_control_rule_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol_type_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"access_control_rule_configuration_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"source_access_control_rule_configuration_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_access_control_rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_control_rule_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudAccessControlRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	var allAccessControlRuleList []sdk.AccessControlRule
	configNo, configNoOk := d.GetOk("access_control_group_configuration_no")
	acgName, acgNameOk := d.GetOk("access_control_group_name")
	isDefaultGroup, isDefaultGroupOk := d.GetOk("is_default_group")

	if !configNoOk && !acgNameOk && !isDefaultGroupOk {
		return fmt.Errorf("either `access_control_group_configuration_no` or `access_control_group_name` or `is_default_group` must be defined")
	}

	if !configNoOk {
		reqParams := new(sdk.RequestAccessControlGroupList)
		if acgNameOk {
			reqParams.AccessControlGroupName = acgName.(string)
		}
		if isDefaultGroupOk {
			reqParams.IsDefault = isDefaultGroup.(string)
		}
		acgResp, err := getAccessControlGroupList(conn, reqParams)
		if err != nil {
			return err
		}
		for _, acg := range acgResp.AccessControlGroup {
			resp, err := getAccessControlRuleList(conn, acg.AccessControlGroupConfigurationNo)
			if err != nil {
				return err
			}
			for _, rule := range resp.AccessControlRuleList {
				allAccessControlRuleList = append(allAccessControlRuleList, rule)
			}
		}
	} else {
		groupConfigNo := configNo.(string)
		resp, err := getAccessControlRuleList(conn, groupConfigNo)
		if err != nil {
			return err
		}
		for _, rule := range resp.AccessControlRuleList {
			allAccessControlRuleList = append(allAccessControlRuleList, rule)
		}
	}

	var filteredAccessControlRuleList []sdk.AccessControlRule
	var accessControlRule sdk.AccessControlRule

	var r *regexp.Regexp
	nameRegex, nameRegexOk := d.GetOk("source_access_control_rule_name_regex")
	sourceIP, sourceIPOk := d.GetOk("source_ip")
	destinationPort, destinationPortOk := d.GetOk("destination_port")
	protocolTypeCode, protocolTypeCodeOk := d.GetOk("protocol_type_code")

	if nameRegexOk || sourceIPOk || destinationPortOk || protocolTypeCodeOk {
		if nameRegexOk {
			r = regexp.MustCompile(nameRegex.(string))
		}

		for _, rule := range allAccessControlRuleList {
			if nameRegexOk && r.MatchString(rule.SourceAccessControlRuleName) {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			} else if sourceIPOk && sourceIP == rule.SourceIP {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			} else if destinationPortOk && destinationPort == rule.DestinationPort {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			} else if protocolTypeCodeOk && protocolTypeCode == rule.ProtocolType.Code {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			}
		}
	} else {
		filteredAccessControlRuleList = allAccessControlRuleList[:]
	}

	if len(filteredAccessControlRuleList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	accessControlRule = filteredAccessControlRuleList[0]
	return accessControlRuleAttributes(d, accessControlRule)
}

func getAccessControlRuleList(conn *sdk.Conn, groupConfigNo string) (*sdk.AccessControlRuleList, error) {
	resp, err := conn.GetAccessControlRuleList(groupConfigNo)
	if err != nil {
		logErrorResponse("GetAccessControlRuleList", err, groupConfigNo)
		return nil, err
	}
	logCommonResponse("GetAccessControlRuleList", groupConfigNo, resp.CommonResponse)
	return resp, nil
}

func accessControlRuleAttributes(d *schema.ResourceData, accessControlRule sdk.AccessControlRule) error {
	d.SetId(accessControlRule.AccessControlRuleConfigurationNo)
	d.Set("access_control_rule_configuration_no", accessControlRule.AccessControlRuleConfigurationNo)
	d.Set("protocol_type", setCommonCode(accessControlRule.ProtocolType))
	d.Set("source_ip", accessControlRule.SourceIP)
	d.Set("destination_port", accessControlRule.DestinationPort)
	d.Set("source_access_control_rule_configuration_no", accessControlRule.SourceAccessControlRuleConfigurationNo)
	d.Set("source_access_control_rule_name", accessControlRule.SourceAccessControlRuleName)
	d.Set("access_control_rule_description", accessControlRule.AccessControlRuleDescription)

	return nil
}
