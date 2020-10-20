package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_vpc_peering", dataSourceNcloudVpcPeering())
}

func dataSourceNcloudVpcPeering() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"vpc_peering_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"source_vpc_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"target_vpc_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudVpcPeering(), fieldMap, dataSourceNcloudVpcPeeringRead)
}

func dataSourceNcloudVpcPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_vpc_peering`")
	}

	resources, err := getVpcPeeringListFiltered(d, config)
	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getVpcPeeringListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vpc.GetVpcPeeringInstanceListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("vpc_peering_no"); ok {
		reqParams.VpcPeeringInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcPeeringName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("target_vpc_name"); ok {
		reqParams.TargetVpcName = ncloud.String(v.(string))
	}

	if v, ok := d.GetOk("source_vpc_name"); ok {
		reqParams.SourceVpcName = ncloud.String(v.(string))
	}

	logCommonRequest("GetVpcPeeringInstanceList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetVpcPeeringInstanceList(reqParams)

	if err != nil {
		logErrorResponse("GetVpcPeeringInstanceList", err, reqParams)
		return nil, err
	}
	logResponse("GetVpcPeeringInstanceList", resp)

	resources := []map[string]interface{}{}

	for _, r := range resp.VpcPeeringInstanceList {
		instance := map[string]interface{}{
			"id":                      *r.VpcPeeringInstanceNo,
			"vpc_peering_no":          *r.VpcPeeringInstanceNo,
			"name":                    *r.VpcPeeringName,
			"description":             *r.VpcPeeringDescription,
			"source_vpc_no":           *r.SourceVpcNo,
			"target_vpc_no":           *r.TargetVpcNo,
			"target_vpc_name":         *r.TargetVpcName,
			"target_vpc_login_id":     *r.TargetVpcLoginId,
			"has_reverse_vpc_peering": *r.HasReverseVpcPeering,
			"is_between_accounts":     *r.IsBetweenAccounts,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudVpcPeering().Schema)
	}

	return resources, nil
}
