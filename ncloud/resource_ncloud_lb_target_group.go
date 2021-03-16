package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("ncloud_lb_target_group", resourceNcloudLbTargetGroup())
}

func resourceNcloudLbTargetGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudTargetGroupCreate,
		ReadContext:   resourceNcloudTargetGroupRead,
		UpdateContext: resourceNcloudTargetGroupUpdate,
		DeleteContext: resourceNcloudTargetGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"target_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringLenBetween(3, 30)),
			},
			"port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 65534)),
			},
			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"TCP", "PROXY_TCP", "HTTP", "HTTPS"}, false)),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cycle": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          30,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(5, 300)),
						},
						"down_threshold": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          2,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(2, 10)),
						},
						"up_threshold": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          2,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(2, 10)),
						},
						"http_method": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"HEAD", "GET"}, false)),
						},
						"port": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          80,
							ValidateDiagFunc: ToDiagFunc(validation.IntBetween(1, 65534)),
						},
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false)),
						},
						"url_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"target_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"VSVR"}, false)),
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"use_sticky_session": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"use_proxy_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"algorithm_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"RR", "SIPHS", "LC", "MH"}, false)),
			},
			"load_balancer_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudTargetGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group`"))
	}
	reqParams := &vloadbalancer.CreateTargetGroupRequest{
		RegionCode: &config.RegionCode,
		// Optional
		TargetGroupPort:        Int32PtrOrNil(d.GetOk("port")),
		TargetGroupDescription: StringPtrOrNil(d.GetOk("description")),
		TargetGroupName:        StringPtrOrNil(d.GetOk("name")),
		TargetNoList:           StringListPtrOrNil(d.GetOk("target_no_list")),
		// Required
		TargetTypeCode:              ncloud.String(d.Get("target_type").(string)),
		VpcNo:                       ncloud.String(d.Get("vpc_no").(string)),
		TargetGroupProtocolTypeCode: ncloud.String(d.Get("protocol").(string)),
	}

	if healthChecks, ok := d.GetOk("health_check"); ok {
		healthCheck := healthChecks.([]interface{})[0].(map[string]interface{})

		reqParams.HealthCheckCycle = ncloud.Int32(int32(healthCheck["cycle"].(int)))
		reqParams.HealthCheckDownThreshold = ncloud.Int32(int32(healthCheck["down_threshold"].(int)))
		reqParams.HealthCheckUpThreshold = ncloud.Int32(int32(healthCheck["up_threshold"].(int)))
		reqParams.HealthCheckPort = ncloud.Int32(int32(healthCheck["port"].(int)))

		// Required
		reqParams.HealthCheckProtocolTypeCode = ncloud.String(healthCheck["protocol"].(string))
		if err := validateHealthCheckProtocolByTargetGroupProtocol(*reqParams.TargetGroupProtocolTypeCode, *reqParams.HealthCheckProtocolTypeCode); err != nil {
			return diag.FromErr(err)
		}

		if *reqParams.HealthCheckProtocolTypeCode == "HTTP" || *reqParams.HealthCheckProtocolTypeCode == "HTTPS" {
			reqParams.HealthCheckUrlPath = ncloud.String(healthCheck["url_path"].(string))
			if healthCheck["http_method"] == "" {
				return diag.FromErr(fmt.Errorf("http_method is required if the health check protocol type is HTTP or HTTPS."))
			}
			reqParams.HealthCheckHttpMethodTypeCode = ncloud.String(healthCheck["http_method"].(string))
		}
	}

	logCommonRequest("resourceNcloudTargetGroupCreate", reqParams)
	resp, err := config.Client.vloadbalancer.V2Api.CreateTargetGroup(reqParams)
	logResponse("resourceNcloudTargetGroupCreate", resp)
	if err != nil {
		logErrorResponse("resourceNcloudTargetGroupCreate", err, reqParams)
		return diag.FromErr(err)
	}

	d.SetId(ncloud.StringValue(resp.TargetGroupList[0].TargetGroupNo))
	return resourceNcloudTargetGroupUpdate(ctx, d, meta)
}

func resourceNcloudTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group`"))
	}
	reqParams := &vloadbalancer.GetTargetGroupListRequest{
		RegionCode:        &config.RegionCode,
		TargetGroupNoList: []*string{ncloud.String(d.Id())},
	}
	logCommonRequest("resourceNcloudTargetGroupRead", reqParams)
	resp, err := config.Client.vloadbalancer.V2Api.GetTargetGroupList(reqParams)
	logResponse("resourceNcloudTargetGroupRead", resp)
	if err != nil {
		logErrorResponse("resourceNcloudTargetGroupRead", err, reqParams)
		return diag.FromErr(err)
	}

	if len(resp.TargetGroupList) < 1 {
		return nil
	}

	respTg := resp.TargetGroupList[0]
	tg := &TargetGroup{
		TargetGroupNo:           respTg.TargetGroupNo,
		TargetGroupName:         respTg.TargetGroupName,
		TargetType:              respTg.TargetType.Code,
		VpcNo:                   respTg.VpcNo,
		TargetGroupProtocolType: respTg.TargetGroupProtocolType.Code,
		TargetGroupPort:         respTg.TargetGroupPort,
		TargetGroupDescription:  respTg.TargetGroupDescription,
		UseStickySession:        respTg.UseStickySession,
		UseProxyProtocol:        respTg.UseProxyProtocol,
		AlgorithmType:           respTg.AlgorithmType.Code,
		LoadBalancerInstanceNo:  respTg.LoadBalancerInstanceNo,
		TargetNoList:            respTg.TargetNoList,
		HealthCheck: []*HealthCheck{
			{
				HealthCheckProtocolType:   respTg.HealthCheckProtocolType.Code,
				HealthCheckPort:           respTg.HealthCheckPort,
				HealthCheckUrlPath:        respTg.HealthCheckUrlPath,
				HealthCheckHttpMethodType: respTg.HealthCheckHttpMethodType.Code,
				HealthCheckCycle:          respTg.HealthCheckCycle,
				HealthCheckUpThreshold:    respTg.HealthCheckUpThreshold,
				HealthCheckDownThreshold:  respTg.HealthCheckDownThreshold,
			},
		},
	}
	tgMap := ConvertToMap(tg)
	SetSingularResourceDataFromMapSchema(resourceNcloudLbTargetGroup(), d, tgMap)
	return nil
}

func resourceNcloudTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group`"))
	}
	if d.HasChange("health_check") {
		reqParams := &vloadbalancer.ChangeTargetGroupHealthCheckConfigurationRequest{
			RegionCode:    &config.RegionCode,
			TargetGroupNo: ncloud.String(d.Id()),
		}
		healthChecks := d.Get("health_check").([]interface{})
		if len(healthChecks) == 1 {
			healthCheck := healthChecks[0].(map[string]interface{})
			healthCheckProtocol := healthCheck["protocol"].(string)

			reqParams.HealthCheckCycle = ncloud.Int32(int32(healthCheck["cycle"].(int)))
			reqParams.HealthCheckDownThreshold = ncloud.Int32(int32(healthCheck["down_threshold"].(int)))
			reqParams.HealthCheckPort = ncloud.Int32(int32(healthCheck["port"].(int)))
			reqParams.HealthCheckUpThreshold = ncloud.Int32(int32(healthCheck["up_threshold"].(int)))

			if healthCheckProtocol == "HTTP" || healthCheckProtocol == "HTTPS" {
				reqParams.HealthCheckUrlPath = ncloud.String(healthCheck["url_path"].(string))
				if healthCheck["http_method"] == "" {
					return diag.FromErr(fmt.Errorf("http_method is required if the health check protocol type is HTTP or HTTPS."))
				}
				reqParams.HealthCheckHttpMethodTypeCode = ncloud.String(healthCheck["http_method"].(string))
			}
		}
		logCommonRequest("resourceNcloudTargetGroupUpdate", reqParams)
		if _, err := config.Client.vloadbalancer.V2Api.ChangeTargetGroupHealthCheckConfiguration(reqParams); err != nil {
			logErrorResponse("resourceNcloudTargetGroupUpdate", err, reqParams)
			return diag.FromErr(err)
		}
	}

	if d.HasChange("algorithm_type") || d.HasChange("use_sticky_session") || d.HasChange("use_proxy_protocol") {
		reqParams := &vloadbalancer.ChangeTargetGroupConfigurationRequest{
			RegionCode:        &config.RegionCode,
			TargetGroupNo:     ncloud.String(d.Id()),
			AlgorithmTypeCode: ncloud.String(d.Get("algorithm_type").(string)),
		}

		targetGroupProtocol := d.Get("protocol").(string)
		switch targetGroupProtocol {
		case "HTTP", "HTTPS", "TCP":
			reqParams.UseStickySession = ncloud.Bool(d.Get("use_sticky_session").(bool))
		case "PROXY_TCP":
			reqParams.UseProxyProtocol = ncloud.Bool(d.Get("use_proxy_protocol").(bool))
		}

		if err := validateAlgorithmTypeByTargetGroupProtocol(*reqParams.AlgorithmTypeCode, targetGroupProtocol); err != nil {
			return diag.FromErr(err)
		}
		logCommonRequest("resourceNcloudTargetGroupUpdate", reqParams)
		if _, err := config.Client.vloadbalancer.V2Api.ChangeTargetGroupConfiguration(reqParams); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNcloudTargetGroupRead(ctx, d, meta)
}

func resourceNcloudTargetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb_target_group`"))
	}
	reqParams := &vloadbalancer.DeleteTargetGroupsRequest{
		RegionCode:        &config.RegionCode,
		TargetGroupNoList: []*string{ncloud.String(d.Id())},
	}
	if _, err := config.Client.vloadbalancer.V2Api.DeleteTargetGroups(reqParams); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func validateAlgorithmTypeByTargetGroupProtocol(algorithmType string, protocol string) error {
	protocolMap := make(map[string][]string)
	protocolMap["PROXY_TCP"] = []string{"RR", "SIPHS", "LC"}
	protocolMap["HTTP"] = []string{"RR", "SIPHS", "LC"}
	protocolMap["HTTPS"] = []string{"RR", "SIPHS", "LC"}
	protocolMap["TCP"] = []string{"MH", "RR"}
	if ok := containsInStringList(protocolMap[protocol], algorithmType); !ok {
		return fmt.Errorf("%s protocol is only suppoort %s algorithm types", protocol, protocolMap[protocol])
	}
	return nil
}

func validateHealthCheckProtocolByTargetGroupProtocol(targetGroupProtocol string, healthCheckProtocol string) error {
	if targetGroupProtocol == "TCP" || targetGroupProtocol == "PROXY_TCP" {
		if healthCheckProtocol != "TCP" {
			return fmt.Errorf("Health check protocol is only support TCP when target group protocol is %s.", targetGroupProtocol)
		}
	} else {
		if healthCheckProtocol != "HTTP" && healthCheckProtocol != "HTTPS" {
			return fmt.Errorf("Health check protocol is only support HTTP, HTTPS when target group protocol is %s.", targetGroupProtocol)
		}
	}
	return nil
}
