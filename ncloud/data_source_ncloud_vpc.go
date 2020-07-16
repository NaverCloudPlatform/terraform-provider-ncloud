package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudVpcRead,

		Schema: map[string]*schema.Schema{
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_cidr_block": {
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

func dataSourceNcloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.GetVpcDetailRequest{
		VpcNo:      ncloud.String(d.Get("vpc_no").(string)),
		RegionCode: regionCode,
	}

	logCommonRequest("GetVpcDetail", reqParams)
	resp, err := client.vpc.V2Api.GetVpcDetail(reqParams)

	if err != nil {
		logErrorResponse("Get Vpc Instance", err, reqParams)
		return err
	}

	logResponse("GetVpcDetail", resp)

	vpcInstanceList := resp.VpcList

	if err := validateOneResult(len(vpcInstanceList)); err != nil {
		return err
	}

	vpcInstance := vpcInstanceList[0]

	d.SetId(*vpcInstance.VpcNo)
	d.Set("vpc_no", vpcInstance.VpcNo)
	d.Set("name", vpcInstance.VpcName)
	d.Set("ipv4_cidr_block", vpcInstance.Ipv4CidrBlock)
	d.Set("status", vpcInstance.VpcStatus.Code)

	return nil
}
