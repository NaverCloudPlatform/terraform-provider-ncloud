package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"subnet_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet No. The id of the subnet for lookup",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_acl_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage_type": {
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

func dataSourceNcloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.GetSubnetDetailRequest{
		RegionCode: &config.RegionCode,
		SubnetNo:   ncloud.String(d.Get("subnet_no").(string)),
	}

	logCommonRequest("data_source_ncloud_subnet > GetSubnetDetail", reqParams)
	resp, err := config.Client.vpc.V2Api.GetSubnetDetail(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_subnet > GetSubnetDetail", err, reqParams)
		return err
	}

	logResponse("data_source_ncloud_subnet > GetSubnetDetail", resp)

	instanceList := resp.SubnetList

	if err := validateOneResult(len(instanceList)); err != nil {
		return err
	}

	instance := instanceList[0]

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
