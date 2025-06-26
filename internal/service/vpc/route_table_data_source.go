package vpc

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudRouteTable() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
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
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudRouteTable(), fieldMap, dataSourceNcloudRouteTableRead)
}

func dataSourceNcloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	resources, err := getRouteTableListFiltered(d, config)
	if err != nil {
		return err
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getRouteTableList(d *schema.ResourceData, config *conn.ProviderConfig) (*vpc.GetRouteTableListResponse, error) {
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

	if v, ok := d.GetOk("id"); ok {
		reqParams.RouteTableNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("GetRouteTableList", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetRouteTableList(reqParams)

	if err != nil {
		LogErrorResponse("GetRouteTableList", err, reqParams)
		return nil, err
	}

	LogResponse("GetRouteTableList", resp)
	return resp, nil
}

func getRouteTableListFiltered(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
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
			"vpc_no":                *r.VpcNo,
			"supported_subnet_type": *r.SupportedSubnetType.Code,
			"is_default":            *r.IsDefault,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, ResourceNcloudRouteTable().Schema)
	}

	return resources, nil
}
