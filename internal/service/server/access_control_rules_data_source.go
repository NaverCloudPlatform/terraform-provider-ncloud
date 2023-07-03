package server

import (
	"fmt"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudAccessControlRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAccessControlRulesRead,

		Schema: map[string]*schema.Schema{
			"access_control_group_configuration_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access control group setting number to search",
			},
			"source_name_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringIsValidRegExp),
				Description:      "A regex string to apply to the ACG rule list returned by ncloud",
			},
			"access_control_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of access control rules configuration no",
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudAccessControlRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client
	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		return NotSupportVpc("data source `ncloud_access_control_rules`")
	}

	d.SetId(time.Now().UTC().String())

	id := d.Get("access_control_group_configuration_no").(string)
	reqParams := server.GetAccessControlRuleListRequest{AccessControlGroupConfigurationNo: ncloud.String(id)}

	LogCommonRequest("GetAccessControlRuleList", reqParams)

	resp, err := client.Server.V2Api.GetAccessControlRuleList(&reqParams)
	if err != nil {
		LogErrorResponse("GetAccessControlRuleList", err, id)
		return err
	}
	LogCommonResponse("GetAccessControlRuleList", GetCommonResponse(resp))

	allAccessControlRuleList := resp.AccessControlRuleList
	var filteredAccessControlRuleList []*server.AccessControlRule
	nameRegex, nameRegexOk := d.GetOk("source_name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, rule := range allAccessControlRuleList {
			if r.MatchString(ncloud.StringValue(rule.SourceAccessControlRuleName)) {
				filteredAccessControlRuleList = append(filteredAccessControlRuleList, rule)
			}
		}
	} else {
		filteredAccessControlRuleList = allAccessControlRuleList[:]
	}

	if len(filteredAccessControlRuleList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return accessControlRulesAttributes(d, filteredAccessControlRuleList)
}

func accessControlRulesAttributes(d *schema.ResourceData, accessControlRules []*server.AccessControlRule) error {
	var ids []string

	for _, accessControlRule := range accessControlRules {
		ids = append(ids, ncloud.StringValue(accessControlRule.AccessControlRuleConfigurationNo))
	}
	d.SetId(DataResourceIdHash(ids))

	if err := d.Set("access_control_rules", flattenAccessControlRules(accessControlRules)); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return WriteToFile(output.(string), d.Get("access_control_rules"))
	}

	return nil
}
