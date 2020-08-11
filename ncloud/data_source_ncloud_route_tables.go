package ncloud

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRouteTablesRead,

		Schema: map[string]*schema.Schema{
			"vpc_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"supported_subnet_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": dataSourceFiltersSchema(),
			"route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(resourceNcloudRouteTable()),
			},
		},
	}
}

func dataSourceNcloudRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	resp, err := getRouteTableList(d, config)

	if err != nil {
		return err
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
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudRouteTables().Schema["route_tables"].Elem.(*schema.Resource).Schema)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("route_tables", resources); err != nil {
		return fmt.Errorf("Error setting route table ids: %s", err)
	}

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

	logCommonRequest("data_source_ncloud_route_tables > GetRouteTableList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetRouteTableList(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_route_tables > GetRouteTableList", err, reqParams)
		return nil, err
	}

	logResponse("data_source_ncloud_route_tables > GetRouteTableList", resp)
	return resp, nil
}
