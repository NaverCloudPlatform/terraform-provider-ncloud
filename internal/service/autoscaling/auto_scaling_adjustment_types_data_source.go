package autoscaling

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func DataSourceNcloudAutoScalingAdjustmentTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudAutoScalingAdjustmentTypesRead,

		Schema: map[string]*schema.Schema{
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"filter": DataSourceFiltersSchema(),
		},
	}
}

func dataSourceNcloudAutoScalingAdjustmentTypesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	resources, err := getAutoScalingAdjustmentListFiltered(d, config)

	if err != nil {
		return err
	}

	types := make([]map[string]interface{}, len(resources))
	for i, r := range resources {
		types[i] = map[string]interface{}{
			"code":      r["code"],
			"code_name": r["code_name"],
		}
	}

	d.SetId("ncloud_auto_scaling_adjustment_types")
	d.Set("types", types)

	return nil
}

func getAutoScalingAdjustmentListFiltered(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	var resources []map[string]interface{}
	var err error

	resources, err = getVpcAutoScalingAdjustmentTypeList(config)
	if err != nil {
		return nil, err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudAutoScalingAdjustmentTypes().Schema)
	}
	return resources, nil
}

func getVpcAutoScalingAdjustmentTypeList(config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vautoscaling.GetAdjustmentTypeListRequest{
		RegionCode: &regionCode,
	}

	LogCommonRequest("GetAdjustmentTypeListRequest", reqParams)

	resp, err := client.Vautoscaling.V2Api.GetAdjustmentTypeList(reqParams)
	if err != nil {
		LogErrorResponse("GetAdjustmentTypeListRequest", err, reqParams)
		return nil, err
	}

	LogResponse("GetAdjustmentTypeListRequest", resp)

	var resources []map[string]interface{}

	for _, r := range resp.AdjustmentTypeList {
		instance := map[string]interface{}{
			"code":      *r.Code,
			"code_name": *r.CodeName,
		}

		resources = append(resources, instance)
	}
	return resources, nil
}
