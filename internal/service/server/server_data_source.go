package server

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudServer() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudServer(), fieldMap, dataSourceNcloudServerRead)
}

func dataSourceNcloudServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	instances, err := getServerList(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudServer().Schema)
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	d.SetId(resources[0]["instance_no"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudServer(), d, resources[0])
	return nil
}

func getServerList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*ServerInstance, error) {
	if config.SupportVPC {
		return getVpcServerList(d, config)
	} else {
		return getClassicServerList(d, config)
	}
}

func getClassicServerList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*ServerInstance, error) {
	regionNo, err := conn.ParseRegionNoParameter(d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.GetServerInstanceListRequest{
		RegionNo: regionNo,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.ServerInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getClassicServerList", reqParams)
	resp, err := config.Client.Server.V2Api.GetServerInstanceList(reqParams)

	if err != nil {
		LogErrorResponse("getClassicServerList", err, reqParams)
		return nil, err
	}

	LogResponse("getClassicServerList", resp)

	var list []*ServerInstance
	for _, r := range resp.ServerInstanceList {
		list = append(list, convertClassicServerInstance(r))
	}

	return list, nil
}

func getVpcServerList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*ServerInstance, error) {
	client := config.Client

	reqParams := &vserver.GetServerInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.ServerInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getVpcServerList", reqParams)

	resp, err := client.Vserver.V2Api.GetServerInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("getVpcServerList", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcServerList", resp)

	var list []*ServerInstance
	for _, r := range resp.ServerInstanceList {
		list = append(list, convertVcpServerInstance(r))
	}

	return list, nil
}
