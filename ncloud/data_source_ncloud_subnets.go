package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceNcloudSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudSubnetsRead,

		Schema: map[string]*schema.Schema{
			"subnet_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of subnet ID to retrieve",
			},
			"vpc_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC ID that you want to filter from",
			},
			"subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CIDR block for the subnet.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Available Zone. Get available values using the `data ncloud_zones`.",
			},
			"network_acl_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network ACL No. Get available values using the `default_network_acl_no` from Resource `ncloud_vpc` or Data source `data.ncloud_network_acls`.",
			},
			"subnet_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false),
				Description:  "Internet Gateway Only. PUBLC(Yes/Public), PRIVATE(No/Private).",
			},
			"usage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"GEN", "LOADB", "BM"}, false),
				Description:  "Usage type. GEN(Normal), LOADB(Load Balance), BM(BareMetal). default : GEN(Normal).",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"INIT", "CREATING", "RUN", "TERMTING"}, false),
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceNcloudSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	regionCode := meta.(*ProviderConfig).RegionCode

	reqParams := &vpc.GetSubnetListRequest{
		RegionCode: &regionCode,
	}

	if v, ok := d.GetOk("subnet_no_list"); ok {
		reqParams.SubnetNoList = expandStringInterfaceList(v.([]interface{}))
	}

	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNo = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("subnet"); ok {
		reqParams.Subnet = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		reqParams.ZoneCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("network_acl_no"); ok {
		reqParams.NetworkAclNo = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_type_code"); ok {
		reqParams.SubnetTypeCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("usage_type_code"); ok {
		reqParams.UsageTypeCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		reqParams.SubnetStatusCode = ncloud.String(v.(string))
	}

	logCommonRequest("data_source_ncloud_subnets > GetSubnetList", reqParams)
	resp, err := client.vpc.V2Api.GetSubnetList(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_subnets > GetSubnetList", err, reqParams)
		return err
	}
	logResponse("data_source_ncloud_subnets > GetSubnetList", resp)

	if resp == nil || len(resp.SubnetList) == 0 {
		return fmt.Errorf("no matching Subnets found")
	}

	instances := make([]string, 0)

	for _, vpc := range resp.SubnetList {
		instances = append(instances, ncloud.StringValue(vpc.VpcNo))
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("ids", instances); err != nil {
		return fmt.Errorf("Error setting Subnets ids: %s", err)
	}

	return nil
}
