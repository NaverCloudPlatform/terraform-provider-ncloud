package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access control group name to search",
			},
			"is_default_group": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether default group",
			},
			"source_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to apply to the source access control rule list returned by ncloud",
			},
			"source_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source IP",
			},
			"destination_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Destination Port",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Protocol type code",
			},
			"configuration_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access control rule configuration no",
			},
			"source_configuration_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source access control rule configuration no",
			},
			"source_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source access control rule name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access control rule description",
			},
		},
	}
}

func dataSourceNcloudAccessControlRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	var allAccessControlRuleList []*server.AccessControlRule
	configNo, configNoOk := d.GetOk("access_control_group_configuration_no")
	acgName, acgNameOk := d.GetOk("access_control_group_name")
	isDefaultGroup, isDefaultGroupOk := d.GetOk("is_default_group")

	if !configNoOk && !acgNameOk && !isDefaultGroupOk {
		return fmt.Errorf("either `access_control_group_configuration_no` or `access_control_group_name` or `is_default_group` must be defined")
	}

	if !configNoOk {
		reqParams := new(server.GetAccessControlGroupListRequest)
		if acgNameOk {
			reqParams.AccessControlGroupName = ncloud.String(acgName.(string))
		}
		if isDefaultGroupOk {
			reqParams.IsDefault = ncloud.Bool(isDefaultGroup.(bool))
		}
		acgResp, err := getAccessControlGroupList(client, reqParams)
		if err != nil {
			return err
		}
		for _, acg := range acgResp.AccessControlGroupList {
			resp, err := getAccessControlRuleList(client, ncloud.StringValue(acg.AccessControlGroupConfigurationNo))
			if err != nil {
				return err
			}
			for _, rule := range resp.AccessControlRuleList {
				allAccessControlRuleList = append(allAccessControlRuleList, rule)
			}
		}
	} else {
		groupConfigNo := configNo.(string)
		resp, err := getAccessControlRuleList(client, groupConfigNo)
		if err != nil {
			return err
		}
		for _, rule := range resp.AccessControlRuleList {
			allAccessControlRuleList = append(allAccessControlRuleList, rule)
		}
	}

	var filteredAccessControlRuleList []*server.AccessControlRule
	var accessControlRule *server.AccessControlRule

	var r *regexp.Regexp
	nameRegex, nameRegexOk := d.GetOk("source_name_regex")
	sourceIP, sourceIPOk := d.GetOk("source_ip")
	destinationPort, destinationPortOk := d.GetOk("destination_port")
	protocolTypeCode, protocolTypeCodeOk := d.GetOk("protocol_type")

	if nameRegexOk || sourceIPOk || destinationPortOk || protocolTypeCodeOk {
		if nameRegexOk {
			r = regexp.MustCompile(nameRegex.(string))
		}

		for _, rule := range allAccessControlRuleList {
			if nameRegexOk && r.MatchString(ncloud.StringValue(rule.SourceAccessControlRuleName)) {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			} else if sourceIPOk && sourceIP == ncloud.StringValue(rule.SourceIp) {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			} else if destinationPortOk && destinationPort == ncloud.StringValue(rule.DestinationPort) {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			} else if protocolTypeCodeOk && protocolTypeCode == ncloud.StringValue(rule.ProtocolType.Code) {
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

func getAccessControlRuleList(client *NcloudAPIClient, groupConfigNo string) (*server.GetAccessControlRuleListResponse, error) {
	reqParams := server.GetAccessControlRuleListRequest{
		AccessControlGroupConfigurationNo: ncloud.String(groupConfigNo),
	}

	logCommonRequest("GetAccessControlRuleList", reqParams)
	resp, err := client.server.V2Api.GetAccessControlRuleList(&reqParams)
	if err != nil {
		logErrorResponse("GetAccessControlRuleList", err, groupConfigNo)
		return nil, err
	}
	logCommonResponse("GetAccessControlRuleList", GetCommonResponse(resp))
	return resp, nil
}

func accessControlRuleAttributes(d *schema.ResourceData, accessControlRule *server.AccessControlRule) error {
	d.SetId(*accessControlRule.AccessControlRuleConfigurationNo)
	d.Set("configuration_no", accessControlRule.AccessControlRuleConfigurationNo)
	d.Set("source_ip", accessControlRule.SourceIp)
	d.Set("destination_port", accessControlRule.DestinationPort)
	d.Set("source_configuration_no", accessControlRule.SourceAccessControlRuleConfigurationNo)
	d.Set("source_name", accessControlRule.SourceAccessControlRuleName)
	d.Set("description", accessControlRule.AccessControlRuleDescription)

	if protocolType := flattenCommonCode(accessControlRule.ProtocolType); protocolType["code"] != nil {
		d.Set("protocol_type", protocolType["code"])
	}

	return nil
}
