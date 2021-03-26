package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func init() {
	RegisterResource("ncloud_lb", resourceNcloudLb())
}

const (
	LoadBalancerInstanceOperationChangeCode             = "CHANG"
	LoadBalancerInstanceOperationCreateCode             = "CREAT"
	LoadBalancerInstanceOperationDisUseCode             = "DISUS"
	LoadBalancerInstanceOperationNullCode               = "NULL"
	LoadBalancerInstanceOperationPendingTerminationCode = "PTERM"
	LoadBalancerInstanceOperationTerminateCode          = "TERMT"
	LoadBalancerInstanceOperationUseCode                = "USE"
)

func resourceNcloudLb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudLbCreate,
		ReadContext:   resourceNcloudLbRead,
		UpdateContext: resourceNcloudLbUpdate,
		DeleteContext: resourceNcloudLbDelete,
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
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ToDiagFunc(validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false)),
				ForceNew:         true,
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
				ForceNew:         true,
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
				ForceNew: true,
			},
			"ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener_no_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNcloudLbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb`"))
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
		LoadBalancerTypeCode: ncloud.String(d.Get("type").(string)),
		SubnetNoList:         ncloud.StringInterfaceList(d.Get("subnet_no_list").([]interface{})),
	}
	subnet, err := getSubnetInstance(config, *reqParams.SubnetNoList[0])
	if err != nil {
		return diag.FromErr(err)
	}

	reqParams.VpcNo = subnet.VpcNo
	logCommonRequest("resourceNcloudLbCreate", reqParams)
	resp, err := config.Client.vloadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudLbCreate", err, reqParams)
		return diag.FromErr(err)
	}
	logResponse("resourceNcloudLbCreate", resp)
	if err := waitForLoadBalancerActive(ctx, d, config, ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
	return resourceNcloudLbRead(ctx, d, meta)
}

func resourceNcloudLbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb`"))
	}

	reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(d.Id()),
	}
	resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resp.LoadBalancerInstanceList) < 1 {
		d.SetId("")
		return nil
	}
	lb := convertLbInstance(resp.LoadBalancerInstanceList[0])
	lbMap := ConvertToMap(lb)
	SetSingularResourceDataFromMapSchema(resourceNcloudLb(), d, lbMap)
	return nil
}

func resourceNcloudLbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb`"))
	}
	if d.HasChanges("idle_timeout", "throughput_type") {
		if err := waitForLoadBalancerActive(ctx, d, config, d.Id()); err != nil {
			return diag.FromErr(err)
		}
		_, err := config.Client.vloadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(&vloadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
			RegionCode:             &config.RegionCode,
			LoadBalancerInstanceNo: ncloud.String(d.Id()),
			IdleTimeout:            Int32PtrOrNil(d.GetOk("idle_timeout")),
			ThroughputTypeCode:     StringPtrOrNil(d.GetOk("throughput_type")),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("description") {
		if err := waitForLoadBalancerActive(ctx, d, config, d.Id()); err != nil {
			return diag.FromErr(err)
		}
		_, err := config.Client.vloadbalancer.V2Api.SetLoadBalancerDescription(&vloadbalancer.SetLoadBalancerDescriptionRequest{
			RegionCode:              &config.RegionCode,
			LoadBalancerInstanceNo:  ncloud.String(d.Id()),
			LoadBalancerDescription: StringPtrOrNil(d.GetOk("description")),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceNcloudLbRead(ctx, d, config)
}

func resourceNcloudLbDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb`"))
	}
	deleteInstanceReqParams := &vloadbalancer.DeleteLoadBalancerInstancesRequest{
		RegionCode:                 &config.RegionCode,
		LoadBalancerInstanceNoList: ncloud.StringList([]string{d.Id()}),
	}

	if err := waitForLoadBalancerActive(ctx, d, config, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	logCommonRequest("resourceNcloudLbDelete", deleteInstanceReqParams)
	if _, err := config.Client.vloadbalancer.V2Api.DeleteLoadBalancerInstances(deleteInstanceReqParams); err != nil {
		logErrorResponse("resourceNcloudLbDelete", err, deleteInstanceReqParams)
		return diag.FromErr(err)
	}

	if err := waitForLoadBalancerDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForLoadBalancerDeletion(ctx context.Context, d *schema.ResourceData, config *ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{LoadBalancerInstanceOperationTerminateCode},
		Target:  []string{LoadBalancerInstanceOperationNullCode},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(d.Id()),
			}
			resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
			if err != nil {
				return nil, "", err
			}

			if len(resp.LoadBalancerInstanceList) < 1 {
				return resp, LoadBalancerInstanceOperationNullCode, nil
			}

			lb := resp.LoadBalancerInstanceList[0]
			return resp, ncloud.StringValue(lb.LoadBalancerInstanceOperation.Code), nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for Load Balancer instance (%s) to become terminating: %s", d.Id(), err)
	}
	return nil
}

func waitForLoadBalancerActive(ctx context.Context, d *schema.ResourceData, config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{LoadBalancerInstanceOperationCreateCode, LoadBalancerInstanceOperationChangeCode},
		Target:  []string{LoadBalancerInstanceOperationNullCode},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(id),
			}
			resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
			if err != nil {
				return nil, "", err
			}

			if len(resp.LoadBalancerInstanceList) < 1 {
				return nil, "", fmt.Errorf("not found load balancer instance(%s)", id)
			}

			lb := resp.LoadBalancerInstanceList[0]
			return resp, ncloud.StringValue(lb.LoadBalancerInstanceOperation.Code), nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 3 * time.Second,
		Delay:      2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for Load Balancer instance (%s) to become activating: %s", id, err)
	}
	return nil
}

func convertLbInstance(instance *vloadbalancer.LoadBalancerInstance) *LoadBalancerInstance {
	return &LoadBalancerInstance{
		LoadBalancerInstanceNo:   instance.LoadBalancerInstanceNo,
		LoadBalancerDescription:  instance.LoadBalancerDescription,
		LoadBalancerName:         instance.LoadBalancerName,
		LoadBalancerDomain:       instance.LoadBalancerDomain,
		LoadBalancerIpList:       instance.LoadBalancerIpList,
		LoadBalancerType:         instance.LoadBalancerType.Code,
		LoadBalancerNetworkType:  instance.LoadBalancerNetworkType.Code,
		ThroughputType:           instance.ThroughputType.Code,
		IdleTimeout:              instance.IdleTimeout,
		VpcNo:                    instance.VpcNo,
		SubnetNoList:             instance.SubnetNoList,
		LoadBalancerListenerList: instance.LoadBalancerListenerNoList,
	}
}
