package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	vpcservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
)

const (
	LoadBalancerInstanceOperationChangeCode             = "CHANG"
	LoadBalancerInstanceOperationCreateCode             = "CREAT"
	LoadBalancerInstanceOperationDisUseCode             = "DISUS"
	LoadBalancerInstanceOperationNullCode               = "NULL"
	LoadBalancerInstanceOperationPendingTerminationCode = "PTERM"
	LoadBalancerInstanceOperationTerminateCode          = "TERMT"
	LoadBalancerInstanceOperationUseCode                = "USE"
)

func ResourceNcloudLb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNcloudLbCreate,
		ReadContext:   resourceNcloudLbRead,
		UpdateContext: resourceNcloudLbUpdate,
		DeleteContext: resourceNcloudLbDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Update: schema.DefaultTimeout(conn.DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false)),
				ForceNew:         true,
			},
			"idle_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 3600)),
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"APPLICATION", "NETWORK", "NETWORK_PROXY"}, false)),
				ForceNew:         true,
			},
			"throughput_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "DYNAMIC"}, false)),
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
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb`"))
	}

	throughput_type := StringPtrOrNil(d.GetOk("throughput_type"))
	if (d.Get("type").(string) == "NETWORK") && (throughput_type != nil && *throughput_type != "DYNAMIC") {
		return diag.FromErr(fmt.Errorf("Network Loadbalancer throughput_type can only be set to empty or DYNAMIC"))
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

	// The subnet must be the same VPC, so the size of vpcNoMap must be 1.
	vpcNoMap := make(map[string]int)
	subnetList := make([]*vpc.Subnet, 0)
	for _, subnetNo := range reqParams.SubnetNoList {
		subnet, err := vpcservice.GetSubnetInstance(config, *subnetNo)
		if err != nil {
			return diag.FromErr(err)
		}
		if subnet == nil {
			return diag.FromErr(fmt.Errorf("not found subnet(%s)", *subnetNo))
		}
		subnetList = append(subnetList, subnet)
		vpcNoMap[*subnet.VpcNo]++
	}

	if len(vpcNoMap) > 1 {
		return diag.FromErr(fmt.Errorf("subnet must be set to the subnet of the same vpc"))
	}

	reqParams.VpcNo = subnetList[0].VpcNo

	LogCommonRequest("createLoadBalancerInstance", reqParams)
	resp, err := config.Client.Vloadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		LogErrorResponse("createLoadBalancerInstance", err, reqParams)
		return diag.FromErr(err)
	}
	LogResponse("createLoadBalancerInstance", resp)

	if err := waitForLoadBalancerActive(ctx, d, config, ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
	return resourceNcloudLbRead(ctx, d, meta)
}

func resourceNcloudLbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb`"))
	}

	lb, err := GetVpcLoadBalancer(config, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if lb == nil {
		d.SetId("")
		return nil
	}

	lbMap := ConvertToMap(lb)
	SetSingularResourceDataFromMapSchema(ResourceNcloudLb(), d, lbMap)
	return nil
}

func resourceNcloudLbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*conn.ProviderConfig)
	if !config.SupportVPC {
		return diag.FromErr(NotSupportClassic("resource `ncloud_lb`"))
	}
	if d.HasChanges("idle_timeout", "throughput_type") {
		if err := waitForLoadBalancerActive(ctx, d, config, d.Id()); err != nil {
			return diag.FromErr(err)
		}
		_, err := config.Client.Vloadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(&vloadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
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
		_, err := config.Client.Vloadbalancer.V2Api.SetLoadBalancerDescription(&vloadbalancer.SetLoadBalancerDescriptionRequest{
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
	config := meta.(*conn.ProviderConfig)
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

	LogCommonRequest("resourceNcloudLbDelete", deleteInstanceReqParams)
	if _, err := config.Client.Vloadbalancer.V2Api.DeleteLoadBalancerInstances(deleteInstanceReqParams); err != nil {
		LogErrorResponse("resourceNcloudLbDelete", err, deleteInstanceReqParams)
		return diag.FromErr(err)
	}

	if err := waitForLoadBalancerDeletion(ctx, d, config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForLoadBalancerDeletion(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{LoadBalancerInstanceOperationTerminateCode},
		Target:  []string{LoadBalancerInstanceOperationNullCode},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(d.Id()),
			}
			resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
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

func waitForLoadBalancerActive(ctx context.Context, d *schema.ResourceData, config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{LoadBalancerInstanceOperationCreateCode, LoadBalancerInstanceOperationChangeCode},
		Target:  []string{LoadBalancerInstanceOperationNullCode},
		Refresh: func() (result interface{}, state string, err error) {
			reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
				RegionCode:             &config.RegionCode,
				LoadBalancerInstanceNo: ncloud.String(id),
			}
			resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
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

func GetVpcLoadBalancer(config *conn.ProviderConfig, id string) (*LoadBalancerInstance, error) {
	reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(id),
	}
	LogCommonRequest("getLoadBalancerInstanceDetail", reqParams)

	resp, err := config.Client.Vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
	if err != nil {
		LogErrorResponse("getLoadBalancerInstanceDetail", err, reqParams)
		return nil, err
	}
	LogResponse("getLoadBalancerInstanceDetail", resp)

	if len(resp.LoadBalancerInstanceList) < 1 {
		return nil, nil
	}

	return convertVpcLoadBalancer(resp.LoadBalancerInstanceList[0]), nil
}

func convertVpcLoadBalancer(instance *vloadbalancer.LoadBalancerInstance) *LoadBalancerInstance {
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
