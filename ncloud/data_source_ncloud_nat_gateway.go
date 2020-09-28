package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterDatasource("ncloud_nat_gateway", dataSourceNcloudNatGateway())
}

func dataSourceNcloudNatGateway() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"nat_gateway_no": {
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
		"status": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudNatGateway(), fieldMap, dataSourceNcloudNatGatewayRead)
}

func dataSourceNcloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

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

	if v, ok := d.GetOk("nat_gateway_no"); ok {
		reqParams.NatGatewayInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.NatGatewayName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_name"); ok {
		reqParams.VpcName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		reqParams.NatGatewayInstanceStatusCode = ncloud.String(v.(string))
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
			"status":         *r.NatGatewayInstanceStatus.Code,
			"vpc_no":         *r.VpcNo,
			"zone":           *r.ZoneCode,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudNatGateway().Schema)
	}

	return resources, nil
}
