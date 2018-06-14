package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudAccessControlGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAccessControlGroupsRead,

		Schema: map[string]*schema.Schema{
			"access_control_group_configuration_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"is_default_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateBoolValue,
			},
			"access_control_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_no": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"access_control_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_control_group_configuration_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_control_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_control_group_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default_group": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_date": {
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

func dataSourceNcloudAccessControlGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	d.SetId(time.Now().UTC().String())

	reqParams := new(sdk.RequestAccessControlGroupList)
	var paramAccessControlGroupConfigurationNoList []string
	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		paramAccessControlGroupConfigurationNoList = StringList(param.([]interface{}))
	}
	reqParams.AccessControlGroupConfigurationNoList = paramAccessControlGroupConfigurationNoList

	reqParams.AccessControlGroupName = d.Get("access_control_group_name").(string)
	reqParams.IsDefault = d.Get("is_default_group").(string)
	reqParams.PageNo = d.Get("page_no").(int)

	resp, err := conn.GetAccessControlGroupList(reqParams)
	if err != nil {
		return err
	}

	var accessControlGroups []sdk.AccessControlGroup

	for _, group := range resp.AccessControlGroup {
		accessControlGroups = append(accessControlGroups, group)
	}

	if len(accessControlGroups) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return accessControlGroupsAttributes(d, accessControlGroups)
}

func accessControlGroupsAttributes(d *schema.ResourceData, accessControlGroups []sdk.AccessControlGroup) error {
	var ids []string
	var s []map[string]interface{}
	for _, accessControlGroup := range accessControlGroups {
		mapping := map[string]interface{}{
			"access_control_group_configuration_no": accessControlGroup.AccessControlGroupConfigurationNo,
			"access_control_group_name":             accessControlGroup.AccessControlGroupName,
			"access_control_group_description":      accessControlGroup.AccessControlGroupDescription,
			"is_default_group":                      accessControlGroup.IsDefault,
			"create_date":                           accessControlGroup.CreateDate,
		}

		ids = append(ids, string(accessControlGroup.AccessControlGroupConfigurationNo))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("access_control_groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
