package ncloud

import (
	"fmt"
	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudRouteTableRead,

		Schema: map[string]*schema.Schema{
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
			"route_table_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
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

func dataSourceNcloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	regionCode, err := parseRegionCodeParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := &vpc.GetRouteTableListRequest{
		RegionCode: regionCode,
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

	logCommonRequest("data_source_ncloud_route_table > GetRouteTableList", reqParams)
	resp, err := client.vpc.V2Api.GetRouteTableList(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_route_table > GetRouteTableList", err, reqParams)
		return err
	}

	logResponse("data_source_ncloud_route_table > GetRouteTableList", resp)

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

	if err := validateOneResult(len(instanceList)); err != nil {
		return err
	}

	instance := instanceList[0]

	d.SetId(*instance.RouteTableNo)
	d.Set("route_table_no", instance.RouteTableNo)
	d.Set("name", instance.RouteTableName)
	d.Set("description", instance.RouteTableDescription)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("supported_subnet_type", instance.SupportedSubnetType.Code)
	d.Set("is_default", instance.IsDefault)
	d.Set("status", instance.RouteTableStatus.Code)

	return nil
}
