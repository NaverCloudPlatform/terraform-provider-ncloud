package server

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudNetworkInterface() *schema.Resource {
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
		"private_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudNetworkInterface(), fieldMap, dataSourceNcloudNetworkInterfaceRead)
}

func dataSourceNcloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)
	var resources []map[string]interface{}
	var err error

	resources, err = getVpcNetworkInterfaceListFiltered(d, config)
	if err != nil {
		return err
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getVpcNetworkInterfaceListFiltered(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vserver.GetNetworkInterfaceListRequest{
		RegionCode:           &config.RegionCode,
		Ip:                   StringPtrOrNil(d.GetOk("private_ip")),
		NetworkInterfaceName: StringPtrOrNil(d.GetOk("name")),
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.NetworkInterfaceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getVpcNetworkInterfaceList", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetNetworkInterfaceList(reqParams)

	if err != nil {
		LogErrorResponse("getVpcNetworkInterfaceList", err, reqParams)
		return nil, err
	}
	LogResponse("getVpcNetworkInterfaceList", resp)

	var resources []map[string]interface{}

	for _, r := range resp.NetworkInterfaceList {
		instance := map[string]interface{}{
			"id":                   *r.NetworkInterfaceNo,
			"network_interface_no": *r.NetworkInterfaceNo,
			"name":                 *r.NetworkInterfaceName,
			"description":          *r.NetworkInterfaceDescription,
			"subnet_no":            *r.SubnetNo,
			"private_ip":           *r.Ip,
			"status":               *r.NetworkInterfaceStatus.Code,
			"is_default":           *r.IsDefault,
			"server_instance_no":   StringOrEmpty(r.InstanceNo),
		}

		if r.AccessControlGroupNoList != nil {
			instance["access_control_groups"] = StringPtrArrToStringArr(r.AccessControlGroupNoList)
		}

		if r.InstanceType != nil {
			instance["instance_type"] = *r.InstanceType.Code
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, ResourceNcloudNetworkInterface().Schema)
	}

	return resources, nil
}
