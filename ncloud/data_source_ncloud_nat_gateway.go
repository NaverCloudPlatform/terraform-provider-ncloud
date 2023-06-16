package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_nat_gateway", dataSourceNcloudNatGateway())
}

func dataSourceNcloudNatGateway() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vpc_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudNatGateway(), fieldMap, dataSourceNcloudNatGatewayRead)
}

func dataSourceNcloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_nat_gateway`")
	}

	resources, err := getNatGatewayListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getNatGatewayListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vpc.GetNatGatewayInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.NatGatewayInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NatGatewayName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_name"); ok {
		reqParams.VpcName = ncloud.String(v.(string))
	}

	logCommonRequest("GetNatGatewayInstanceList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNatGatewayInstanceList(reqParams)

	if err != nil {
		logErrorResponse("GetNatGatewayInstanceList", err, reqParams)
		return nil, err
	}

	logResponse("GetNatGatewayInstanceList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.NatGatewayInstanceList {
		instance := map[string]interface{}{
			"id":             *r.NatGatewayInstanceNo,
			"nat_gateway_no": *r.NatGatewayInstanceNo,
			"name":           *r.NatGatewayName,
			"description":    *r.NatGatewayDescription,
			"public_ip":      *r.PublicIp,
			"vpc_no":         *r.VpcNo,
			"vpc_name":       *r.VpcName,
			"zone":           *r.ZoneCode,
			"subnet_no":      *r.SubnetNo,
			"subnet_name":    *r.SubnetName,
			"private_ip":     *r.PrivateIp,
			"public_ip_no":   *r.PublicIpInstanceNo,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudNatGateway().Schema)
	}

	return resources, nil
}
