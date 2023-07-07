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

const (
	SubnetPleaseTryAgainErrorCode = "3000"
)

func ResourceNcloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudSubnetCreate,
		Read:   resourceNcloudSubnetRead,
		Update: resourceNcloudSubnetUpdate,
		Delete: resourceNcloudSubnetDelete,
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
				Description:      "Subnet name to create. default: Assigned by NAVER CLOUD PLATFORM.",
			},
			"vpc_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the VPC that the desired subnet belongs to.",
			},
			"subnet": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.IsCIDRNetwork(16, 28)),
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_acl_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false)),
			},
			"usage_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"GEN", "LOADB", "BM", "NATGW"}, false)),
			},
			"subnet_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_subnet`")
	}

	reqParams := &vpc.CreateSubnetRequest{
		RegionCode:     &config.RegionCode,
		Subnet:         ncloud.String(d.Get("subnet").(string)),
		SubnetTypeCode: ncloud.String(d.Get("subnet_type").(string)),
		UsageTypeCode:  ncloud.String(d.Get("usage_type").(string)),
		NetworkAclNo:   ncloud.String(d.Get("network_acl_no").(string)),
		VpcNo:          ncloud.String(d.Get("vpc_no").(string)),
		ZoneCode:       ncloud.String(d.Get("zone").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.SubnetName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("usage_type"); ok {
		reqParams.UsageTypeCode = ncloud.String(v.(string))
	}

	var resp *vpc.CreateSubnetResponse
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		LogCommonRequest("CreateSubnet", reqParams)
		resp, err = config.Client.Vpc.V2Api.CreateSubnet(reqParams)

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == "1001015" || errBody.ReturnCode == SubnetPleaseTryAgainErrorCode {
				LogErrorResponse("retry CreateSubnet", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	instance := resp.SubnetList[0]
	d.SetId(*instance.SubnetNo)
	log.Printf("[INFO] Subnet ID: %s", d.Id())

	if err := waitForNcloudSubnetCreation(config, d.Id()); err != nil {
		return err
	}

	return resourceNcloudSubnetRead(d, meta)
}

func resourceNcloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instance, err := GetSubnetInstance(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.SubnetNo)
	d.Set("subnet_no", instance.SubnetNo)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("zone", instance.ZoneCode)
	d.Set("name", instance.SubnetName)
	d.Set("subnet", instance.Subnet)
	d.Set("subnet_type", instance.SubnetType.Code)
	d.Set("usage_type", instance.UsageType.Code)
	d.Set("network_acl_no", instance.NetworkAclNo)

	return nil
}

func resourceNcloudSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if d.HasChange("network_acl_no") {
		reqParams := &vpc.SetSubnetNetworkAclRequest{
			RegionCode:   &config.RegionCode,
			SubnetNo:     ncloud.String(d.Get("subnet_no").(string)),
			NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
		}

		LogCommonRequest("SetSubnetNetworkAcl", reqParams)
		resp, err := config.Client.Vpc.V2Api.SetSubnetNetworkAcl(reqParams)
		if err != nil {
			LogErrorResponse("SetSubnetNetworkAcl", err, reqParams)
			return err
		}
		LogResponse("SetSubnetNetworkAcl", resp)

		if err := waitForNcloudNetworkACLUpdate(config, d.Get("network_acl_no").(string)); err != nil {
			return err
		}
	}

	return resourceNcloudSubnetRead(d, meta)
}

func resourceNcloudSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	reqParams := &vpc.DeleteSubnetRequest{
		RegionCode: &config.RegionCode,
		SubnetNo:   ncloud.String(d.Get("subnet_no").(string)),
	}

	LogCommonRequest("DeleteSubnet", reqParams)
	resp, err := config.Client.Vpc.V2Api.DeleteSubnet(reqParams)
	if err != nil {
		LogErrorResponse("DeleteSubnet", err, reqParams)
		return err
	}
	LogResponse("DeleteSubnet", resp)

	if err := WaitForNcloudSubnetDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForNcloudSubnetCreation(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetSubnetInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "SubnetStatus")
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Subnet (%s) to become available: %s", id, err)
	}

	return nil
}

func waitForNcloudNetworkACLUpdate(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetNetworkACLInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Set network ACL for Subnet (%s) to become running: %s", id, err)
	}

	return nil
}

func WaitForNcloudSubnetDeletion(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetSubnetInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "SubnetStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Subnet (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func GetSubnetInstance(config *conn.ProviderConfig, id string) (*vpc.Subnet, error) {
	reqParams := &vpc.GetSubnetDetailRequest{
		RegionCode: &config.RegionCode,
		SubnetNo:   ncloud.String(id),
	}

	LogCommonRequest("GetSubnetDetail", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetSubnetDetail(reqParams)
	if err != nil {
		LogErrorResponse("GetSubnetDetail", err, reqParams)
		return nil, err
	}
	LogResponse("GetSubnetDetail", resp)

	if len(resp.SubnetList) > 0 {
		instance := resp.SubnetList[0]
		return instance, nil
	}

	return nil, nil
}
