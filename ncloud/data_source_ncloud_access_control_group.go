package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudAccessControlGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAccessControlGroupRead,

		Schema: map[string]*schema.Schema{
			"access_control_group_configuration_no": {
				Type:     schema.TypeString,
				Optional: true,
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
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"access_control_group_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudAccessControlGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	configNo, configNoOk := d.GetOk("access_control_group_configuration_no")
	acgName, acgNameOk := d.GetOk("access_control_group_name")
	mostRecent, mostRecentOk := d.GetOk("most_recent")

	if !configNoOk && !acgNameOk && !mostRecentOk {
		return fmt.Errorf("either access_control_group_configuration_no or access_control_group_name or most_recent is required")
	}

	reqParams := new(sdk.RequestAccessControlGroupList)
	if configNoOk {
		reqParams.AccessControlGroupConfigurationNoList = []string{configNo.(string)}
	}
	if acgNameOk {
		reqParams.AccessControlGroupName = acgName.(string)
	}
	reqParams.IsDefault = d.Get("is_default_group").(string)
	reqParams.PageNo = 1

	resp, err := conn.GetAccessControlGroupList(reqParams)
	if err != nil {
		logErrorResponse("GetAccessControlGroupList", err, reqParams)
		return err
	}
	logCommonResponse("GetAccessControlGroupList", reqParams, resp.CommonResponse)

	var accessControlGroup sdk.AccessControlGroup
	var accessControlGroups []sdk.AccessControlGroup

	for _, acg := range resp.AccessControlGroup {
		accessControlGroups = append(accessControlGroups, acg)
	}

	if len(accessControlGroups) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	if len(accessControlGroups) > 1 && mostRecent.(bool) {
		// Query returned single result.
		accessControlGroup = mostRecentAccessControlGroup(accessControlGroups)
	} else {
		accessControlGroup = accessControlGroups[0]
	}

	return accessControlGroupAttributes(d, accessControlGroup)
}

func accessControlGroupAttributes(d *schema.ResourceData, accessControlGroup sdk.AccessControlGroup) error {
	d.SetId(string(accessControlGroup.AccessControlGroupConfigurationNo))
	d.Set("access_control_group_configuration_no", accessControlGroup.AccessControlGroupConfigurationNo)
	d.Set("access_control_group_name", accessControlGroup.AccessControlGroupName)
	d.Set("access_control_group_description", accessControlGroup.AccessControlGroupDescription)
	d.Set("is_default_group", accessControlGroup.IsDefault)
	d.Set("create_date", accessControlGroup.CreateDate)

	return nil
}
