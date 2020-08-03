package ncloud

import (
	"fmt"
	"strconv"
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

			"is_default": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceNcloudRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

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

	logCommonRequest("data_source_ncloud_route_tables > GetRouteTableList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetRouteTableList(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_route_tables > GetRouteTableList", err, reqParams)
		return err
	}

	logResponse("data_source_ncloud_route_tables > GetRouteTableList", resp)

	var instanceList []*vpc.RouteTable

	if v, ok := d.GetOk("is_default"); ok {
		isDefault, err := strconv.ParseBool(v.(string))
		if err != nil {
			return fmt.Errorf("invalid attribute: invalid value for is_default: %s", v)
		}

		for _, i := range resp.RouteTableList {
			if *i.IsDefault == isDefault {
				instanceList = append(instanceList, i)
			}
		}
	} else {
		instanceList = resp.RouteTableList
	}

	ids := make([]string, 0)

	for _, instance := range instanceList {
		ids = append(ids, ncloud.StringValue(instance.RouteTableNo))
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("ids", ids); err != nil {
		return fmt.Errorf("Error setting route table ids: %s", err)
	}

	return nil
}
