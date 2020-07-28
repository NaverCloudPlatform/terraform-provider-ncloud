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
		CustomizeDiff: resourceNcloudNetworkACLCustomizeDiff,
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
	client := meta.(*NcloudAPIClient)
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.CreateNetworkAclRequest{
		RegionCode: regionCode,
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

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INIT", "CREATING"},
		Target:     []string{"RUN"},
		Refresh:    NetworkACLStateRefreshFunc(client, d.Id()),
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for Network ACL (%s) to become available: %s",
			d.Id(), err)
	}

	log.Printf("[INFO] Network ACL ID: %s", d.Id())

	return resourceNcloudNetworkACLRead(d, meta)
}

func resourceNcloudNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

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
	client := meta.(*NcloudAPIClient)

	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.DeleteNetworkAclRequest{
		NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
		RegionCode:   regionCode,
	}

	logCommonRequest("resource_ncloud_network_acl > DeleteNetworkAcl", reqParams)
	resp, err := client.vpc.V2Api.DeleteNetworkAcl(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_network_acl > DeleteNetworkAcl", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_network_acl > DeleteNetworkAcl", resp)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"RUN", "TERMTING"},
		Target:     []string{"TERMINATED"},
		Refresh:    NetworkACLStateRefreshFunc(client, d.Id()),
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for Network ACL (%s) to become termintaing: %s",
			d.Id(), err)
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

// NetworkACLStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch a Network ACL
func NetworkACLStateRefreshFunc(client *NcloudAPIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getNetworkACLInstance(client, id)

		if err != nil {
			return nil, "", err
		}

		if instance == nil {
			return instance, "TERMINATED", nil
		}

		return instance, *instance.NetworkAclStatus.Code, nil
	}
}

func resourceNcloudNetworkACLCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChange("name") {
		old, new := diff.GetChange("name")
		if len(old.(string)) > 0 {
			return fmt.Errorf("Change 'name' is not support, Please set name as a old value = [%s -> %s]", new, old)
		}
	}

	return nil
}
