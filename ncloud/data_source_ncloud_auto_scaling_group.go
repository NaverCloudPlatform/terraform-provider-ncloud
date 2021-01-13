package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	RegisterDataSource("ncloud_auto_scaling_group", dataSourceNcloudAutoScalingGroup())
}

func dataSourceNcloudAutoScalingGroup() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"auto_scaling_group_no": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": dataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchema(resourceNcloudAutoScalingGroup(), fieldMap, dataSourceNcloudAUtoScalingGroupRead)
}

func dataSourceNcloudAUtoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if v, ok := d.GetOk("auto_scaling_group_no"); ok {
		d.SetId(v.(string))
	}

	autoScalingGroupList, err := getAutoScalingGroupList(config, d.Id())

	if err != nil {
		return err
	}

	autoScalingGroupListMap := ConvertToArrayMap(autoScalingGroupList)
	if f, ok := d.GetOk("filter"); ok {
		autoScalingGroupListMap = ApplyFilters(f.(*schema.Set), autoScalingGroupListMap, dataSourceNcloudAutoScalingGroup().Schema)
	}

	if err := validateOneResult(len(autoScalingGroupListMap)); err != nil {
		return err
	}

	d.SetId(autoScalingGroupListMap[0]["auto_scaling_group_no"].(string))
	SetSingularResourceDataFromMapSchema(dataSourceNcloudAutoScalingGroup(), d, autoScalingGroupListMap[0])
	return nil
}

func getAutoScalingGroupList(config *ProviderConfig, id string) ([]*AutoScalingGroup, error) {
	if config.SupportVPC {
		return getVpcAutoScalingGroupList(config, id)
	} else {
		return getClassicAutoScalingGroupList(config, id)
	}

}

func getVpcAutoScalingGroupList(config *ProviderConfig, id string) ([]*AutoScalingGroup, error) {
	reqParams := &vautoscaling.GetAutoScalingGroupListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.AutoScalingGroupNoList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.vautoscaling.V2Api.GetAutoScalingGroupList(reqParams)
	if err != nil {
		return nil, err
	}

	if len(resp.AutoScalingGroupList) < 1 {
		return nil, nil
	}

	list := make([]*AutoScalingGroup, 0)
	for _, a := range resp.AutoScalingGroupList {
		list = append(list, &AutoScalingGroup{
			AutoScalingGroupNo:                   a.AutoScalingGroupNo,
			AutoScalingGroupName:                 a.AutoScalingGroupName,
			LaunchConfigurationNo:                a.LaunchConfigurationNo,
			DesiredCapacity:                      a.DesiredCapacity,
			MinSize:                              a.MinSize,
			MaxSize:                              a.MaxSize,
			DefaultCooldown:                      a.DefaultCoolDown,
			HealthCheckGracePeriod:               a.HealthCheckGracePeriod,
			HealthCheckTypeCode:                  a.HealthCheckType.Code,
			InAutoScalingGroupServerInstanceList: flattenVpcAutoScalingGroupServerInstanceList(a.InAutoScalingGroupServerInstanceList),
			SuspendedProcessList:                 flattenVpcSuspendedProcessList(a.SuspendedProcessList),
			VpcNo:                                a.VpcNo,
			SubnetNo:                             a.SubnetNo,
			ServerNamePrefix:                     a.ServerNamePrefix,
			TargetGroupNoList:                    a.TargetGroupNoList,
			AccessControlGroupNoList:             a.AccessControlGroupNoList,
		})
	}

	return list, nil
}

func getClassicAutoScalingGroupList(config *ProviderConfig, id string) ([]*AutoScalingGroup, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetAutoScalingGroupListRequest{
		RegionNo: &config.RegionNo,
	}

	resp, err := config.Client.autoscaling.V2Api.GetAutoScalingGroupList(reqParams)
	if err != nil {
		return nil, err
	}

	list := make([]*AutoScalingGroup, 0)
	for _, a := range resp.AutoScalingGroupList {
		autoScalingGroup := &AutoScalingGroup{
			AutoScalingGroupNo:                   a.AutoScalingGroupNo,
			AutoScalingGroupName:                 a.AutoScalingGroupName,
			LaunchConfigurationNo:                a.LaunchConfigurationNo,
			DesiredCapacity:                      a.DesiredCapacity,
			MinSize:                              a.MinSize,
			MaxSize:                              a.MaxSize,
			DefaultCooldown:                      a.DefaultCooldown,
			LoadBalancerInstanceSummaryList:      flattenLoadBalancerInstanceSummaryList(a.LoadBalancerInstanceSummaryList),
			HealthCheckGracePeriod:               a.HealthCheckGracePeriod,
			HealthCheckTypeCode:                  a.HealthCheckType.Code,
			InAutoScalingGroupServerInstanceList: flattenClassicAutoScalingGroupServerInstanceList(a.InAutoScalingGroupServerInstanceList),
			SuspendedProcessList:                 flattenClassicSuspendedProcessList(a.SuspendedProcessList),
			ZoneList:                             flattenZoneList(a.ZoneList),
		}
		if *a.AutoScalingGroupNo == *no {
			return []*AutoScalingGroup{autoScalingGroup}, nil
		}
		list = append(list, autoScalingGroup)
	}

	if *no != "" {
		return nil, nil
	}

	return list, nil
}
