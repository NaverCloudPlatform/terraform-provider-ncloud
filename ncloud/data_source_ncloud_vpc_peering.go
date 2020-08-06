package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudVpcPeering() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudVpcPeeringRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_vpc_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_vpc_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_vpc_login_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceNcloudVpcPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.GetVpcPeeringInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcPeeringName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("target_vpc_name"); ok {
		reqParams.TargetVpcName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("source_vpc_name"); ok {
		reqParams.SourceVpcName = ncloud.String(v.(string))
	}

	logCommonRequest("data_source_ncloud_vpc_peering > GetVpcPeeringInstanceList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetVpcPeeringInstanceList(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_vpc_peering > GetVpcPeeringInstanceList", err, reqParams)
		return err
	}
	logResponse("data_source_ncloud_vpc_peering > GetVpcPeeringInstanceList", resp)

	if err := validateOneResult(len(resp.VpcPeeringInstanceList)); err != nil {
		return err
	}

	instance := resp.VpcPeeringInstanceList[0]

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
