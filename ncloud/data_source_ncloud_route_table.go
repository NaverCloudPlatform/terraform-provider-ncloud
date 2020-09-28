package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterDatasource("ncloud_route_table", dataSourceNcloudRouteTable())
}

func dataSourceNcloudRouteTable() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"route_table_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vpc_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"supported_subnet_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudRouteTable(), fieldMap, dataSourceNcloudRouteTableRead)
}

func dataSourceNcloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	resources, err := getRouteTableListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getRouteTableList(d *schema.ResourceData, config *ProviderConfig) (*vpc.GetRouteTableListResponse, error) {
	reqParams := &vpc.GetRouteTableListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("vpc_no"); ok {
		reqParams.VpcNo = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.RouteTableName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("supported_subnet_type"); ok {
		reqParams.SupportedSubnetTypeCode = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("route_table_no"); ok {
		reqParams.RouteTableNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("GetRouteTableList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetRouteTableList(reqParams)

	if err != nil {
		logErrorResponse("GetRouteTableList", err, reqParams)
		return nil, err
	}

	logResponse("GetRouteTableList", resp)
	return resp, nil
}

func getRouteTableListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	resp, err := getRouteTableList(d, config)

	if err != nil {
		return nil, err
	}

	resources := []map[string]interface{}{}

	for _, r := range resp.RouteTableList {
		instance := map[string]interface{}{
			"id":                    *r.RouteTableNo,
			"route_table_no":        *r.RouteTableNo,
			"name":                  *r.RouteTableName,
			"description":           *r.RouteTableDescription,
			"status":                *r.RouteTableStatus.Code,
			"vpc_no":                *r.VpcNo,
			"supported_subnet_type": *r.SupportedSubnetType.Code,
			"is_default":            *r.IsDefault,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudRouteTable().Schema)
	}

	return resources, nil
}
