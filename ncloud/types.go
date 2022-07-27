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

type AutoScalingPolicy struct {
	AutoScalingPolicyNo   *string `json:"auto_scaling_policy_no,omitempty"`
	AutoScalingPolicyName *string `json:"name,omitempty"`
	AutoScalingGroupNo    *string `json:"auto_scaling_group_no,omitempty"`
	AdjustmentTypeCode    *string `json:"adjustment_type_code,omitempty"`
	ScalingAdjustment     *int32  `json:"scaling_adjustment,omitempty"`
	Cooldown              *int32  `json:"cooldown,omitempty"`
	MinAdjustmentStep     *int32  `json:"min_adjustment_step,omitempty"`
}

type AutoScalingSchedule struct {
	ScheduledActionNo   *string `json:"auto_scaling_schedule_no,omitempty"`
	ScheduledActionName *string `json:"name,omitempty"`
	AutoScalingGroupNo  *string `json:"auto_scaling_group_no,omitempty"`
	DesiredCapacity     *int32  `json:"desired_capacity,omitempty"`
	MinSize             *int32  `json:"min_size,omitempty"`
	MaxSize             *int32  `json:"max_size,omitempty"`
	StartTime           *string `json:"start_time,omitempty"`
	EndTime             *string `json:"end_time,omitempty"`
	RecurrenceInKST     *string `json:"recurrence,omitempty"`
	TimeZone            *string `json:"time_zone,omitempty"`
}

type InAutoScalingGroupServerInstance struct {
	HealthStatus        *string `json:"health_status,omitempty"`
	LifecycleState      *string `json:"lifecycle_state,omitempty"`
	LaunchConfiguration *string `json:"launch_configuration,omitempty"`
	ServerInstanceNo    *string `json:"server_instance_no,omitempty"`
	ServerInstanceName  *string `json:"server_instance_name,omitempty"`
}

type TargetGroup struct {
	TargetGroupNo           *string        `json:"target_group_no,omitempty"`
	TargetGroupName         *string        `json:"name,omitempty"`
	TargetType              *string        `json:"target_type,omitempty"`
	VpcNo                   *string        `json:"vpc_no,omitempty"`
	TargetGroupProtocolType *string        `json:"protocol,omitempty"`
	TargetGroupPort         *int32         `json:"port,omitempty"`
	TargetGroupDescription  *string        `json:"description,omitempty"`
	UseStickySession        *bool          `json:"use_sticky_session,omitempty"`
	UseProxyProtocol        *bool          `json:"use_proxy_protocol,omitempty"`
	AlgorithmType           *string        `json:"algorithm_type,omitempty"`
	LoadBalancerInstanceNo  *string        `json:"load_balancer_instance_no,omitempty"`
	TargetNoList            []*string      `json:"target_no_list"`
	HealthCheck             []*HealthCheck `json:"health_check"`
}

type HealthCheck struct {
	HealthCheckProtocolType   *string `json:"protocol,omitempty"`
	HealthCheckPort           *int32  `json:"port,omitempty"`
	HealthCheckUrlPath        *string `json:"url_path,omitempty"`
	HealthCheckHttpMethodType *string `json:"http_method,omitempty"`
	HealthCheckCycle          *int32  `json:"cycle,omitempty"`
	HealthCheckUpThreshold    *int32  `json:"up_threshold,omitempty"`
	HealthCheckDownThreshold  *int32  `json:"down_threshold,omitempty"`
}

type LoadBalancerInstance struct {
	LoadBalancerInstanceNo   *string   `json:"load_balancer_no,omitempty"`
	LoadBalancerDescription  *string   `json:"description,omitempty"`
	LoadBalancerName         *string   `json:"name,omitempty"`
	LoadBalancerDomain       *string   `json:"domain,omitempty"`
	LoadBalancerIpList       []*string `json:"ip_list,omitempty"`
	LoadBalancerType         *string   `json:"type,omitempty"`
	LoadBalancerNetworkType  *string   `json:"network_type,omitempty"`
	ThroughputType           *string   `json:"throughput_type,omitempty"`
	IdleTimeout              *int32    `json:"idle_timeout,omitempty"`
	VpcNo                    *string   `json:"vpc_no,omitempty"`
	SubnetNoList             []*string `json:"subnet_no_list,omitempty"`
	LoadBalancerListenerList []*string `json:"listener_no_list"`
}

type LoadBalancerListener struct {
	LoadBalancerListenerNo *string   `json:"listener_no,omitempty"`
	ProtocolType           *string   `json:"protocol,omitempty"`
	Port                   *int32    `json:"port,omitempty"`
	UseHttp2               *bool     `json:"use_http2,omitempty"`
	SslCertificateNo       *string   `json:"ssl_certificate_no,omitempty"`
	TlsMinVersionType      *string   `json:"tls_min_version_type,omitempty"`
	LoadBalancerRuleNoList []*string `json:"rule_no_list"`
	TargetGroupNo          *string   `json:"target_group_no,omitempty"`
}

type LoginKey struct {
	KeyName     *string `json:"key_name,omitempty"`
	Fingerprint *string `json:"fingerprint,omitempty"`
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

func flattenAccessControlGroupList(asgs []*autoscaling.AccessControlGroup) []*string {
	l := make([]*string, 0)
	for _, asg := range asgs {
		l = append(l, asg.AccessControlGroupConfigurationNo)
	}
	return l
}
