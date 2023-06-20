package server

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func init() {
	RegisterResource("ncloud_access_control_group", resourceNcloudAccessControlGroup())
}

func resourceNcloudAccessControlGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudAccessControlGroupCreate,
		Read:   resourceNcloudAccessControlGroupRead,
		Delete: resourceNcloudAccessControlGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(ValidateInstanceName),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(0, 1000)),
			},
			"access_control_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNcloudAccessControlGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := createAccessControlGroup(d, config)

	if err != nil {
		return err
	}

	d.SetId(*instance.AccessControlGroupNo)
	log.Printf("[INFO] ACG ID: %s", d.Id())

	return resourceNcloudAccessControlGroupRead(d, meta)
}

func resourceNcloudAccessControlGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getAccessControlGroup(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.AccessControlGroupNo)
	d.Set("access_control_group_no", instance.AccessControlGroupNo)
	d.Set("name", instance.AccessControlGroupName)
	d.Set("description", instance.AccessControlGroupDescription)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("is_default", instance.IsDefault)

	return nil
}

func resourceNcloudAccessControlGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if err := deleteAccessControlGroup(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func getAccessControlGroup(config *ProviderConfig, id string) (*vserver.AccessControlGroup, error) {
	if config.SupportVPC {
		return getVpcAccessControlGroup(config, id)
	}

	return nil, NotSupportClassic("resource `ncloud_access_control_group`")
}

func getVpcAccessControlGroup(config *ProviderConfig, id string) (*vserver.AccessControlGroup, error) {
	reqParams := &vserver.GetAccessControlGroupDetailRequest{
		RegionCode:           &config.RegionCode,
		AccessControlGroupNo: ncloud.String(id),
	}

	LogCommonRequest("getVpcAccessControlGroup", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetAccessControlGroupDetail(reqParams)
	if err != nil {
		LogErrorResponse("getVpcAccessControlGroup", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcAccessControlGroup", resp)

	if len(resp.AccessControlGroupList) > 0 {
		return resp.AccessControlGroupList[0], nil
	}

	return nil, nil
}

func createAccessControlGroup(d *schema.ResourceData, config *ProviderConfig) (*vserver.AccessControlGroup, error) {
	if config.SupportVPC {
		return createVpcAccessControlGroup(d, config)
	}

	return nil, NotSupportClassic("resource `ncloud_access_control_group`")
}

func createVpcAccessControlGroup(d *schema.ResourceData, config *ProviderConfig) (*vserver.AccessControlGroup, error) {
	reqParams := &vserver.CreateAccessControlGroupRequest{
		RegionCode:                    &config.RegionCode,
		VpcNo:                         ncloud.String(d.Get("vpc_no").(string)),
		AccessControlGroupName:        StringPtrOrNil(d.GetOk("name")),
		AccessControlGroupDescription: StringPtrOrNil(d.GetOk("description")),
	}

	LogCommonRequest("createVpcAccessControlGroup", reqParams)
	resp, err := config.Client.Vserver.V2Api.CreateAccessControlGroup(reqParams)
	if err != nil {
		LogErrorResponse("createVpcAccessControlGroup", err, reqParams)
		return nil, err
	}
	LogResponse("createVpcAccessControlGroup", resp)

	return resp.AccessControlGroupList[0], nil
}

func deleteAccessControlGroup(config *ProviderConfig, id string) error {
	if config.SupportVPC {
		return deleteVpcAccessControlGroup(config, id)
	}

	return NotSupportClassic("resource `ncloud_access_control_group`")
}

func deleteVpcAccessControlGroup(config *ProviderConfig, id string) error {
	accessControlGroup, err := getAccessControlGroup(config, id)
	if err != nil {
		return err
	}

	if accessControlGroup == nil {
		return fmt.Errorf("no matching Access Control Group: %s", id)
	}

	reqParams := &vserver.DeleteAccessControlGroupRequest{
		RegionCode:           &config.RegionCode,
		VpcNo:                accessControlGroup.VpcNo,
		AccessControlGroupNo: ncloud.String(id),
	}

	LogCommonRequest("deleteVpcAccessControlGroup", reqParams)
	resp, err := config.Client.Vserver.V2Api.DeleteAccessControlGroup(reqParams)
	if err != nil {
		LogErrorResponse("deleteVpcAccessControlGroup", err, reqParams)
		return err
	}
	LogResponse("deleteVpcAccessControlGroup", resp)

	if err := waitForVpcAccessControlGroupDeletion(config, id); err != nil {
		return err
	}

	return nil
}

func waitForVpcAccessControlGroupDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getAccessControlGroup(config, id)
			return vpc.VpcCommonStateRefreshFunc(instance, err, "AccessControlGroupStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Access Control Group (%s) to become terminated: %s", id, err)
	}

	return nil
}

func waitForVpcAccessControlGroupRunning(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getAccessControlGroup(config, id)
			return vpc.VpcCommonStateRefreshFunc(instance, err, "AccessControlGroupStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Access Control Group (%s) to become running: %s", id, err)
	}

	return nil
}
