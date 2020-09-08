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

func resourceNcloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudSubnetCreate,
		Read:   resourceNcloudSubnetRead,
		Update: resourceNcloudSubnetUpdate,
		Delete: resourceNcloudSubnetDelete,
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
			"vpc_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the VPC that the desired subnet belongs to.",
			},
			"subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(16, 28),
				Description:  "The CIDR block for the subnet.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Available Zone. Get available values using the `data ncloud_zones`.",
			},
			"network_acl_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network ACL No. Get available values using the `default_network_acl_no` from Resource `ncloud_vpc` or Data source `data.ncloud_network_acls`.",
			},
			"subnet_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false),
				Description:  "Internet Gateway Only. PUBLIC(Yes/Public), PRIVATE(No/Private).",
			},
			"usage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"GEN", "LOADB", "BM"}, false),
				Description:  "Usage type. GEN(Normal), LOADB(Load Balance), BM(BareMetal). default : GEN(Normal).",
			},
			"subnet_no": {
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

func resourceNcloudSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

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
		logCommonRequest("resource_ncloud_subnet > CreateSubnet", reqParams)
		resp, err = config.Client.vpc.V2Api.CreateSubnet(reqParams)

		if err == nil {
			return resource.NonRetryableError(err)
		}

		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == "1001015" {
			logErrorResponse("retry resource_ncloud_subnet > CreateSubnet", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
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
	config := meta.(*ProviderConfig)

	instance, err := getSubnetInstance(config, d.Id())
	if err != nil {
		d.SetId("")
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
	d.Set("status", instance.SubnetStatus.Code)
	d.Set("subnet_type", instance.SubnetType.Code)
	d.Set("usage_type", instance.UsageType.Code)
	d.Set("network_acl_no", instance.NetworkAclNo)

	return nil
}

func resourceNcloudSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("network_acl_no") {
		reqParams := &vpc.SetSubnetNetworkAclRequest{
			RegionCode:   &config.RegionCode,
			SubnetNo:     ncloud.String(d.Get("subnet_no").(string)),
			NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
		}

		logCommonRequest("resource_ncloud_subnet > SetSubnetNetworkAcl", reqParams)
		resp, err := config.Client.vpc.V2Api.SetSubnetNetworkAcl(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_subnet > SetSubnetNetworkAcl", err, reqParams)
			return err
		}
		logResponse("resource_ncloud_subnet > SetSubnetNetworkAcl", resp)

		if err := waitForNcloudNetworkACLUpdate(config, d.Id()); err != nil {
			return err
		}
	}

	return resourceNcloudSubnetRead(d, meta)
}

func resourceNcloudSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.DeleteSubnetRequest{
		RegionCode: &config.RegionCode,
		SubnetNo:   ncloud.String(d.Get("subnet_no").(string)),
	}

	logCommonRequest("resource_ncloud_subnet > DeleteSubnet", reqParams)
	resp, err := config.Client.vpc.V2Api.DeleteSubnet(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_subnet > DeleteSubnet", err, reqParams)
		return err
	}
	logResponse("resource_ncloud_subnet > DeleteSubnet", resp)

	if err := waitForNcloudSubnetDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForNcloudSubnetCreation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getSubnetInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "SubnetStatus")
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Subnet (%s) to become available: %s", id, err)
	}

	return nil
}

func waitForNcloudNetworkACLUpdate(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getNetworkACLInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "NetworkAclStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Set network ACL for Subnet (%s) to become running: %s", id, err)
	}

	return nil
}

func waitForNcloudSubnetDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getSubnetInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "SubnetStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Subnet (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func getSubnetInstance(config *ProviderConfig, id string) (*vpc.Subnet, error) {
	reqParams := &vpc.GetSubnetDetailRequest{
		RegionCode: &config.RegionCode,
		SubnetNo:   ncloud.String(id),
	}

	logCommonRequest("resource_ncloud_subnet > GetSubnetDetail", reqParams)
	resp, err := config.Client.vpc.V2Api.GetSubnetDetail(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_subnet > GetSubnetDetail", err, reqParams)
		return nil, err
	}
	logResponse("resource_ncloud_subnet > GetSubnetDetail", resp)

	if len(resp.SubnetList) > 0 {
		instance := resp.SubnetList[0]
		return instance, nil
	}

	return nil, nil
}
