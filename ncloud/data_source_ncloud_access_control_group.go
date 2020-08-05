package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudAccessControlGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAccessControlGroupRead,

		Schema: map[string]*schema.Schema{
			"configuration_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNcloudAccessControlGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	configNo, configNoOk := d.GetOk("configuration_no")
	acgName, acgNameOk := d.GetOk("name")
	isDefaultGroup, isDefaultGroupOk := d.GetOk("is_default_group")

	if !configNoOk && !acgNameOk && !isDefaultGroupOk {
		return fmt.Errorf("either configuration_no or name or is_default_group is required")
	}

	reqParams := server.GetAccessControlGroupListRequest{}
	if configNoOk {
		reqParams.AccessControlGroupConfigurationNoList = []*string{ncloud.String(configNo.(string))}
	}
	if acgNameOk {
		reqParams.AccessControlGroupName = ncloud.String(acgName.(string))
	}
	if isDefaultGroupOk {
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

	if err := validateOneResult(len(accessControlGroups)); err != nil {
		return err
	}
	accessControlGroup = accessControlGroups[0]
	return accessControlGroupAttributes(d, accessControlGroup)
}

func accessControlGroupAttributes(d *schema.ResourceData, accessControlGroup *server.AccessControlGroup) error {
	d.SetId(*accessControlGroup.AccessControlGroupConfigurationNo)
	d.Set("configuration_no", accessControlGroup.AccessControlGroupConfigurationNo)
	d.Set("name", accessControlGroup.AccessControlGroupName)
	d.Set("description", accessControlGroup.AccessControlGroupDescription)
	d.Set("is_default_group", accessControlGroup.IsDefaultGroup)

	return nil
}
