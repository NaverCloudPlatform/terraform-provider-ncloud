package ncloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudAccessControlRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAccessControlRulesRead,

		Schema: map[string]*schema.Schema{
			"access_control_group_configuration_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access control group setting number to search",
			},
			"source_access_control_rule_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp,
			},
			"access_control_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_control_rule_configuration_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     commonCodeSchemaResource,
						},
						"source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_port": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudAccessControlRulesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	d.SetId(time.Now().UTC().String())

	id := d.Get("access_control_group_configuration_no").(string)
	resp, err := conn.GetAccessControlRuleList(id)
	if err != nil {
		logErrorResponse("GetAccessControlRuleList", err, id)
	}
	logCommonResponse("GetAccessControlRuleList", id, resp.CommonResponse)

	allAccessControlRuleList := resp.AccessControlRuleList
	var filtereAccessControlRuleList []sdk.AccessControlRule
	nameRegex, nameRegexOk := d.GetOk("source_access_control_rule_name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, serverImage := range allAccessControlRuleList {
			if r.MatchString(serverImage.SourceAccessControlRuleName) {
				filtereAccessControlRuleList = append(filtereAccessControlRuleList, serverImage)
			}
		}
	} else {
		filtereAccessControlRuleList = allAccessControlRuleList[:]
	}

	if len(filtereAccessControlRuleList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return accessControlRulesAttributes(d, filtereAccessControlRuleList)
}

func accessControlRulesAttributes(d *schema.ResourceData, accessControlRules []sdk.AccessControlRule) error {
	var ids []string
	var s []map[string]interface{}
	for _, accessControlRule := range accessControlRules {
		mapping := map[string]interface{}{
			"access_control_rule_configuration_no":        accessControlRule.AccessControlRuleConfigurationNo,
			"protocol_type":                               setCommonCode(accessControlRule.ProtocolType),
			"source_ip":                                   accessControlRule.SourceIP,
			"destination_port":                            accessControlRule.DestinationPort,
			"source_access_control_rule_configuration_no": accessControlRule.SourceAccessControlRuleConfigurationNo,
			"source_access_control_rule_name":             accessControlRule.SourceAccessControlRuleName,
			"access_control_rule_description":             accessControlRule.AccessControlRuleDescription,
		}

		ids = append(ids, accessControlRule.AccessControlRuleConfigurationNo)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("access_control_rules", s); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
