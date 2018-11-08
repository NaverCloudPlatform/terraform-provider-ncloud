package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
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
				Type:     schema.TypeBool,
				Optional: true,
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
	client := meta.(*NcloudAPIClient)

	configNo, configNoOk := d.GetOk("access_control_group_configuration_no")
	acgName, acgNameOk := d.GetOk("access_control_group_name")
	mostRecent, mostRecentOk := d.GetOk("most_recent")

	if !configNoOk && !acgNameOk && !mostRecentOk {
		return fmt.Errorf("either access_control_group_configuration_no or access_control_group_name or most_recent is required")
	}

	reqParams := server.GetAccessControlGroupListRequest{}
	if configNoOk {
		reqParams.AccessControlGroupConfigurationNoList = []*string{ncloud.String(configNo.(string))}
	}
	if acgNameOk {
		reqParams.AccessControlGroupName = ncloud.String(acgName.(string))
	}

	if isDefaultGroup, ok := d.GetOk("is_default_group"); ok {
		reqParams.IsDefault = ncloud.Bool(isDefaultGroup.(bool))
	}
	reqParams.PageNo = ncloud.Int32(1)

	logCommonRequest("GetAccessControlGroupList", reqParams)

	resp, err := client.server.V2Api.GetAccessControlGroupList(&reqParams)
	if err != nil {
		logErrorResponse("GetAccessControlGroupList", err, reqParams)
		return err
	}
	logCommonResponse("GetAccessControlGroupList", GetCommonResponse(resp))

	var accessControlGroup *server.AccessControlGroup
	var accessControlGroups []*server.AccessControlGroup

	for _, acg := range resp.AccessControlGroupList {
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

func accessControlGroupAttributes(d *schema.ResourceData, accessControlGroup *server.AccessControlGroup) error {
	d.SetId(*accessControlGroup.AccessControlGroupConfigurationNo)
	d.Set("access_control_group_configuration_no", accessControlGroup.AccessControlGroupConfigurationNo)
	d.Set("access_control_group_name", accessControlGroup.AccessControlGroupName)
	d.Set("access_control_group_description", accessControlGroup.AccessControlGroupDescription)
	d.Set("is_default_group", accessControlGroup.IsDefaultGroup)
	d.Set("create_date", accessControlGroup.CreateDate)

	return nil
}
