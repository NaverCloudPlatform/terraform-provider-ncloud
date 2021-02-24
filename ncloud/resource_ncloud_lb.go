package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"reflect"
	"sort"
	"strconv"
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_group_no": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"listener_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS", "TCP", "TLS"}, false)),
						},
						"port": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 65534)),
						},
						"use_http2": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"ssl_certificate_no": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tls_min_version_type": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"TLSV10", "TLSV11", "TLSV12"}, false)),
						},
						"rule_no_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
		LoadBalancerListenerList: expandLoadBalancerListenerList(d.Get("listener_list").([]interface{})),
	}
	subnet, err := getSubnetInstance(config, *reqParams.SubnetNoList[0])
	if err != nil {
		return err
	}

	reqParams.VpcNo = subnet.VpcNo
	resp, err := config.Client.vloadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		return err
	}
	if err := waitForLoadBalancerActive(d, config, resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo); err != nil {
		return err
	}
	d.SetId(ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
	return resourceNcloudLbRead(d, meta)
}

func resourceNcloudLbRead(d *schema.ResourceData, meta interface{}) error {

	//It is used to maintain the diff Sync of target_group_no
	dataListenerList := d.Get("listener_list").([]interface{})
	targetGroupMap := make(map[int32]string, 0)
	for _, listener := range dataListenerList {
		listenerMap := listener.(map[string]interface{})
		targetGroupMap[int32(listenerMap["port"].(int))] = listenerMap["target_group_no"].(string)
	}

	// received Load Balancer instance detail from API
	config := meta.(*ProviderConfig)
	reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(d.Id()),
	}
	resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
	if err != nil {
		return err
	}
	lb := convertLbInstance(resp.LoadBalancerInstanceList[0])

	// receive all Listener belonging to Load Balancer from API
	listenerReqParams := &vloadbalancer.GetLoadBalancerListenerListRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: lb.LoadBalancerInstanceNo,
	}

	logCommonRequest("resourceNcloudLbRead Listener", listenerReqParams)
	listenerResp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerListenerList(listenerReqParams)
	if err != nil {
		logErrorResponse("resourceNcloudLbRead Listener", err, reqParams)
		return err
	}
	logResponse("resourceNcloudLbRead Listener", listenerResp)

	listenerList := make([]*LoadBalancerListener, 0)
	sort.Slice(listenerResp.LoadBalancerListenerList, func(i, j int) bool {
		a, _ := strconv.Atoi(*listenerResp.LoadBalancerListenerList[i].LoadBalancerListenerNo)
		b, _ := strconv.Atoi(*listenerResp.LoadBalancerListenerList[j].LoadBalancerListenerNo)
		return a < b
	})
	for _, respListener := range listenerResp.LoadBalancerListenerList {
		listener := convertListener(respListener)
		listener.TargetGroupNo = ncloud.String(targetGroupMap[ncloud.Int32Value(listener.Port)])
		listenerList = append(listenerList, listener)
	}

	sort.Slice(listenerList, func(i, j int) bool {
		a, _ := strconv.Atoi(*listenerList[i].LoadBalancerListenerNo)
		b, _ := strconv.Atoi(*listenerList[j].LoadBalancerListenerNo)
		return a < b
	})
	lb.LoadBalancerListenerList = listenerList
	lbMap := ConvertToMap(lb)
	SetSingularResourceDataFromMapSchema(resourceNcloudLb(), d, lbMap)
	return nil
}

func resourceNcloudLbUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	if d.HasChanges("idle_timeout", "throughput_type") {
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

	if d.HasChange("listener_list") {
		listenerList := d.Get("listener_list").([]interface{})

		o, n := d.GetChange("listener_list")

		oldListenerNoMap := make(map[string]bool)
		ol := o.([]interface{})
		for _, old := range ol {
			oldMap := old.(map[string]interface{})
			oldListenerNoMap[oldMap["listener_no"].(string)] = true
		}

		nl := n.([]interface{})
		newListenerNoMap := make(map[string]bool)
		for _, newInterface := range nl {
			newMap := newInterface.(map[string]interface{})
			newListenerNoMap[newMap["listener_no"].(string)] = true
		}

		for k, _ := range oldListenerNoMap {
			if newListenerNoMap[k] {
				delete(oldListenerNoMap, k)
			}
		}

		for k, v := range oldListenerNoMap {
			log.Printf("Removed listener (%s)", k)
			if v {
				if err := waitForLoadBalancerActive(d, config, ncloud.String(d.Id())); err != nil {
					return err
				}
				reqParams := &vloadbalancer.DeleteLoadBalancerListenersRequest{
					RegionCode:                 &config.RegionCode,
					LoadBalancerListenerNoList: []*string{ncloud.String(k)},
				}
				if _, err := config.Client.vloadbalancer.V2Api.DeleteLoadBalancerListeners(reqParams); err != nil {
					return err
				}
			}
		}

		for _, listener := range listenerList {
			listenerMap := listener.(map[string]interface{})
			if err := waitForLoadBalancerActive(d, config, ncloud.String(d.Id())); err != nil {
				return err
			}
			// non listener_no is new
			if reflect.Value.IsZero(reflect.ValueOf(listenerMap["listener_no"])) {
				reqParams := &vloadbalancer.CreateLoadBalancerListenerRequest{
					RegionCode:             &config.RegionCode,
					LoadBalancerInstanceNo: ncloud.String(d.Id()),
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["target_group_no"])); !ok {
					reqParams.TargetGroupNo = ncloud.String(listenerMap["target_group_no"].(string))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["protocol"])); !ok {
					reqParams.ProtocolTypeCode = ncloud.String(listenerMap["protocol"].(string))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["port"])); !ok {
					reqParams.Port = ncloud.Int32(int32(listenerMap["port"].(int)))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["ssl_certificate_no"])); !ok {
					reqParams.SslCertificateNo = ncloud.String(listenerMap["ssl_certificate_no"].(string))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["use_http2"])); !ok {
					reqParams.UseHttp2 = ncloud.Bool(listenerMap["use_http2"].(bool))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["tls_min_version_type"])); !ok {
					reqParams.TlsMinVersionTypeCode = ncloud.String(listenerMap["tls_min_version_type"].(string))
				}

				if _, err := config.Client.vloadbalancer.V2Api.CreateLoadBalancerListener(reqParams); err != nil {
					return err
				}
			} else {
				reqParams := &vloadbalancer.ChangeLoadBalancerListenerConfigurationRequest{
					RegionCode:             &config.RegionCode,
					LoadBalancerListenerNo: ncloud.String(listenerMap["listener_no"].(string)),
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["protocol"])); !ok {
					reqParams.ProtocolTypeCode = ncloud.String(listenerMap["protocol"].(string))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["port"])); !ok {
					reqParams.Port = ncloud.Int32(int32(listenerMap["port"].(int)))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["ssl_certificate_no"])); !ok {
					reqParams.SslCertificateNo = ncloud.String(listenerMap["ssl_certificate_no"].(string))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["use_http2"])); !ok {
					reqParams.UseHttp2 = ncloud.Bool(listenerMap["use_http2"].(bool))
				}

				if ok := reflect.Value.IsZero(reflect.ValueOf(listenerMap["tls_min_version_type"])); !ok {
					reqParams.TlsMinVersionTypeCode = ncloud.String(listenerMap["tls_min_version_type"].(string))
				}

				if _, err := config.Client.vloadbalancer.V2Api.ChangeLoadBalancerListenerConfiguration(reqParams); err != nil {
					return err
				}
			}
		}
	}
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
		Pending: []string{LoadBalancerInstanceStatusNameCreating, LoadBalancerInstanceStatusNameChanging},
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

func expandLoadBalancerListenerList(listeners []interface{}) []*vloadbalancer.LoadBalancerListenerParameter {
	listenerParameterList := make([]*vloadbalancer.LoadBalancerListenerParameter, 0)

	for _, listenerParameter := range listeners {
		parameterMap := listenerParameter.(map[string]interface{})

		parameter := &vloadbalancer.LoadBalancerListenerParameter{
			TargetGroupNo: ncloud.String(parameterMap["target_group_no"].(string)),
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["protocol"])); !ok {
			parameter.ProtocolTypeCode = ncloud.String(parameterMap["protocol"].(string))
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["port"])); !ok {
			parameter.Port = ncloud.Int32(int32(parameterMap["port"].(int)))
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["ssl_certificate_no"])); !ok {
			parameter.SslCertificateNo = ncloud.String(parameterMap["ssl_certificate_no"].(string))
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["use_http2"])); !ok {
			parameter.UseHttp2 = ncloud.Bool(parameterMap["use_http2"].(bool))
		}

		if ok := reflect.Value.IsZero(reflect.ValueOf(parameterMap["tls_min_version_type"])); !ok {
			parameter.TlsMinVersionTypeCode = ncloud.String(parameterMap["tls_min_version_type"].(string))
		}

		listenerParameterList = append(listenerParameterList, parameter)
	}

	return listenerParameterList
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
