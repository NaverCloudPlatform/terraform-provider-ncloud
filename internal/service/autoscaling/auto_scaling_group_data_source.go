package autoscaling

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

func DataSourceNcloudAutoScalingGroup() *schema.Resource {
	fieldMap := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"filter": DataSourceFiltersSchema(),
	}
	return GetSingularDataSourceItemSchema(ResourceNcloudAutoScalingGroup(), fieldMap, dataSourceNcloudAutoScalingGroupRead)
}

func dataSourceNcloudAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	autoScalingGroupList, err := getAutoScalingGroupList(config, d.Id())

	if err != nil {
		return err
	}

	autoScalingGroupListMap := ConvertToArrayMap(autoScalingGroupList)
	if f, ok := d.GetOk("filter"); ok {
		autoScalingGroupListMap = ApplyFilters(f.(*schema.Set), autoScalingGroupListMap, DataSourceNcloudAutoScalingGroup().Schema)
	}

	if err := ValidateOneResult(len(autoScalingGroupListMap)); err != nil {
		return err
	}

	d.SetId(autoScalingGroupListMap[0]["auto_scaling_group_no"].(string))
	SetSingularResourceDataFromMapSchema(DataSourceNcloudAutoScalingGroup(), d, autoScalingGroupListMap[0])
	return nil
}

func getAutoScalingGroupList(config *conn.ProviderConfig, id string) ([]*AutoScalingGroup, error) {
	if config.SupportVPC {
		return getVpcAutoScalingGroupList(config, id)
	} else {
		return getClassicAutoScalingGroupList(config, id)
	}

}

func getVpcAutoScalingGroupList(config *conn.ProviderConfig, id string) ([]*AutoScalingGroup, error) {
	reqParams := &vautoscaling.GetAutoScalingGroupListRequest{
		RegionCode: &config.RegionCode,
	}

	if id != "" {
		reqParams.AutoScalingGroupNoList = []*string{ncloud.String(id)}
	}

	resp, err := config.Client.Vautoscaling.V2Api.GetAutoScalingGroupList(reqParams)
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

func getClassicAutoScalingGroupList(config *conn.ProviderConfig, id string) ([]*AutoScalingGroup, error) {
	no := ncloud.String(id)
	reqParams := &autoscaling.GetAutoScalingGroupListRequest{
		RegionNo: &config.RegionNo,
	}

	resp, err := config.Client.Autoscaling.V2Api.GetAutoScalingGroupList(reqParams)
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
