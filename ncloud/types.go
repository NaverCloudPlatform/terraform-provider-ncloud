package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/autoscaling"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vautoscaling"
)

type AutoScalingGroup struct {
	AutoScalingGroupNo                   *string   `json:"auto_scaling_group_no,omitempty"`
	AutoScalingGroupName                 *string   `json:"name,omitempty"`
	LaunchConfigurationNo                *string   `json:"launch_configuration_no,omitempty"`
	DesiredCapacity                      *int32    `json:"desired_capacity,omitempty"`
	MinSize                              *int32    `json:"min_size,omitempty"`
	MaxSize                              *int32    `json:"max_size,omitempty"`
	DefaultCooldown                      *int32    `json:"default_cooldown,omitempty"`
	LoadBalancerInstanceSummaryList      []*string `json:"loadBalancerInstanceSummaryList,omitempty"` // Check
	HealthCheckGracePeriod               *int32    `json:"health_check_grace_period,omitempty"`
	HealthCheckTypeCode                  *string   `json:"health_check_type_code,omitempty"`
	InAutoScalingGroupServerInstanceList []*string `json:"server_instance_no_list,omitempty"`
	SuspendedProcessList                 []*string `json:"suspendedProcessList,omitempty"` // Check
	ZoneList                             []*string `json:"zone_no_list,omitempty"`

	VpcNo                    *string   `json:"vpc_no,omitempty"`
	SubnetNo                 *string   `json:"subnet_no,omitempty"`
	ServerNamePrefix         *string   `json:"server_name_prefix,omitempty"`
	TargetGroupNoList        []*string `json:"target_group_list,omitempty"`
	AccessControlGroupNoList []*string `json:"access_control_group_no_list,omitempty"`
}

type InAutoScalingGroupServerInstance struct {
	HealthStatus        *string `json:"health_status,omitempty"`
	LifecycleState      *string `json:"lifecycle_state,omitempty"`
	LaunchConfiguration *string `json:"launch_configuration,omitempty"`
	ServerInstanceNo    *string `json:"server_instance_no,omitempty"`
	ServerInstanceName  *string `json:"server_instance_name,omitempty"`
}

func flattenZoneList(zoneList []*autoscaling.Zone) []*string {
	noList := make([]*string, 0)
	for _, z := range zoneList {
		noList = append(noList, z.ZoneNo)
	}
	return noList
}

func flattenClassicSuspendedProcessList(suspendedProcessList []*autoscaling.SuspendedProcess) []*string {
	codeList := make([]*string, 0)
	for _, p := range suspendedProcessList {
		codeList = append(codeList, p.Process.Code)
	}
	return codeList
}

func flattenVpcSuspendedProcessList(suspendedProcessList []*vautoscaling.SuspendedProcess) []*string {
	codeList := make([]*string, 0)
	for _, p := range suspendedProcessList {
		codeList = append(codeList, p.Process.Code)
	}
	return codeList
}

func flattenClassicAutoScalingGroupServerInstanceList(sl []*autoscaling.InAutoScalingGroupServerInstance) []*string {
	l := make([]*string, 0)
	for _, s := range sl {
		l = append(l, s.ServerInstanceNo)
	}
	return l
}

func flattenVpcAutoScalingGroupServerInstanceList(sl []*vautoscaling.InAutoScalingGroupServerInstance) []*string {
	l := make([]*string, 0)
	for _, s := range sl {
		l = append(l, s.ServerInstanceNo)
	}
	return l
}

func flattenLoadBalancerInstanceSummaryList(lbs []*autoscaling.LoadBalancerInstanceSummary) []*string {
	noList := make([]*string, 0)
	for _, lb := range lbs {
		noList = append(noList, lb.LoadBalancerInstanceNo)
	}
	return noList
}
