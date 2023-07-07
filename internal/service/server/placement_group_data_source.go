package server

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudPlacementGroup() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: verify.ToDiagFunc(verify.ValidateInstanceName),
		},
		"placement_group_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"AA"}, false)),
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudPlacementGroup(), fieldMap, dataSourceNcloudPlacementGroupRead)
}

func dataSourceNcloudPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_placement_group`")
	}

	resources, err := getPlacementGroupListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getPlacementGroupList(d *schema.ResourceData, config *conn.ProviderConfig) (*vserver.GetPlacementGroupListResponse, error) {
	reqParams := &vserver.GetPlacementGroupListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.PlacementGroupNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.PlacementGroupName = ncloud.String(v.(string))
	}

	LogCommonRequest("GetPlacementGroupList", reqParams)
	resp, err := config.Client.Vserver.V2Api.GetPlacementGroupList(reqParams)

	if err != nil {
		LogErrorResponse("GetPlacementGroupList", err, reqParams)
		return nil, err
	}

	LogResponse("GetPlacementGroupList", resp)
	return resp, nil
}

func getPlacementGroupListFiltered(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	resp, err := getPlacementGroupList(d, config)

	if err != nil {
		return nil, err
	}

	resources := []map[string]interface{}{}

	for _, r := range resp.PlacementGroupList {
		instance := map[string]interface{}{
			"id":                   *r.PlacementGroupNo,
			"placement_group_no":   *r.PlacementGroupNo,
			"name":                 *r.PlacementGroupName,
			"placement_group_type": *r.PlacementGroupType.Code,
		}

		resources = append(resources, instance)
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, ResourceNcloudPlacementGroup().Schema)
	}

	return resources, nil
}
