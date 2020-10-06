package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_placement_group", dataSourceNcloudPlacementGroup())
}

func dataSourceNcloudPlacementGroup() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"placement_group_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validateInstanceName,
		},
		"placement_group_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"AA"}, false),
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudPlacementGroup(), fieldMap, dataSourceNcloudPlacementGroupRead)
}

func dataSourceNcloudPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("data source `ncloud_placement_group`")
	}

	resources, err := getPlacementGroupListFiltered(d, config)

	if err != nil {
		return err
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getPlacementGroupList(d *schema.ResourceData, config *ProviderConfig) (*vserver.GetPlacementGroupListResponse, error) {
	reqParams := &vserver.GetPlacementGroupListRequest{
		RegionCode: &config.RegionCode,
	}

	if v, ok := d.GetOk("placement_group_no"); ok {
		reqParams.PlacementGroupNoList = []*string{ncloud.String(v.(string))}
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.PlacementGroupName = ncloud.String(v.(string))
	}

	logCommonRequest("GetPlacementGroupList", reqParams)
	resp, err := config.Client.vserver.V2Api.GetPlacementGroupList(reqParams)

	if err != nil {
		logErrorResponse("GetPlacementGroupList", err, reqParams)
		return nil, err
	}

	logResponse("GetPlacementGroupList", resp)
	return resp, nil
}

func getPlacementGroupListFiltered(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
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
		resources = ApplyFilters(f.(*schema.Set), resources, resourceNcloudPlacementGroup().Schema)
	}

	return resources, nil
}
