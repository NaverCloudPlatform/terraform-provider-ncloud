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

func resourceNcloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNatGatewayCreate,
		Read:   resourceNcloudNatGatewayRead,
		Update: resourceNcloudNatGatewayUpdate,
		Delete: resourceNcloudNatGatewayDelete,
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
				Description:  "NAT Gateway name to create. default: Assigned by NAVER CLOUD PLATFORM.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
				Description:  "Description of a NAT Gateway to create.",
			},
			"vpc_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the VPC that the desired nat gateway belongs to.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Available Zone. Get available values using the `data ncloud_zones`.",
			},
			"nat_gateway_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	regionCode := meta.(*ProviderConfig).RegionCode

	reqParams := &vpc.CreateNatGatewayInstanceRequest{
		RegionCode: &regionCode,
		VpcNo:      ncloud.String(d.Get("vpc_no").(string)),
		ZoneCode:   ncloud.String(d.Get("zone").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NatGatewayName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		reqParams.NatGatewayDescription = ncloud.String(v.(string))
	}

	logCommonRequest("resource_ncloud_nat_gateway > CreateNatGatewayInstance", reqParams)
	resp, err := client.vpc.V2Api.CreateNatGatewayInstance(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_nat_gateway > CreateNatGatewayInstance", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_nat_gateway > CreateNatGatewayInstance", resp)

	instance := resp.NatGatewayInstanceList[0]
	d.SetId(*instance.NatGatewayInstanceNo)
	log.Printf("[INFO] NAT Gateway ID: %s", d.Id())

	waitForNcloudNatGatewayCreation(client, d.Id())

	return resourceNcloudNatGatewayRead(d, meta)
}

func resourceNcloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	instance, err := getNatGatewayInstance(client, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.NatGatewayInstanceNo)
	d.Set("nat_gateway_no", instance.NatGatewayInstanceNo)
	d.Set("name", instance.NatGatewayName)
	d.Set("description", instance.NatGatewayDescription)
	d.Set("public_ip", instance.PublicIp)
	d.Set("status", instance.NatGatewayInstanceStatus.Code)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("zone", instance.ZoneCode)

	return nil
}

func resourceNcloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudNatGatewayRead(d, meta)
}

func resourceNcloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	regionCode := meta.(*ProviderConfig).RegionCode

	reqParams := &vpc.DeleteNatGatewayInstanceRequest{
		NatGatewayInstanceNo: ncloud.String(d.Get("nat_gateway_no").(string)),
		RegionCode:           &regionCode,
	}

	logCommonRequest("resource_ncloud_nat_gateway > DeleteNatGatewayInstance", reqParams)
	resp, err := client.vpc.V2Api.DeleteNatGatewayInstance(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_nat_gateway > DeleteNatGatewayInstance", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_nat_gateway > DeleteNatGatewayInstance", resp)

	waitForNcloudNatGatewayDeletion(client, d.Id())

	return nil
}

func waitForNcloudNatGatewayCreation(client *NcloudAPIClient, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNatGatewayInstance(client, id)
			return VpcCommonStateRefreshFunc(instance, err, "NatGatewayInstanceStatus")
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become available: %s", id, err)
	}

	return nil
}

func waitForNcloudNatGatewayDeletion(client *NcloudAPIClient, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNatGatewayInstance(client, id)
			return VpcCommonStateRefreshFunc(instance, err, "NatGatewayInstanceStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func getNatGatewayInstance(client *NcloudAPIClient, id string) (*vpc.NatGatewayInstance, error) {
	reqParams := &vpc.GetNatGatewayInstanceDetailRequest{
		NatGatewayInstanceNo: ncloud.String(id),
	}

	logCommonRequest("resource_ncloud_nat_gateway > GetNatGatewayInstanceDetail", reqParams)
	resp, err := client.vpc.V2Api.GetNatGatewayInstanceDetail(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_nat_gateway > GetNatGatewayInstanceDetail", err, reqParams)
		return nil, err
	}
	logResponse("resource_ncloud_nat_gateway > GetNatGatewayInstanceDetail", resp)

	if len(resp.NatGatewayInstanceList) > 0 {
		instance := resp.NatGatewayInstanceList[0]
		return instance, nil
	}

	return nil, nil
}
