package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterResource("ncloud_network_acl", resourceNcloudNetworkACL())
}

func resourceNcloudNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNetworkACLCreate,
		Read:   resourceNcloudNetworkACLRead,
		Update: resourceNcloudNetworkACLUpdate,
		Delete: resourceNcloudNetworkACLDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: ncloudVpcCommonCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_acl_no": {
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

func resourceNcloudNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_network_acl`")
	}

	reqParams := &vpc.CreateNetworkAclRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(d.Get("vpc_no").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NetworkAclName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		reqParams.NetworkAclDescription = ncloud.String(v.(string))
	}

	logCommonRequest("CreateNetworkAcl", reqParams)
	resp, err := config.Client.vpc.V2Api.CreateNetworkAcl(reqParams)
	if err != nil {
		logErrorResponse("CreateNetworkAcl", err, reqParams)
		return err
	}

	logResponse("CreateNetworkAcl", resp)

	instance := resp.NetworkAclList[0]
	d.SetId(*instance.NetworkAclNo)
	log.Printf("[INFO] Network ACL ID: %s", d.Id())

	if err := waitForNcloudNetworkACLCreation(config, d.Id()); err != nil {
		return err
	}

	return resourceNcloudNetworkACLRead(d, meta)
}

func resourceNcloudNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getNetworkACLInstance(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.NetworkAclNo)
	d.Set("network_acl_no", instance.NetworkAclNo)
	d.Set("name", instance.NetworkAclName)
	d.Set("description", instance.NetworkAclDescription)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("is_default", instance.IsDefault)

	return nil
}

func resourceNcloudNetworkACLUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("description") {
		if err := setNetworkACLDescription(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudNetworkACLRead(d, meta)
}

func resourceNcloudNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.DeleteNetworkAclRequest{
		RegionCode:   &config.RegionCode,
		NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
	}

	logCommonRequest("DeleteNetworkAcl", reqParams)
	resp, err := config.Client.vpc.V2Api.DeleteNetworkAcl(reqParams)
	if err != nil {
		logErrorResponse("DeleteNetworkAcl", err, reqParams)
		return err
	}

	logResponse("DeleteNetworkAcl", resp)

	if err := waitForNcloudNetworkACLDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForNcloudNetworkACLCreation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkACLInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclStatus")
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Network ACL (%s) to become available: %s", id, err)
	}

	return nil
}

func waitForNcloudNetworkACLDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkACLInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Network ACL (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func getNetworkACLInstance(config *ProviderConfig, id string) (*vpc.NetworkAcl, error) {
	reqParams := &vpc.GetNetworkAclDetailRequest{
		RegionCode:   &config.RegionCode,
		NetworkAclNo: ncloud.String(id),
	}

	logCommonRequest("GetNetworkAclDetail", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNetworkAclDetail(reqParams)
	if err != nil {
		logErrorResponse("GetNetworkAclDetail", err, reqParams)
		return nil, err
	}
	logResponse("GetNetworkAclDetail", resp)

	if len(resp.NetworkAclList) > 0 {
		instance := resp.NetworkAclList[0]
		return instance, nil
	}

	return nil, nil
}

func setNetworkACLDescription(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vpc.SetNetworkAclDescriptionRequest{
		RegionCode:            &config.RegionCode,
		NetworkAclNo:          ncloud.String(d.Id()),
		NetworkAclDescription: StringPtrOrNil(d.GetOk("description")),
	}

	logCommonRequest("setNetworkAclDescription", reqParams)
	resp, err := config.Client.vpc.V2Api.SetNetworkAclDescription(reqParams)
	if err != nil {
		logErrorResponse("setNetworkAclDescription", err, reqParams)
		return err
	}
	logResponse("setNetworkAclDescription", resp)

	return nil
}
