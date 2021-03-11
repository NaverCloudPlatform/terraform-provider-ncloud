package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		},
	}
}

func resourceNcloudLbCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_lb`")
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
		return err
	}

	reqParams.VpcNo = subnet.VpcNo
	logCommonRequest("resourceNcloudLbCreate", reqParams)
	resp, err := config.Client.vloadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		logErrorResponse("resourceNcloudLbCreate", err, reqParams)
		return err
	}
	logResponse("resourceNcloudLbCreate", resp)
	if err := waitForLoadBalancerActive(d, config, resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo); err != nil {
		return err
	}
	d.SetId(ncloud.StringValue(resp.LoadBalancerInstanceList[0].LoadBalancerInstanceNo))
	return resourceNcloudLbRead(d, meta)
}

func resourceNcloudLbRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_lb`")
	}

	reqParams := &vloadbalancer.GetLoadBalancerInstanceDetailRequest{
		RegionCode:             &config.RegionCode,
		LoadBalancerInstanceNo: ncloud.String(d.Id()),
	}
	resp, err := config.Client.vloadbalancer.V2Api.GetLoadBalancerInstanceDetail(reqParams)
	if err != nil {
		return err
	}
	lb := convertLbInstance(resp.LoadBalancerInstanceList[0])
	lbMap := ConvertToMap(lb)
	SetSingularResourceDataFromMapSchema(resourceNcloudLb(), d, lbMap)
	return nil
}

func resourceNcloudLbUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_lb`")
	}
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

	if d.HasChanges("description") {
		_, err := config.Client.vloadbalancer.V2Api.SetLoadBalancerDescription(&vloadbalancer.SetLoadBalancerDescriptionRequest{
			RegionCode:              &config.RegionCode,
			LoadBalancerInstanceNo:  ncloud.String(d.Id()),
			LoadBalancerDescription: StringPtrOrNil(d.GetOk("description")),
		})
		if err != nil {
			return err
		}
	}
	return resourceNcloudLbRead(d, config)
}

func resourceNcloudLbDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_lb`")
	}
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
