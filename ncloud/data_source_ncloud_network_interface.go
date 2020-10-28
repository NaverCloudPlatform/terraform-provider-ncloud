package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_network_interface", dataSourceNcloudNetworkInterface())
}

func dataSourceNcloudNetworkInterface() *schema.Resource {
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
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudNetworkInterface(), fieldMap, dataSourceNcloudNetworkInterfaceRead)
}

func dataSourceNcloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcNetworkInterfaceListFiltered(d, config)
	} else {
		return NotSupportClassic("data source `ncloud_network_interface`")
	}

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getVpcNetworkInterfaceListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vserver.GetNetworkInterfaceListRequest{
		RegionCode:           &config.RegionCode,
		Ip:                   StringPtrOrNil(d.GetOk("private_ip")),
		NetworkInterfaceName: StringPtrOrNil(d.GetOk("name")),
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.NetworkInterfaceNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getVpcNetworkInterfaceList", reqParams)
	resp, err := config.Client.vserver.V2Api.GetNetworkInterfaceList(reqParams)

	if err != nil {
		logErrorResponse("getVpcNetworkInterfaceList", err, reqParams)
		return nil, err
	}
	logResponse("getVpcNetworkInterfaceList", resp)

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
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudNetworkInterface().Schema)
	}

	return resources, nil
}
