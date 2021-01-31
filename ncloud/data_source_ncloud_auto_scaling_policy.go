package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_auto_scaling_policy", dataSourceNcloudAutoScalingPolicy())
}

func dataSourceNcloudAutoScalingPolicy() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}

	return GetSingularDataSourceItemSchema(resourceNcloudAutoScalingPolicy(), fieldMap, dataSourceNcloudAutoScalingPolicyRead)
}

func dataSourceNcloudAutoScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if v, ok := d.GetOk("name"); ok {
		d.SetId(v.(string))
	}

	policyList, err := getAutoScalingPolicyList(config, d.Id())
	if err != nil {
		return err
	}

	policyListMap := ConvertToArrayMap(policyList)
	if f, ok := d.GetOk("filter"); ok {
		policyListMap = ApplyFilters(f.(*schema.Set), policyListMap, dataSourceNcloudAutoScalingPolicy().Schema)
	}

	if err := validateOneResult(len(policyListMap)); err != nil {
		return err
	}

	d.SetId(policyListMap[0]["name"].(string))
	SetSingularResourceDataFromMapSchema(dataSourceNcloudAutoScalingPolicy(), d, policyListMap[0])
	return nil
}

func getAutoScalingPolicyList(config *ProviderConfig, id string) ([]*AutoScalingPolicy, error) {
	if config.SupportVPC {
		return getVpcAutoScalingPolicyList(config, id)
	} else {
		return getClassicAutoScalingPolicyList(config, id)
	}
}

func getVpcAutoScalingPolicyList(config *ProviderConfig, id string) ([]*AutoScalingPolicy, error) {
	reqParams := &vautoscaling.GetAutoScalingPolicyListRequest{
		RegionCode: &config.RegionCode,
	}

	resp, err := config.Client.vautoscaling.V2Api.GetAutoScalingPolicyList(reqParams)
	if err != nil {
		return nil, err
	}

	list := make([]*AutoScalingPolicy, 0)
	for _, p := range resp.ScalingPolicyList {
		policy := &AutoScalingPolicy{
			AutoScalingPolicyNo: p.PolicyNo,
			AutoScalingGroupNo:  p.AutoScalingGroupNo,
			AdjustmentTypeCode:  p.AdjustmentType.Code,
			ScalingAdjustment:   p.ScalingAdjustment,
			Cooldown:            p.CoolDown,
			MinAdjustmentStep:   p.MinAdjustmentStep,
		}
		if *p.PolicyName == id {
			return []*AutoScalingPolicy{policy}, nil
		}
		list = append(list, policy)
	}

	if id == "" {
		return nil, nil
	}
	return list, nil
}

func getClassicAutoScalingPolicyList(config *ProviderConfig, id string) ([]*AutoScalingPolicy, error) {
	reqParams := &autoscaling.GetAutoScalingPolicyListRequest{}

	if id != "" {
		reqParams.PolicyNameList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.autoscaling.V2Api.GetAutoScalingPolicyList(reqParams)
	if err != nil {
		return nil, err
	}

	list := make([]*AutoScalingPolicy, 0)
	for _, p := range resp.ScalingPolicyList {
		asg, err := getClassicAutoScalingGroupByName(config, *p.AutoScalingGroupName)
		if err != nil {
			return nil, err
		}
		policy := &AutoScalingPolicy{
			AutoScalingPolicyName: p.PolicyName,
			AdjustmentTypeCode:    p.AdjustmentType.Code,
			ScalingAdjustment:     p.ScalingAdjustment,
			Cooldown:              p.Cooldown,
			MinAdjustmentStep:     p.MinAdjustmentStep,
			AutoScalingGroupNo:    asg.AutoScalingGroupNo,
		}
		list = append(list, policy)
	}

	if len(list) < 1 {
		return nil, nil
	}

	return list, nil
}
