package loadbalancer

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
