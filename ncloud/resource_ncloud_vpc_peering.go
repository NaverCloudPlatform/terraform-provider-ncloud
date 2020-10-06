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
	RegisterResource("ncloud_vpc_peering", resourceNcloudVpcPeering())
}

func resourceNcloudVpcPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudVpcPeeringCreate,
		Read:   resourceNcloudVpcPeeringRead,
		Update: resourceNcloudVpcPeeringUpdate,
		Delete: resourceNcloudVpcPeeringDelete,
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
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"source_vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_vpc_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"target_vpc_login_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_peering_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_reverse_vpc_peering": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_between_accounts": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNcloudVpcPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_vpc_peering`")
	}

	reqParams := &vpc.CreateVpcPeeringInstanceRequest{
		RegionCode:  &config.RegionCode,
		SourceVpcNo: ncloud.String(d.Get("source_vpc_no").(string)),
		TargetVpcNo: ncloud.String(d.Get("target_vpc_no").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcPeeringName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		reqParams.VpcPeeringDescription = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("target_vpc_name"); ok {
		reqParams.TargetVpcName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("target_vpc_login_id"); ok {
		reqParams.TargetVpcLoginId = ncloud.String(v.(string))
	}

	logCommonRequest("CreateVpcPeeringInstance", reqParams)
	resp, err := config.Client.vpc.V2Api.CreateVpcPeeringInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateVpcPeeringInstance", err, reqParams)
		return err
	}

	logResponse("CreateVpcPeeringInstance", resp)

	instance := resp.VpcPeeringInstanceList[0]
	d.SetId(*instance.VpcPeeringInstanceNo)
	log.Printf("[INFO] VPC Peering ID: %s", d.Id())

	if err := waitForNcloudVpcPeeringCreation(config, d.Id()); err != nil {
		return err
	}

	return resourceNcloudVpcPeeringRead(d, meta)
}

func resourceNcloudVpcPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instance, err := getVpcPeeringInstance(config, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.VpcPeeringInstanceNo)
	d.Set("vpc_peering_no", instance.VpcPeeringInstanceNo)
	d.Set("name", instance.VpcPeeringName)
	d.Set("description", instance.VpcPeeringDescription)
	d.Set("source_vpc_no", instance.SourceVpcNo)
	d.Set("target_vpc_no", instance.TargetVpcNo)
	d.Set("target_vpc_name", instance.TargetVpcName)
	d.Set("target_vpc_login_id", instance.TargetVpcLoginId)
	d.Set("status", instance.VpcPeeringInstanceStatus.Code)
	d.Set("has_reverse_vpc_peering", instance.HasReverseVpcPeering)
	d.Set("is_between_accounts", instance.IsBetweenAccounts)

	return nil
}

func resourceNcloudVpcPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("description") {
		if err := setVpcPeeringDescription(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudVpcPeeringRead(d, meta)
}

func resourceNcloudVpcPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.DeleteVpcPeeringInstanceRequest{
		RegionCode:           &config.RegionCode,
		VpcPeeringInstanceNo: ncloud.String(d.Get("vpc_peering_no").(string)),
	}

	logCommonRequest("DeleteVpcPeeringInstance", reqParams)
	resp, err := config.Client.vpc.V2Api.DeleteVpcPeeringInstance(reqParams)
	if err != nil {
		logErrorResponse("DeleteVpcPeeringInstance", err, reqParams)
		return err
	}

	logResponse("DeleteVpcPeeringInstance", resp)

	if err := waitForNcloudVpcPeeringDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForNcloudVpcPeeringCreation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getVpcPeeringInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcPeeringInstanceStatus")
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Peering (%s) to become available: %s", id, err)
	}

	return nil
}

func waitForNcloudVpcPeeringDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getVpcPeeringInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcPeeringInstanceStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Peering (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func getVpcPeeringInstance(config *ProviderConfig, id string) (*vpc.VpcPeeringInstance, error) {
	reqParams := &vpc.GetVpcPeeringInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		VpcPeeringInstanceNo: ncloud.String(id),
	}

	logCommonRequest("GetVpcPeeringInstanceDetail", reqParams)
	resp, err := config.Client.vpc.V2Api.GetVpcPeeringInstanceDetail(reqParams)
	if err != nil {
		logErrorResponse("GetVpcPeeringInstanceDetail", err, reqParams)
		return nil, err
	}
	logResponse("GetVpcPeeringInstanceDetail", resp)

	if len(resp.VpcPeeringInstanceList) > 0 {
		instance := resp.VpcPeeringInstanceList[0]
		return instance, nil
	}

	return nil, nil
}

func setVpcPeeringDescription(d *schema.ResourceData, config *ProviderConfig) error {
	reqParams := &vpc.SetVpcPeeringDescriptionRequest{
		RegionCode:            &config.RegionCode,
		VpcPeeringInstanceNo:  ncloud.String(d.Id()),
		VpcPeeringDescription: StringPtrOrNil(d.GetOk("description")),
	}

	logCommonRequest("setVpcPeeringDescription", reqParams)
	resp, err := config.Client.vpc.V2Api.SetVpcPeeringDescription(reqParams)
	if err != nil {
		logErrorResponse("setVpcPeeringDescription", err, reqParams)
		return err
	}
	logResponse("setVpcPeeringDescription", resp)

	return nil
}
