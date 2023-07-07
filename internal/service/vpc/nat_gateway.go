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

func ResourceNcloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudNatGatewayCreate,
		Read:   resourceNcloudNatGatewayRead,
		Update: resourceNcloudNatGatewayUpdate,
		Delete: resourceNcloudNatGatewayDelete,
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
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringLenBetween(0, 1000)),
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nat_gateway_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_nat_gateway`")
	}

	if _, ok := d.GetOk("subnet_no"); !ok {
		return fmt.Errorf("subnet_no is required when creating a new NATGW")
	}

	reqParams := &vpc.CreateNatGatewayInstanceRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      ncloud.String(d.Get("vpc_no").(string)),
		ZoneCode:   ncloud.String(d.Get("zone").(string)),
		SubnetNo:   ncloud.String(d.Get("subnet_no").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NatGatewayName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		reqParams.NatGatewayDescription = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("private_ip"); ok {
		reqParams.PrivateIp = ncloud.String(v.(string))
	}

	LogCommonRequest("CreateNatGatewayInstance", reqParams)
	resp, err := config.Client.Vpc.V2Api.CreateNatGatewayInstance(reqParams)
	if err != nil {
		LogErrorResponse("CreateNatGatewayInstance", err, reqParams)
		return err
	}

	LogResponse("CreateNatGatewayInstance", resp)

	instance := resp.NatGatewayInstanceList[0]
	d.SetId(*instance.NatGatewayInstanceNo)
	log.Printf("[INFO] NAT Gateway ID: %s", d.Id())

	if err := waitForNcloudNatGatewayCreation(config, d.Id()); err != nil {
		return err
	}

	return resourceNcloudNatGatewayRead(d, meta)
}

func resourceNcloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instance, err := GetNatGatewayInstance(config, d.Id())
	if err != nil {
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
	d.Set("vpc_no", instance.VpcNo)
	d.Set("zone", instance.ZoneCode)
	d.Set("subnet_name", instance.SubnetName)
	d.Set("subnet_no", instance.SubnetNo)
	d.Set("private_ip", instance.PrivateIp)
	d.Set("public_ip_no", instance.PublicIpInstanceNo)

	return nil
}

func resourceNcloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if d.HasChange("description") {
		if err := setNatGatewayDescription(d, config); err != nil {
			return err
		}
	}

	return resourceNcloudNatGatewayRead(d, meta)
}

func resourceNcloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	reqParams := &vpc.DeleteNatGatewayInstanceRequest{
		RegionCode:           &config.RegionCode,
		NatGatewayInstanceNo: ncloud.String(d.Get("nat_gateway_no").(string)),
	}

	LogCommonRequest("DeleteNatGatewayInstance", reqParams)
	resp, err := config.Client.Vpc.V2Api.DeleteNatGatewayInstance(reqParams)
	if err != nil {
		LogErrorResponse("DeleteNatGatewayInstance", err, reqParams)
		return err
	}

	LogResponse("DeleteNatGatewayInstance", resp)

	if err := WaitForNcloudNatGatewayDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForNcloudNatGatewayCreation(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNatGatewayInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NatGatewayInstanceStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become available: %s", id, err)
	}

	return nil
}

func WaitForNcloudNatGatewayDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNatGatewayInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NatGatewayInstanceStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NAT Gateway (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetNatGatewayInstance(config *conn.ProviderConfig, id string) (*vpc.NatGatewayInstance, error) {
	reqParams := &vpc.GetNatGatewayInstanceDetailRequest{
		RegionCode:           &config.RegionCode,
		NatGatewayInstanceNo: ncloud.String(id),
	}

	LogCommonRequest("GetNatGatewayInstanceDetail", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetNatGatewayInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("GetNatGatewayInstanceDetail", err, reqParams)
		return nil, err
	}
	LogResponse("GetNatGatewayInstanceDetail", resp)

	if len(resp.NatGatewayInstanceList) > 0 {
		instance := resp.NatGatewayInstanceList[0]
		return instance, nil
	}

	return nil, nil
}

func setNatGatewayDescription(d *schema.ResourceData, config *conn.ProviderConfig) error {
	reqParams := &vpc.SetNatGatewayDescriptionRequest{
		RegionCode:            &config.RegionCode,
		NatGatewayInstanceNo:  ncloud.String(d.Id()),
		NatGatewayDescription: StringPtrOrNil(d.GetOk("description")),
	}

	LogCommonRequest("setNatGatewayDescription", reqParams)
	resp, err := config.Client.Vpc.V2Api.SetNatGatewayDescription(reqParams)
	if err != nil {
		LogErrorResponse("setNatGatewayDescription", err, reqParams)
		return err
	}
	LogResponse("setNatGatewayDescription", resp)

	return nil
}
