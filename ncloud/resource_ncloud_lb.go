package ncloud

import (
	"bytes"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"reflect"
	"time"
)

func init() {
	RegisterResource("ncloud_lb", resourceNcloudLb())
}

const (
	LoadBalancerInstanceStatusNameCreating    = "Creating"
	LoadBalancerInstanceStatusNameRunning     = "Running"
	LoadBalancerInstanceStatusNameChanging    = "Changing"
	LoadBalancerInstanceStatusNameTerminating = "Terminating"
	LoadBalancerInstanceStatusNameTerminated  = "Terminated"
	LoadBalancerInstanceStatusNameRepairing   = "Repairing"
)

func resourceNcloudLb() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLbCreate,
		Read:   resourceNcloudLbRead,
		Update: resourceNcloudLbUpdate,
		Delete: resourceNcloudLbDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false)),
			},
			"idle_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 3600)),
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"APPLICATION", "NETWORK", "NETWORK_PROXY"}, false)),
			},
			"throughput_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false)),
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_no_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_group_no": {
							//타겟 그룹 번호
							//선택한 타겟 그룹은 각 리스너의 DEFAULT 규칙에 적용됩니다.
							//다른 로드밸런서에서 이미 사용중인 타겟 그룹은 이용할 수 없습니다.
							//로드밸런서 유형과 타겟 그룹 프로토콜 유형에 따라서 사용 가능한 타겟 그룹이 제한됩니다.
							//NETWORK : TCP
							//NETWORK_PROXY : PROXY_TCP
							//APPLICATION : HTTP / HTTPS
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							//로드밸런서 리스너 프로토콜 유형 코드
							//최소 한개의 리스너를 등록해야 합니다.
							//로드밸런서 유형에 따라서 사용 가능한 리스너 프로토콜 유형과 기본값이 결정됩니다.
							//APPLICATION : HTTP (Default) / HTTPS
							//NETWORK : TCP (Default)
							//NETWORK_PROXY : TCP (Default) / TLS
							Type:             schema.TypeString,
							Optional:         true,
							Default:          80,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS", "TCP", "TLS"}, false)),
						},
						"port": {
							//로드밸런서 리스너 포트
							//리스너 프로토콜 유형에 따라서 포트 기본값이 결정됩니다.
							//Default :
							//HTTP / TCP : 80
							//HTTPS / TLS : 443
							//포트 번호는 중복될 수 없습니다.
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 65534)),
						},
						"use_http2": {
							//HTTP/2 프로토콜 사용 여부
							//Options : true | false
							//Default : false
							//리스너 프로토콜 유형이 HTTPS 인 경우에만 유효합니다.
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"ssl_certificate_no": {
							//SSL 인증서 번호
							//리스너 프로토콜 유형이 HTTPS 또는 TLS 인 경우 SSL 인증서를 반드시 설정해야 합니다.
							//sslCertificateNo는 GET https://certificatemanager.apigw.ntruss.com/api/v1/certificates 액션을 통해서 획득할 수 있습니다.
							Type:     schema.TypeString,
							Optional: true,
						},
						"tls_min_version_type": {
							//TLS 최소 지원 버전 유형 코드
							//리스너 프로토콜 유형이 HTTPS 또는 TLS 인 경우에만 유효합니다.
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"TLSV10", "TLSV11", "TLSV12"}, false)),
						},
					},
				},
			},
		},
	}
}

func resourceNcloudLbCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_lb_target_group`")
	}
	reqParams := &vloadbalancer.CreateLoadBalancerInstanceRequest{
		RegionCode: &config.RegionCode,
		// Optional
		IdleTimeout:                 Int32PtrOrNil(d.GetOk("idle_timeout")),
		LoadBalancerDescription:     StringPtrOrNil(d.GetOk("description")),
		LoadBalancerNetworkTypeCode: StringPtrOrNil(d.GetOk("network_type")),
		LoadBalancerName:            StringPtrOrNil(d.GetOk("name")),
		ThroughputTypeCode:          StringPtrOrNil(d.GetOk("throughput_type")),

		// Required
		LoadBalancerTypeCode:     ncloud.String(d.Get("type").(string)),
		SubnetNoList:             ncloud.StringInterfaceList(d.Get("subnet_no_list").([]interface{})),
		LoadBalancerListenerList: expandLoadBalancerListenerList(d.Get("listener_list").(*schema.Set)),
	}

	//var listenerParameterList []*vloadbalancer.LoadBalancerListenerParameter
	//for _, listener := range d.Get("listener_list").([]interface{}) {
	//	listenerMap := listener.(map[string]interface{})
	//
	//	fmt.Print(listenerMap)
	//	listenerParameter := &vloadbalancer.LoadBalancerListenerParameter{
	//		SslCertificateNo: ncloud.String(listenerMap["ssl_certificate_no"].(string)),
	//		UseHttp2:         ncloud.Bool(listenerMap["use_http2"].(bool)),
	//		TargetGroupNo:    ncloud.String(listenerMap["target_group_no"].(string)),
	//	}
	//
	//	if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["protocol"])); !ok {
	//		listenerParameter.ProtocolTypeCode = ncloud.String(listenerMap["protocol"].(string))
	//	}
	//
	//	if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["port"])); !ok {
	//		listenerParameter.Port = ncloud.Int32(int32(listenerMap["port"].(int)))
	//	}
	//
	//	if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["tls_min_version_type"])); !ok {
	//		listenerParameter.TlsMinVersionTypeCode = ncloud.String(listenerMap["tls_min_version_type"].(string))
	//	}
	//
	//	listenerParameterList = append(listenerParameterList, listenerParameter)
	//}
	//reqParams.LoadBalancerListenerList = listenerParameterList
	subnet, err := getSubnetInstance(config, *reqParams.SubnetNoList[0])
	if err != nil {
		return err
	}

	reqParams.VpcNo = subnet.VpcNo
	logCommonRequest("resourceNcloudLbCreate", reqParams)
	resp, err := config.Client.vloadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	logResponse("resourceNcloudLbCreate", resp)
	if err != nil {
		logErrorResponse("resourceNcloudLbCreate", err, reqParams)
		return err
	}
	if err := waitForLoadBalancerActive(d, config, resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo); err != nil {
		return err
	}
	d.SetId(ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
	return resourceNcloudLbRead(d, meta)
}

func resourceNcloudLbRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(d.Id()),
	}
	logCommonRequest("resourceNcloudLbRead", reqParams)
	resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
	logResponse("resourceNcloudLbRead", resp)
	if err != nil {
		logErrorResponse("resourceNcloudLbRead", err, reqParams)
		return err
	}

	lb := convertLbInstance(resp.LoadBalancerInstanceList[0])

	listenerReqParams := &vloadbalancer.GetLoadBalancerListenerListRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: lb.LoadBalancerInstanceNo,
	}
	logCommonRequest("resourceNcloudLbRead Listener", listenerReqParams)
	listenerResp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerListenerList(listenerReqParams)
	if err != nil {
		return err
	}
	logResponse("resourceNcloudLbRead Listener", listenerResp)

	listenerList := make([]*LoadBalancerListener, 0)
	for _, respListener := range listenerResp.LoadBalancerListenerList {
		listener := convertListener(respListener)
		listenerList = append(listenerList, listener)
	}

	lb.LoadBalancerListenerList = listenerList
	l := d.Get("listener_list").(*schema.Set)

	// Create empty set for getAccessControlGroupRuleList
	lSet := schema.NewSet(schema.HashResource(resourceNcloudLb().Schema["listener_list"].Elem.(*schema.Resource)), []interface{}{})

	for _, listener := range lb.LoadBalancerListenerList {
		m := map[string]interface{}{
			"listener_no":          ncloud.StringValue(listener.LoadBalancerListenerNo),
			"protocol":             ncloud.StringValue(listener.ProtocolType),
			"port":                 int(ncloud.Int32Value(listener.Port)),
			"use_http2":            ncloud.BoolValue(listener.UseHttp2),
			"ssl_certificate_no":   ncloud.StringValue(listener.SslCertificateNo),
			"tls_min_version_type": ncloud.StringValue(listener.TlsMinVersionType),
			//"rule_no_list":         ncloud.StringListValue(listener.LoadBalancerRuleNoList),
		}

		logResponse("yoogle-test-4", m)
		lSet.Add(m)
	}

	if err := d.Set("listener_list", l.Intersection(lSet).List()); err != nil {
		log.Printf("[WARN] Error setting outbound rule set for (%s): %s", d.Id(), err)
	}

	//lbMap := ConvertToMap(lb)
	//SetSingularResourceDataFromMapSchema(resourceNcloudLb(), d, lbMap)
	//if err := d.Set("listener_list", flattenLoadBalancerListenerList(lb.LoadBalancerListenerList)); err != nil {
	//	return err
	//}
	return nil
}

func convertListener(listener *vloadbalancer.LoadBalancerListener) *LoadBalancerListener {
	return &LoadBalancerListener{
		LoadBalancerListenerNo: listener.LoadBalancerListenerNo,
		ProtocolType:           listener.ProtocolType.Code,
		Port:                   listener.Port,
		UseHttp2:               listener.UseHttp2,
		SslCertificateNo:       listener.SslCertificateNo,
		TlsMinVersionType:      listener.TlsMinVersionType.Code,
		LoadBalancerRuleNoList: listener.LoadBalancerRuleNoList,
	}
}

func resourceNcloudLbUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChange("idle_timeout") || d.HasChange("throughput_type") {
		_, err := config.Client.vloadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(&vloadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
			RegionCode:             &config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(d.Id()),
			IdleTimeout:            Int32PtrOrNil(d.GetOk("idle_timeout")),
			ThroughputTypeCode:     StringPtrOrNil(d.GetOk("throughput_type")),
		})
		if err != nil {
			return err
		}
	}

	//if d.HasChange("listener_list") {
	//	_, err := config.Client.vloadbalancer.V2Api.ChangeLoadBalancerListenerConfiguration(&vloadbalancer.ChangeLoadBalancerListenerConfigurationRequest{
	//		RegionCode:             &config.RegionCode,
	//		SslCertificateNo:       nil,
	//		UseHttp2:               nil,
	//		LoadBalancerListenerNo: nil,
	//		Port:                   nil,
	//		ProtocolTypeCode:       nil,
	//		TlsMinVersionTypeCode:  nil,
	//	})
	//}
	return nil
}

func resourceNcloudLbDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	deleteInstanceReqParams := &vloadbalancer.DeleteLoadBalancerInstancesRequest{
		RegionCode:                 &config.RegionCode,
		LoadBalancerInstanceNoList: ncloud.StringList([]string{d.Id()}),
	}

	logCommonRequest("resourceNcloudLbDelete", deleteInstanceReqParams)
	if _, err := config.Client.vloadbalancer.V2Api.DeleteLoadBalancerInstances(deleteInstanceReqParams); err != nil {
		logErrorResponse("resourceNcloudLbDelete", err, deleteInstanceReqParams)
		return err
	}

	if err := waitForLoadBalancerDeletion(d, config); err != nil {
		return err
	}

	return nil
}

func waitForLoadBalancerDeletion(d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{LoadBalancerInstanceStatusNameTerminating},
		Target:  []string{LoadBalancerInstanceStatusNameTerminated},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(d.Id()),
			}
			resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
			if err != nil {
				return nil, "", nil
			}

			if len(resp.LoadBalancerInstanceList) < 1 {
				return resp, LoadBalancerInstanceStatusNameTerminated, nil
			}

			lb := resp.LoadBalancerInstanceList[0]
			return resp, ncloud.StringValue(lb.LoadBalancerInstanceStatusName), nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Load Balancer instance (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForLoadBalancerActive(d *schema.ResourceData, config *ProviderConfig, no *string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{LoadBalancerInstanceStatusNameCreating},
		Target:  []string{LoadBalancerInstanceStatusNameRunning},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: no,
			}
			resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
			if err != nil {
				return nil, "", nil
			}

			if len(resp.LoadBalancerInstanceList) < 1 {
				return nil, "", fmt.Errorf("Not found load balancer instance(%s)", *no)
			}

			lb := resp.LoadBalancerInstanceList[0]
			return resp, ncloud.StringValue(lb.LoadBalancerInstanceStatusName), nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Load Balancer instance (%s) to become activating: %s", *no, err)
	}
	return nil
}

func expandLoadBalancerListenerList(listeners *schema.Set) []*vloadbalancer.LoadBalancerListenerParameter {
	listenerParameterList := make([]*vloadbalancer.LoadBalancerListenerParameter, 0)

	for _, listenerParameter := range listeners.List() {
		parameterMap := listenerParameter.(map[string]interface{})

		parameter := &vloadbalancer.LoadBalancerListenerParameter{
			TargetGroupNo:    ncloud.String(parameterMap["target_group_no"].(string)),
			SslCertificateNo: ncloud.String(parameterMap["ssl_certificate_no"].(string)),
			UseHttp2:         ncloud.Bool(parameterMap["use_http2"].(bool)),
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["protocol"])); !ok {
			parameter.ProtocolTypeCode = ncloud.String(parameterMap["protocol"].(string))
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["port"])); !ok {
			parameter.Port = ncloud.Int32(int32(parameterMap["port"].(int)))
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["tls_min_version_type"])); !ok {
			parameter.TlsMinVersionTypeCode = ncloud.String(parameterMap["tls_min_version_type"].(string))
		}

		listenerParameterList = append(listenerParameterList, parameter)
	}

	return listenerParameterList
}

func flattenLoadBalancerListenerList(listenerList []*LoadBalancerListener) *schema.Set {
	vListener := make([]interface{}, 0)

	for _, listener := range listenerList {
		mListener := map[string]interface{}{
			"listener_no":          ncloud.StringValue(listener.LoadBalancerListenerNo),
			"protocol":             ncloud.StringValue(listener.ProtocolType),
			"port":                 int(ncloud.Int32Value(listener.Port)),
			"use_http2":            ncloud.BoolValue(listener.UseHttp2),
			"ssl_certificate_no":   ncloud.StringValue(listener.SslCertificateNo),
			"tls_min_version_type": ncloud.StringValue(listener.TlsMinVersionType),
			"rule_no_list":         ncloud.StringListValue(listener.LoadBalancerRuleNoList),
		}
		vListener = append(vListener, mListener)
	}

	logResponse("yoogle-test-2", vListener)
	return schema.NewSet(LbHash, vListener)
}

func LbHash(vLb interface{}) int {
	var buf bytes.Buffer

	mRb := vLb.(map[string]interface{})

	if v, ok := mRb["listener_no"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	if v, ok := mRb["target_group_no"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	if v, ok := mRb["protocol"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	if v, ok := mRb["port"].(int); ok {
		buf.WriteString(fmt.Sprintf("%d-", v))

	}

	if v, ok := mRb["use_http2"].(bool); ok {
		buf.WriteString(fmt.Sprintf("%t-", v))
	}

	if v, ok := mRb["ssl_certificate_no"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	if v, ok := mRb["tls_min_version_type"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	//if v, ok := mRb["use_http2"].(string); ok {
	//	buf.WriteString(fmt.Sprintf("%s-", v))
	//}

	return hashcode(buf.String())
}

func convertLbInstance(instance *vloadbalancer.LoadBalancerInstance) *LoadBalancerInstance {
	return &LoadBalancerInstance{
		LoadBalancerInstanceNo:         instance.LoadBalancerInstanceNo,
		LoadBalancerInstanceStatus:     instance.LoadBalancerInstanceStatus.Code,
		LoadBalancerInstanceOperation:  instance.LoadBalancerInstanceOperation.Code,
		LoadBalancerInstanceStatusName: instance.LoadBalancerInstanceStatusName,
		LoadBalancerDescription:        instance.LoadBalancerDescription,
		LoadBalancerName:               instance.LoadBalancerName,
		LoadBalancerDomain:             instance.LoadBalancerDomain,
		LoadBalancerIpList:             instance.LoadBalancerIpList,
		LoadBalancerType:               instance.LoadBalancerType.Code,
		LoadBalancerNetworkType:        instance.LoadBalancerNetworkType.Code,
		ThroughputType:                 instance.ThroughputType.Code,
		IdleTimeout:                    instance.IdleTimeout,
		VpcNo:                          instance.VpcNo,
		SubnetNoList:                   instance.SubnetNoList,
	}
}
