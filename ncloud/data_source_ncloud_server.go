package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_server", dataSourceNcloudServer())
}

func dataSourceNcloudServer() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudServer(), fieldMap, dataSourceNcloudServerRead)
}

func dataSourceNcloudServerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	instances, err := getServerList(d, config)
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	resources := ConvertToArrayMap(instances)
	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudServer().Schema)
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	d.SetId(resources[0]["instance_no"].(string))
	SetSingularResourceDataFromMapSchema(dataSourceNcloudServer(), d, resources[0])
	return nil
}

func getServerList(d *schema.ResourceData, config *ProviderConfig) ([]*ServerInstance, error) {
	if config.SupportVPC {
		return getVpcServerList(d, config)
	} else {
		return getClassicServerList(d, config)
	}
}

func getClassicServerList(d *schema.ResourceData, config *ProviderConfig) ([]*ServerInstance, error) {
	regionNo, err := parseRegionNoParameter(config.Client, d)
	if err != nil {
		return nil, err
	}

	reqParams := &server.GetServerInstanceListRequest{
		RegionNo: regionNo,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.ServerInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getClassicServerList", reqParams)
	resp, err := config.Client.server.V2Api.GetServerInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getClassicServerList", err, reqParams)
		return nil, err
	}

	logResponse("getClassicServerList", resp)

	var list []*ServerInstance
	for _, r := range resp.ServerInstanceList {
		list = append(list, convertClassicServerInstance(r))
	}

	return list, nil
}

func getVpcServerList(d *schema.ResourceData, config *ProviderConfig) ([]*ServerInstance, error) {
	client := config.Client

	reqParams := &vserver.GetServerInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.ServerInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getVpcServerList", reqParams)

	resp, err := client.vserver.V2Api.GetServerInstanceList(reqParams)
	if err != nil {
		logErrorResponse("getVpcServerList", err, reqParams)
		return nil, err
	}
	logResponse("getVpcServerList", resp)

	var list []*ServerInstance
	for _, r := range resp.ServerInstanceList {
		list = append(list, convertVcpServerInstance(r))
	}

	return list, nil
}
