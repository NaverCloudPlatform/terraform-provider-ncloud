package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudAccessControlGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAccessControlGroupsRead,

		Schema: map[string]*schema.Schema{
			"configuration_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Description: "List of ACG configuration numbers you want to get",
			},
			"is_default_group": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to get default groups only",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the ACG you want to get",
			},
			"page_no": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number based on the page size if the number of items is large.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of items to be shown per page",
			},

			"access_control_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A List of access control group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACG configuration number",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACG name",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACG description",
						},
						"is_default_group": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether default group",
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
	client := meta.(*NcloudAPIClient)

	d.SetId(time.Now().UTC().String())

	reqParams := &server.GetAccessControlGroupListRequest{}
	var paramAccessControlGroupConfigurationNoList []*string
	if param, ok := d.GetOk("configuration_no_list"); ok {
		paramAccessControlGroupConfigurationNoList = expandStringInterfaceList(param.([]interface{}))
	}
	reqParams.AccessControlGroupConfigurationNoList = paramAccessControlGroupConfigurationNoList
	reqParams.AccessControlGroupName = ncloud.String(d.Get("name").(string))
	if isDefaultGroup, ok := d.GetOk("is_default_group"); ok {
		reqParams.IsDefault = ncloud.Bool(isDefaultGroup.(bool))
	}
	if pageNo, ok := d.GetOk("page_no"); ok {
		reqParams.PageNo = ncloud.Int32(int32(pageNo.(int)))
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		reqParams.PageSize = ncloud.Int32(int32(pageSize.(int)))
	}

	resp, err := getAccessControlGroupList(client, reqParams)
	if err != nil {
		return err
	}
	var accessControlGroups []*server.AccessControlGroup

	for _, group := range resp.AccessControlGroupList {
		accessControlGroups = append(accessControlGroups, group)
	}

	if len(accessControlGroups) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return accessControlGroupsAttributes(d, accessControlGroups)
}

func getAccessControlGroupList(client *NcloudAPIClient, reqParams *server.GetAccessControlGroupListRequest) (*server.GetAccessControlGroupListResponse, error) {
	logCommonRequest("GetAccessControlGroupList", reqParams)
	resp, err := client.server.V2Api.GetAccessControlGroupList(reqParams)
	if err != nil {
		logErrorResponse("GetAccessControlGroupList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetAccessControlGroupList", GetCommonResponse(resp))
	return resp, nil
}

func accessControlGroupsAttributes(d *schema.ResourceData, accessControlGroups []*server.AccessControlGroup) error {
	var ids []string

	for _, accessControlGroup := range accessControlGroups {
		ids = append(ids, ncloud.StringValue(accessControlGroup.AccessControlGroupConfigurationNo))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("access_control_groups", flattenAccessControlGroups(accessControlGroups)); err != nil {
		return err
	}

	// create a json file in current directory and write d source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), d.Get("access_control_groups"))
	}

	return nil
}
