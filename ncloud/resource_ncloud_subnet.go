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

func resourceNcloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudSubnetCreate,
		Read:   resourceNcloudSubnetRead,
		Update: resourceNcloudSubnetUpdate,
		Delete: resourceNcloudSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceNcloudSubnetCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 30),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-*]+$`), "Composed of alphabets, numbers, hyphen (-) and wild card (*)."),
					validation.StringMatch(regexp.MustCompile(`.*[^\\-]$`), "Hyphen (-) cannot be used for the last character and if wild card (*) is used, other characters cannot be input."),
				),
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(16, 28),
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false),
			},
			"usage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"GEN", "LOADB"}, false),
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
	client := meta.(*NcloudAPIClient)
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.CreateSubnetRequest{
		SubnetName:     ncloud.String(d.Get("name").(string)),
		Subnet:         ncloud.String(d.Get("subnet").(string)),
		SubnetTypeCode: ncloud.String(d.Get("subnet_type").(string)),
		UsageTypeCode:  ncloud.String(d.Get("usage_type").(string)),
		NetworkAclNo:   ncloud.String(d.Get("network_acl_no").(string)),
		VpcNo:          ncloud.String(d.Get("vpc_no").(string)),
		ZoneCode:       ncloud.String(d.Get("zone").(string)),
		RegionCode:     regionCode,
	}

	if v, ok := d.GetOk("usage_type"); ok {
		reqParams.UsageTypeCode = ncloud.String(v.(string))
	}

	logCommonRequest("resource_ncloud_subnet > CreateSubnet", reqParams)
	resp, err := client.vpc.V2Api.CreateSubnet(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_subnet > CreateSubnet", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_subnet > CreateSubnet", resp)

	instance := resp.SubnetList[0]
	d.SetId(*instance.SubnetNo)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INIT", "CREATING"},
		Target:     []string{"RUN"},
		Refresh:    SubnetStateRefreshFunc(client, d.Id()),
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for Subnet (%s) to become available: %s",
			d.Id(), err)
	}

	log.Printf("[INFO] Subnet ID: %s", d.Id())

	return resourceNcloudSubnetRead(d, meta)
}

func resourceNcloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	instance, err := getSubnetInstance(client, d.Id())
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
	client := meta.(*NcloudAPIClient)

	if d.HasChange("network_acl_no") {
		reqParams := &vpc.SetSubnetNetworkAclRequest{
			SubnetNo:     ncloud.String(d.Get("subnet_no").(string)),
			NetworkAclNo: ncloud.String(d.Get("network_acl_no").(string)),
		}

		logCommonRequest("resource_ncloud_subnet > SetSubnetNetworkAcl", reqParams)
		resp, err := client.vpc.V2Api.SetSubnetNetworkAcl(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_subnet > SetSubnetNetworkAcl", err, reqParams)
			return err
		}
		logResponse("resource_ncloud_subnet > SetSubnetNetworkAcl", resp)
	}

	return resourceNcloudSubnetRead(d, meta)
}

func resourceNcloudSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.DeleteSubnetRequest{
		SubnetNo:   ncloud.String(d.Get("subnet_no").(string)),
		RegionCode: regionCode,
	}

	logCommonRequest("resource_ncloud_subnet > DeleteSubnet", reqParams)
	resp, err := client.vpc.V2Api.DeleteSubnet(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_subnet > DeleteSubnet", err, reqParams)
		return err
	}
	logResponse("resource_ncloud_subnet > DeleteSubnet", resp)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"RUN", "TERMTING"},
		Target:     []string{"TERMINATED"},
		Refresh:    SubnetStateRefreshFunc(client, d.Id()),
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for Subnet (%s) to become termintaing: %s",
			d.Id(), err)
	}

	return nil
}

func getSubnetInstance(client *NcloudAPIClient, id string) (*vpc.Subnet, error) {
	reqParams := &vpc.GetSubnetDetailRequest{
		SubnetNo: ncloud.String(id),
	}

	logCommonRequest("resource_ncloud_subnet > GetSubnetDetail", reqParams)
	resp, err := client.vpc.V2Api.GetSubnetDetail(reqParams)
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

// SubnetStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch a Subnet
func SubnetStateRefreshFunc(client *NcloudAPIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getSubnetInstance(client, id)

		if err != nil {
			return nil, "", err
		}

		if instance == nil {
			return instance, "TERMINATED", nil
		}

		return instance, *instance.SubnetStatus.Code, nil
	}
}

func resourceNcloudSubnetCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChange("name") {
		old, new := diff.GetChange("name")
		if len(old.(string)) > 0 {
			return fmt.Errorf("Change 'name' is not support, Please set name as a old value = [%s -> %s]", new, old)
		}
	}

	return nil
}
