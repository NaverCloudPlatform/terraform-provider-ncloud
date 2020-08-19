package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNcloudVpcPeering() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
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
		"status": {
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

	resources, err := getVpcPeeringListFiltered(d, config)
	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	for k, v := range resources[0] {
		if k == "id" {
			d.SetId(v.(string))
			continue
		}
		d.Set(k, v)
	}

	return nil
}

func getVpcPeeringListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	reqParams := &vpc.GetVpcPeeringInstanceListRequest{
		RegionCode: &config.RegionCode,
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

	if v, ok := d.GetOk("status"); ok {
		reqParams.VpcPeeringInstanceStatusCode = ncloud.String(v.(string))
	}

	logCommonRequest("data_source_ncloud_vpc_peering > GetVpcPeeringInstanceList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetVpcPeeringInstanceList(reqParams)

	if err != nil {
		logErrorResponse("data_source_ncloud_vpc_peering > GetVpcPeeringInstanceList", err, reqParams)
		return nil, err
	}
	logResponse("data_source_ncloud_vpc_peering > GetVpcPeeringInstanceList", resp)

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
			"status":                  *r.VpcPeeringInstanceStatus.Code,
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
