package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterResource("ncloud_network_acl_deny_allow_group", resourceNcloudNetworkACLDenyAllowGroup())
}

func resourceNcloudNetworkACLDenyAllowGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNetworkACLDenyAllowGroupCreate,
		Read:   resourceNcloudNetworkACLDenyAllowGroupRead,
		Update: resourceNcloudNetworkACLDenyAllowGroupUpdate,
		Delete: resourceNcloudNetworkACLDenyAllowGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_acl_deny_allow_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validateInstanceName),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 1000)),
			},
			"ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNcloudNetworkACLDenyAllowGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_network_acl_deny_allow_group`")
	}

	reqParams := &vpc.CreateNetworkAclDenyAllowGroupRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(d.Get("vpc_no").(string))}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NetworkAclDenyAllowGroupName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		reqParams.NetworkAclDenyAllowGroupDescription = ncloud.String(v.(string))
	}

	logCommonRequest("CreateNetworkAclDenyAllowGroup", reqParams)
	resp, err := config.Client.vpc.V2Api.CreateNetworkAclDenyAllowGroup(reqParams)
	if err != nil {
		logErrorResponse("CreateNetworkAclDenyAllowGroup", err, reqParams)
		return err
	}

	logResponse("CreateNetworkAclDenyAllowGroup", resp)

	instance := resp.NetworkAclDenyAllowGroupList[0]
	d.SetId(*instance.NetworkAclDenyAllowGroupNo)
	log.Printf("[INFO] Network ACL DenyAllowGroup ID: %s", d.Id())

	if err := waitForVpcNetworkAclDenyAllowGroupState(config, d.Id(), []string{InstanceStatusInit, InstanceStatusCreate}, []string{InstanceStatusRunning}, DefaultCreateTimeout); err != nil {
		return err
	}

	if err := setNetworkAclDenyAllowGroupIpList(d, config); err != nil {
		return err
	}

	if err := waitForVpcNetworkAclDenyAllowGroupState(config, d.Id(), []string{InstanceStatusSetting}, []string{InstanceStatusRunning}, DefaultCreateTimeout); err != nil {
		return err
	}

	return resourceNcloudNetworkACLDenyAllowGroupRead(d, meta)
}

func resourceNcloudNetworkACLDenyAllowGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getNetworkAclDenyAllowGroupDetail(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	m := map[string]interface{}{
		"id":                              *instance.NetworkAclDenyAllowGroupNo,
		"network_acl_deny_allow_group_no": *instance.NetworkAclDenyAllowGroupNo,
		"vpc_no":                          *instance.VpcNo,
		"name":                            *instance.NetworkAclDenyAllowGroupName,
		"description":                     *instance.NetworkAclDenyAllowGroupDescription,
		"ip_list":                         instance.IpList,
	}

	SetSingularResourceDataFromMapSchema(resourceNcloudNetworkACLDenyAllowGroup(), d, m)

	return nil
}

func resourceNcloudNetworkACLDenyAllowGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("ip_list") {
		if err := setNetworkAclDenyAllowGroupIpList(d, config); err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		if err := setNetworkAclDenyAllowGroupDescription(d, config); err != nil {
			return err
		}
	}

	if err := waitForVpcNetworkAclDenyAllowGroupState(config, d.Id(), []string{InstanceStatusSetting}, []string{InstanceStatusRunning}, DefaultTimeout); err != nil {
		return err
	}

	return resourceNcloudNetworkACLDenyAllowGroupRead(d, meta)
}

func resourceNcloudNetworkACLDenyAllowGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.DeleteNetworkAclDenyAllowGroupRequest{
		RegionCode:                 &config.RegionCode,
		NetworkAclDenyAllowGroupNo: ncloud.String(d.Id()),
	}

	logCommonRequest("DeleteNetworkAclDenyAllowGroup", reqParams)
	resp, err := config.Client.vpc.V2Api.DeleteNetworkAclDenyAllowGroup(reqParams)
	if err != nil {
		logErrorResponse("DeleteNetworkAclDenyAllowGroup", err, reqParams)
		return err
	}

	logResponse("DeleteNetworkAclDenyAllowGroup", resp)

	if err := waitForVpcNetworkAclDenyAllowGroupState(config, d.Id(), []string{InstanceStatusRunning, InstanceStatusTerminating}, []string{InstanceStatusTerminated}, DefaultTimeout); err != nil {
		return err
	}

	return nil
}

func waitForVpcNetworkAclDenyAllowGroupState(config *ProviderConfig, id string, pending []string, target []string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkAclDenyAllowGroupDetail(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclDenyAllowGroupStatus")
		},
		Timeout:    timeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for NetworkAclDenyAllowGroupStatus (%s) to become (%v): %s", id, target, err)
	}

	return nil
}

func getNetworkAclDenyAllowGroupDetail(config *ProviderConfig, id string) (*vpc.NetworkAclDenyAllowGroup, error) {
	reqParams := &vpc.GetNetworkAclDenyAllowGroupDetailRequest{
		RegionCode:                 &config.RegionCode,
		NetworkAclDenyAllowGroupNo: &id,
	}

	logCommonRequest("GetNetworkAclDenyAllowGroupDetail", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNetworkAclDenyAllowGroupDetail(reqParams)
	if err != nil {
		logErrorResponse("GetNetworkAclDenyAllowGroupDetail", err, reqParams)
		return nil, err
	}
	logResponse("GetNetworkAclDenyAllowGroupDetail", resp)

	if len(resp.NetworkAclDenyAllowGroupList) > 0 {
		instance := resp.NetworkAclDenyAllowGroupList[0]
		return instance, nil
	}

	return nil, nil
}

func setNetworkAclDenyAllowGroupDescription(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vpc.SetNetworkAclDenyAllowGroupDescriptionRequest{
		RegionCode:                          &config.RegionCode,
		NetworkAclDenyAllowGroupNo:          ncloud.String(d.Id()),
		NetworkAclDenyAllowGroupDescription: StringPtrOrNil(d.GetOk("description")),
	}

	logCommonRequest("SetNetworkAclDenyAllowGroupDescription", reqParams)
	resp, err := config.Client.vpc.V2Api.SetNetworkAclDenyAllowGroupDescription(reqParams)
	if err != nil {
		logErrorResponse("SetNetworkAclDenyAllowGroupDescription", err, reqParams)
		return err
	}
	logResponse("SetNetworkAclDenyAllowGroupDescription", resp)

	return nil
}

func setNetworkAclDenyAllowGroupIpList(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vpc.SetNetworkAclDenyAllowGroupIpListRequest{
		RegionCode:                 &config.RegionCode,
		NetworkAclDenyAllowGroupNo: ncloud.String(d.Id()),
		IpList:                     ncloud.StringInterfaceList(d.Get("ip_list").([]interface{})),
	}

	logCommonRequest("SetNetworkAclDenyAllowGroupIpList", reqParams)
	resp, err := config.Client.vpc.V2Api.SetNetworkAclDenyAllowGroupIpList(reqParams)
	if err != nil {
		logErrorResponse("SetNetworkAclDenyAllowGroupIpList", err, reqParams)
		return err
	}
	logResponse("SetNetworkAclDenyAllowGroupIpList", resp)

	return nil
}
