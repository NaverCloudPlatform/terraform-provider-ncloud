package vpc

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudVpc() *schema.Resource {
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
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudVpc(), fieldMap, dataSourceNcloudVpcRead)
}

func dataSourceNcloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_vpc`")
	}

	resources, err := getVpcListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getVpcListFiltered(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vpc.GetVpcListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.VpcNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("GetVpcList", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetVpcList(reqParams)

	if err != nil {
		LogErrorResponse("GetVpcList", err, reqParams)
		return nil, err
	}
	LogResponse("GetVpcList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.VpcList {
		id := *r.VpcNo
		defaultNetworkACLNo, err := getDefaultNetworkACL(config, id)
		if err != nil {
			return nil, fmt.Errorf("error get default network acl for VPC (%s): %s", id, err)
		}
		defaultAcgNo, err := GetDefaultAccessControlGroup(config, id)
		if err != nil {
			return nil, fmt.Errorf("error get default Access Control Group for VPC (%s): %s", id, err)
		}
		publicRouteTableNo, privateRouteTableNo, err := getDefaultRouteTable(config, id)
		if err != nil {
			return nil, fmt.Errorf("error get default Route Table for VPC (%s): %s", id, err)
		}

		instance := map[string]interface{}{
			"id":                              *r.VpcNo,
			"vpc_no":                          *r.VpcNo,
			"name":                            *r.VpcName,
			"ipv4_cidr_block":                 *r.Ipv4CidrBlock,
			"default_network_acl_no":          defaultNetworkACLNo,
			"default_access_control_group_no": defaultAcgNo,
			"default_public_route_table_no":   publicRouteTableNo,
			"default_private_route_table_no":  privateRouteTableNo,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, ResourceNcloudVpc().Schema)
	}

	return resources, nil
}
