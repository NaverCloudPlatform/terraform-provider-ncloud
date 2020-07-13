package ncloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNcloudVpc() *schema.Resource {
	return &schema.Resource{
		Create:        resourceNcloudVpcCreate,
		Read:          resourceNcloudVpcRead,
		Update:        resourceNcloudVpcUpdate,
		Delete:        resourceNcloudVpcDelete,
		SchemaVersion: 1,
		CustomizeDiff: resourceNcloudVpcCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 30),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-*]+$`), "Composed of alphabets, numbers, hyphen (-) and wild card (*)."),
					validation.StringMatch(regexp.MustCompile(`.*[^\\-]$`), "Hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input."),
				),
			},
			"ipv4_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
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
		VpcName:       ncloud.String(d.Get("name").(string)),
		Ipv4CidrBlock: ncloud.String(d.Get("ipv4_cidr_block").(string)),
		RegionCode:    regionCode,
	}

	logCommonRequest("CreateVpc", reqParams)
	resp, err := client.vpc.V2Api.CreateVpc(reqParams)
	if err != nil {
		logErrorResponse("Create Vpc Instance", err, reqParams)
		return err
	}

	logCommonResponse("CreateVpc", GetCommonResponse(resp))

	vpcInstance := resp.VpcList[0]
	d.SetId(*vpcInstance.VpcNo)
	d.Set("vpc_no", vpcInstance.VpcNo)

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

	vpc, err := getVpcInstance(client, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if vpc == nil {
		d.SetId("")
		return nil
	}

	return vpcInstanceAttributes(d, vpc)
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

	logCommonRequest("DeleteVpc", reqParams)
	resp, err := client.vpc.V2Api.DeleteVpc(reqParams)
	if err != nil {
		logErrorResponse("DeleteVpc Vpc Instance", err, reqParams)
		return err
	}

	logCommonResponse("DeleteVpc", GetCommonResponse(resp))

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
		logErrorResponse("Get Vpc Instance", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetVpcDetail", GetCommonResponse(resp))

	if len(resp.VpcList) > 0 {
		vpc := resp.VpcList[0]
		return vpc, nil
	}

	return nil, nil
}

// VPCStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a VPC.
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

func resourceNcloudVpcCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChange("name") {
		old, new := diff.GetChange("name")
		return fmt.Errorf("Change 'name' is not support, Please set name as a old value = [%s -> %s]", new, old)
	}

	return nil
}
