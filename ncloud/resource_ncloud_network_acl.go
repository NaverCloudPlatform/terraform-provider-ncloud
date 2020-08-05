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
				Description:  "Network ACL name to create. default: Assigned by NAVER CLOUD PLATFORM.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
				Description:  "Description of a Network ACL to create.",
			},
			"vpc_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the VPC that the desired acl belongs to.",
			},
			"network_acl_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	regionCode := meta.(*ProviderConfig).RegionCode

	reqParams := &vpc.CreateNetworkAclRequest{
		RegionCode: &regionCode,
		VpcNo:      ncloud.String(d.Get("vpc_no").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NetworkAclName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		reqParams.NetworkAclDescription = ncloud.String(v.(string))
	}

	logCommonRequest("resource_ncloud_network_acl > CreateNetworkAcl", reqParams)
	resp, err := client.vpc.V2Api.CreateNetworkAcl(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_network_acl > CreateNetworkAcl", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_network_acl > CreateNetworkAcl", resp)

	instance := resp.NetworkAclList[0]
	d.SetId(*instance.NetworkAclNo)
	log.Printf("[INFO] Network ACL ID: %s", d.Id())

	waitForNcloudNetworkACLCreation(client, d.Id())

	return resourceNcloudNetworkACLRead(d, meta)
}

func resourceNcloudNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	instance, err := getNetworkACLInstance(client, d.Id())
	if err != nil {
		d.SetId("")
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
	d.Set("status", instance.NetworkAclStatus.Code)

	return nil
}

func resourceNcloudNetworkACLUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudNetworkACLRead(d, meta)
}

func resourceNcloudNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	regionCode := meta.(*ProviderConfig).RegionCode

	reqParams := &vpc.DeleteNetworkAclRequest{
		NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
		RegionCode:   &regionCode,
	}

	logCommonRequest("resource_ncloud_network_acl > DeleteNetworkAcl", reqParams)
	resp, err := client.vpc.V2Api.DeleteNetworkAcl(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_network_acl > DeleteNetworkAcl", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_network_acl > DeleteNetworkAcl", resp)

	waitForNcloudNetworkACLDeletion(client, d.Id())

	return nil
}

func waitForNcloudNetworkACLCreation(client *NcloudAPIClient, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkACLInstance(client, id)
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

func waitForNcloudNetworkACLDeletion(client *NcloudAPIClient, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkACLInstance(client, id)
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

func getNetworkACLInstance(client *NcloudAPIClient, id string) (*vpc.NetworkAcl, error) {
	reqParams := &vpc.GetNetworkAclDetailRequest{
		NetworkAclNo: ncloud.String(id),
	}

	logCommonRequest("resource_ncloud_network_acl > GetNetworkAclDetail", reqParams)
	resp, err := client.vpc.V2Api.GetNetworkAclDetail(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_network_acl > GetNetworkAclDetail", err, reqParams)
		return nil, err
	}
	logResponse("resource_ncloud_network_acl > GetNetworkAclDetail", resp)

	if len(resp.NetworkAclList) > 0 {
		instance := resp.NetworkAclList[0]
		return instance, nil
	}

	return nil, nil
}
