package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func ResourceNcloudVpcPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudVpcPeeringCreate,
		Read:   resourceNcloudVpcPeeringRead,
		Update: resourceNcloudVpcPeeringUpdate,
		Delete: resourceNcloudVpcPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(verify.ValidateInstanceName),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringLenBetween(0, 1000)),
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
				ForceNew: true,
			},
			"target_vpc_login_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_peering_no": {
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
	config := meta.(*conn.ProviderConfig)

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

	LogCommonRequest("CreateVpcPeeringInstance", reqParams)
	resp, err := config.Client.Vpc.V2Api.CreateVpcPeeringInstance(reqParams)
	if err != nil {
		LogErrorResponse("CreateVpcPeeringInstance", err, reqParams)
		return err
	}

	LogResponse("CreateVpcPeeringInstance", resp)

	instance := resp.VpcPeeringInstanceList[0]
	d.SetId(*instance.VpcPeeringInstanceNo)
	log.Printf("[INFO] VPC Peering ID: %s", d.Id())

	if err := waitForNcloudVpcPeeringCreation(config, d.Id()); err != nil {
		return err
	}

	return resourceNcloudVpcPeeringRead(d, meta)
}

func resourceNcloudVpcPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instance, err := GetVpcPeeringInstance(config, d.Id())
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
	d.Set("has_reverse_vpc_peering", instance.HasReverseVpcPeering)
	d.Set("is_between_accounts", instance.IsBetweenAccounts)

	return nil
}

func resourceNcloudVpcPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if d.HasChange("description") {
		if err := setVpcPeeringDescription(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudVpcPeeringRead(d, meta)
}

func resourceNcloudVpcPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	reqParams := &vpc.DeleteVpcPeeringInstanceRequest{
		RegionCode:           &config.RegionCode,
		VpcPeeringInstanceNo: ncloud.String(d.Get("vpc_peering_no").(string)),
	}

	LogCommonRequest("DeleteVpcPeeringInstance", reqParams)
	resp, err := config.Client.Vpc.V2Api.DeleteVpcPeeringInstance(reqParams)
	if err != nil {
		LogErrorResponse("DeleteVpcPeeringInstance", err, reqParams)
		return err
	}

	LogResponse("DeleteVpcPeeringInstance", resp)

	if err := WaitForNcloudVpcPeeringDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForNcloudVpcPeeringCreation(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetVpcPeeringInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcPeeringInstanceStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Peering (%s) to become available: %s", id, err)
	}

	return nil
}

func WaitForNcloudVpcPeeringDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetVpcPeeringInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcPeeringInstanceStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Peering (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetVpcPeeringInstance(config *conn.ProviderConfig, id string) (*vpc.VpcPeeringInstance, error) {
	reqParams := &vpc.GetVpcPeeringInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		VpcPeeringInstanceNo: ncloud.String(id),
	}

	LogCommonRequest("GetVpcPeeringInstanceDetail", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetVpcPeeringInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("GetVpcPeeringInstanceDetail", err, reqParams)
		return nil, err
	}
	LogResponse("GetVpcPeeringInstanceDetail", resp)

	if len(resp.VpcPeeringInstanceList) > 0 {
		instance := resp.VpcPeeringInstanceList[0]
		return instance, nil
	}

	return nil, nil
}

func setVpcPeeringDescription(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vpc.SetVpcPeeringDescriptionRequest{
		RegionCode:            &config.RegionCode,
		VpcPeeringInstanceNo:  ncloud.String(d.Id()),
		VpcPeeringDescription: StringPtrOrNil(d.GetOk("description")),
	}

	LogCommonRequest("setVpcPeeringDescription", reqParams)
	resp, err := config.Client.Vpc.V2Api.SetVpcPeeringDescription(reqParams)
	if err != nil {
		LogErrorResponse("setVpcPeeringDescription", err, reqParams)
		return err
	}
	LogResponse("setVpcPeeringDescription", resp)

	return nil
}
