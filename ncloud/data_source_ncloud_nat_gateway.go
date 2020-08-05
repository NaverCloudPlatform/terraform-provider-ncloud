package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudNatGatewayRead,

		Schema: map[string]*schema.Schema{
			"nat_gateway_no": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "NAT Gateway No. The id of the NAT Gateway for lookup",
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
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

func dataSourceNcloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	regionCode := meta.(*ProviderConfig).RegionCode

	reqParams := &vpc.GetNatGatewayInstanceDetailRequest{
		NatGatewayInstanceNo: ncloud.String(d.Get("nat_gateway_no").(string)),
		RegionCode:           &regionCode,
	}

	logCommonRequest("data_source_ncloud_nat_gateway > GetNatGatewayInstanceDetail", reqParams)
	resp, err := client.vpc.V2Api.GetNatGatewayInstanceDetail(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_nat_gateway > GetNatGatewayInstanceDetail", err, reqParams)
		return err
	}

	logResponse("data_source_ncloud_nat_gateway > GetNatGatewayInstanceDetail", resp)

	instanceList := resp.NatGatewayInstanceList

	if err := validateOneResult(len(instanceList)); err != nil {
		return err
	}

	instance := instanceList[0]

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
