package ncloud

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func init() {
	RegisterResource("ncloud_lb_listener", resourceNcloudLbListener())
}

const (
	LoadBalancerListenerBusyStateErrorCode = "1200004"
	LoadBalancerListenerServerErrorCode    = "1250000"
)

func resourceNcloudLbListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudLbListenerCreate,
		ReadContext:   resourceNcloudLbListenerRead,
		UpdateContext: resourceNcloudLbListenerUpdate,
		DeleteContext: resourceNcloudLbListenerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"listener_no": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_balancer_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_group_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 65534)),
			},
			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS", "TCP", "TLS"}, false)),
			},
			"tls_min_version_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"TLSV10", "TLSV11", "TLSV12"}, false)),
			},
			"use_http2": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ssl_certificate_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_no_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNcloudLbListenerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_listener`"))
	}

	reqParams := &vloadbalancer.CreateLoadBalancerListenerRequest{
		RegionCode: &config.RegionCode,
		// Required
		LoadBalancerInstanceNo: ncloud.String(d.Get("load_balancer_no").(string)),
		TargetGroupNo:          ncloud.String(d.Get("target_group_no").(string)),
		Port:                   ncloud.Int32(int32(d.Get("port").(int))),
		ProtocolTypeCode:       ncloud.String(d.Get("protocol").(string)),

		// Optional
		SslCertificateNo:      StringPtrOrNil(d.GetOk("ssl_certificate_no")),
		UseHttp2:              BoolPtrOrNil(d.GetOk("use_http2")),
		TlsMinVersionTypeCode: StringPtrOrNil(d.GetOk("tls_min_version_type")),
	}

	listener := &vloadbalancer.LoadBalancerListener{}
	err := resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		resp, err := config.Client.vloadbalancer.V2Api.CreateLoadBalancerListener(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == LoadBalancerListenerBusyStateErrorCode || errBody.ReturnCode == LoadBalancerListenerServerErrorCode {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		listener = getListenerFromCreateResponseByPort(resp.LoadBalancerListenerList, reqParams.Port)
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ncloud.StringValue(listener.LoadBalancerListenerNo))
	return resourceNcloudLbListenerRead(ctx, d, meta)
}

func resourceNcloudLbListenerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_listener`"))
	}
	reqParams := &vloadbalancer.GetLoadBalancerListenerListRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(d.Get("load_balancer_no").(string)),
	}
	resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerListenerList(reqParams)
	if err != nil {
		return diag.FromErr(err)
	}

	listener := &LoadBalancerListener{}
	for _, l := range resp.LoadBalancerListenerList {
		if d.Id() == *l.LoadBalancerListenerNo {
			listener.LoadBalancerListenerNo = l.LoadBalancerListenerNo
			listener.ProtocolType = l.ProtocolType.Code
			listener.Port = l.Port
			listener.UseHttp2 = l.UseHttp2
			listener.SslCertificateNo = l.SslCertificateNo
			listener.TlsMinVersionType = l.TlsMinVersionType.Code
			listener.LoadBalancerRuleNoList = l.LoadBalancerRuleNoList
		}
	}

	listerMap := ConvertToMap(listener)
	SetSingularResourceDataFromMapSchema(resourceNcloudLbListener(), d, listerMap)
	return nil
}

func resourceNcloudLbListenerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_listener`"))
	}
	reqParams := &vloadbalancer.ChangeLoadBalancerListenerConfigurationRequest{
		RegionCode: &config.RegionCode,
		// Required
		LoadBalancerListenerNo: ncloud.String(d.Id()),
		Port:                   ncloud.Int32(int32(d.Get("port").(int))),
		ProtocolTypeCode:       ncloud.String(d.Get("protocol").(string)),

		// Optional
		SslCertificateNo:      StringPtrOrNil(d.GetOk("ssl_certificate_no")),
		UseHttp2:              BoolPtrOrNil(d.GetOk("use_http2")),
		TlsMinVersionTypeCode: StringPtrOrNil(d.GetOk("tls_min_version_type")),
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := config.Client.vloadbalancer.V2Api.ChangeLoadBalancerListenerConfiguration(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == LoadBalancerListenerBusyStateErrorCode || errBody.ReturnCode == LoadBalancerListenerServerErrorCode {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNcloudLbListenerRead(ctx, d, config)
}

func resourceNcloudLbListenerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_listener`"))
	}
	reqParams := &vloadbalancer.DeleteLoadBalancerListenersRequest{
		RegionCode:                 &config.RegionCode,
		LoadBalancerListenerNoList: []*string{ncloud.String(d.Id())},
	}

	err := resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		_, err := config.Client.vloadbalancer.V2Api.DeleteLoadBalancerListeners(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == LoadBalancerListenerBusyStateErrorCode || errBody.ReturnCode == LoadBalancerListenerServerErrorCode {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getListenerFromCreateResponseByPort(listenerList []*vloadbalancer.LoadBalancerListener, port *int32) *vloadbalancer.LoadBalancerListener {
	for _, listener := range listenerList {
		if *listener.Port == *port {
			return listener
		}
	}
	return nil
}
