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

func resourceNcloudVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudVpcCreate,
		Read:   resourceNcloudVpcRead,
		Update: resourceNcloudVpcUpdate,
		Delete: resourceNcloudVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: ncloudVpcCommonCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
				Description:  "Subnet name to create. default: Assigned by NAVER CLOUD PLATFORM.",
			},
			"ipv4_cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(16, 28),
				Description:  "The CIDR block for the vpc.",
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_network_acl_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.CreateVpcRequest{
		Ipv4CidrBlock: ncloud.String(d.Get("ipv4_cidr_block").(string)),
		RegionCode:    regionCode,
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcName = ncloud.String(v.(string))
	}

	logCommonRequest("resource_ncloud_vpc > CreateVpc", reqParams)
	resp, err := client.vpc.V2Api.CreateVpc(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_vpc > Create Vpc Instance", err, reqParams)
		return err
	}

	logCommonResponse("resource_ncloud_vpc > CreateVpc", GetCommonResponse(resp))

	vpcInstance := resp.VpcList[0]
	d.SetId(*vpcInstance.VpcNo)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INIT", "CREATING"},
		Target:     []string{"RUN"},
		Refresh:    VPCStateRefreshFunc(client, d.Id()),
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for VPC (%s) to become available: %s",
			d.Id(), err)
	}

	log.Printf("[INFO] VPC ID: %s", d.Id())

	return resourceNcloudVpcRead(d, meta)
}

func resourceNcloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	instance, err := getVpcInstance(client, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.VpcNo)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("name", instance.VpcName)
	d.Set("ipv4_cidr_block", instance.Ipv4CidrBlock)
	d.Set("status", instance.VpcStatus.Code)

	defaultNetworkACLNo, err := getDefaultNetworkACL(client, d.Id())
	if err != nil {
		return fmt.Errorf("Error get default network acl for VPC (%s): %s", d.Id(), err)
	}

	d.Set("default_network_acl_no", defaultNetworkACLNo)

	return nil
}

func getDefaultNetworkACL(client *NcloudAPIClient, id string) (string, error) {
	reqParams := &vpc.GetNetworkAclListRequest{
		VpcNo: ncloud.String(id),
	}

	logCommonRequest("resource_ncloud_vpc > GetNetworkAclList", reqParams)
	resp, err := client.vpc.V2Api.GetNetworkAclList(reqParams)

	if err != nil {
		logErrorResponse("resource_ncloud_vpc > GetNetworkAclList", err, reqParams)
		return "", err
	}

	logResponse("resource_ncloud_vpc > GetNetworkAclList", resp)

	if resp == nil || len(resp.NetworkAclList) == 0 {
		return "", fmt.Errorf("no matching Network ACL found")
	}

	for _, i := range resp.NetworkAclList {
		if *i.IsDefault {
			return *i.NetworkAclNo, nil
		}
	}

	return "", fmt.Errorf("No matching default network ACL found")
}

func resourceNcloudVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudVpcRead(d, meta)
}

func resourceNcloudVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.DeleteVpcRequest{
		VpcNo:      ncloud.String(d.Get("vpc_no").(string)),
		RegionCode: regionCode,
	}

	logCommonRequest("resource_ncloud_vpc > DeleteVpc", reqParams)
	resp, err := client.vpc.V2Api.DeleteVpc(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_vpc > DeleteVpc Vpc Instance", err, reqParams)
		return err
	}
	logResponse("resource_ncloud_vpc > DeleteVpc", resp)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"RUN", "TERMTING"},
		Target:     []string{"TERMINATED"},
		Refresh:    VPCStateRefreshFunc(client, d.Id()),
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for VPC (%s) to become termintaing: %s",
			d.Id(), err)
	}

	return nil
}

func getVpcInstance(client *NcloudAPIClient, id string) (*vpc.Vpc, error) {
	reqParams := &vpc.GetVpcDetailRequest{
		VpcNo: ncloud.String(id),
	}

	resp, err := client.vpc.V2Api.GetVpcDetail(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_vpc > Get Vpc Instance", err, reqParams)
		return nil, err
	}
	logResponse("resource_ncloud_vpc > GetVpcDetail", resp)

	if len(resp.VpcList) > 0 {
		vpc := resp.VpcList[0]
		return vpc, nil
	}

	return nil, nil
}

// VPCStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch a VPC
func VPCStateRefreshFunc(client *NcloudAPIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vpc, err := getVpcInstance(client, id)

		if err != nil {
			return nil, "", err
		}

		if vpc == nil {
			return vpc, "TERMINATED", nil
		}

		return vpc, *vpc.VpcStatus.Code, nil
	}
}
