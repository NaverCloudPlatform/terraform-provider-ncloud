package autoscaling

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudAutoScalingPolicy() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"auto_scaling_group_no": {
			Type:     schema.TypeString,
			Required: true,
		},
		"filter": DataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(ResourceNcloudAutoScalingPolicy(), fieldMap, dataSourceNcloudAutoScalingPolicyRead)
}

func dataSourceNcloudAutoScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	policyList, err := getAutoScalingPolicyList(d, config)
	if err != nil {
		return err
	}

	policyListMap := ConvertToArrayMap(policyList)
	if f, ok := d.GetOk("filter"); ok {
		policyListMap = ApplyFilters(f.(*schema.Set), policyListMap, DataSourceNcloudAutoScalingPolicy().Schema)
	}

	if err := ValidateOneResult(len(policyListMap)); err != nil {
		return err
	}

	d.SetId(policyListMap[0]["name"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudAutoScalingPolicy(), d, policyListMap[0])
	return nil
}

func getAutoScalingPolicyList(d *schema.ResourceData, config *conn.ProviderConfig) ([]*AutoScalingPolicy, error) {
	reqParams := &vautoscaling.GetAutoScalingPolicyListRequest{
		RegionCode:         &config.RegionCode,
		AutoScalingGroupNo: ncloud.String(d.Get("auto_scaling_group_no").(string)),
	}

	resp, err := config.Client.Vautoscaling.V2Api.GetAutoScalingPolicyList(reqParams)
	if err != nil {
		return nil, err
	}

	list := make([]*AutoScalingPolicy, 0)
	for _, p := range resp.ScalingPolicyList {
		policy := &AutoScalingPolicy{
			AutoScalingPolicyName: p.PolicyName,
			AutoScalingPolicyNo:   p.PolicyNo,
			AutoScalingGroupNo:    p.AutoScalingGroupNo,
			AdjustmentTypeCode:    p.AdjustmentType.Code,
			ScalingAdjustment:     p.ScalingAdjustment,
			Cooldown:              p.CoolDown,
			MinAdjustmentStep:     p.MinAdjustmentStep,
		}
		if *p.PolicyName == d.Id() {
			return []*AutoScalingPolicy{policy}, nil
		}
		list = append(list, policy)
	}

	if d.Id() != "" {
		return nil, nil
	}
	return list, nil
}
