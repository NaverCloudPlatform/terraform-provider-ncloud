package classicloadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLoadBalancerCreate,
		Read:   resourceNcloudLoadBalancerRead,
		Update: resourceNcloudLoadBalancerUpdate,
		Delete: resourceNcloudLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Update: schema.DefaultTimeout(conn.DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(3, 30)),
				Description:      "Name of a load balancer to create. Default: Automatically specified by Ncloud.",
			},
			"algorithm_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RR", "LC", "SIPHS"}, false)),
				Description:      "Load balancer algorithm type code. The available algorithms are as follows: [ROUND ROBIN (RR) | LEAST_CONNECTION (LC)]. Default: ROUND ROBIN (RR)",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1000)),
				Description:      "Description of a load balancer to create",
			},
			"rule_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        loadBalancerRuleSchemaResource,
				Description: "Load balancer rules are required to create a load balancer.",
			},
			"server_instance_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of server instance numbers to be bound to the load balancer",
			},
			// Deprecated
			"internet_line_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
				Description:      "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
				Deprecated:       "This parameter is no longer used.",
			},
			"network_usage_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PBLIP", "PRVT"}, false)),
				Description:      "Network usage identification code. PBLIP(PublicIp), PRVT(PrivateIP). default : PBLIP(PublicIp)",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_operation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_http_keep_alive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"connection_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balanced_server_instance_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNcloudLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client
	config := meta.(*conn.ProviderConfig)

	if config.SupportVPC {
		return NotSupportVpc("resource `ncloud_load_balancer`")
	}

	reqParams, err := buildCreateLoadBalancerInstanceParams(d)
	if err != nil {
		return err
	}
	LogCommonRequest("CreateLoadBalancerInstance", reqParams)
	resp, err := client.Loadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		LogErrorResponse("CreateLoadBalancerInstance", err, reqParams)
		return err
	}
	LogCommonResponse("CreateLoadBalancerInstance", GetCommonResponse(resp))

	loadBalancerInstance := resp.LoadBalancerInstanceList[0]
	d.SetId(*loadBalancerInstance.LoadBalancerInstanceNo)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "USE"},
		Target:  []string{"USED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetLoadBalancerInstance(client, ncloud.StringValue(loadBalancerInstance.LoadBalancerInstanceNo))
			if err != nil {
				return 0, "", err
			}

			if ncloud.StringValue(instance.LoadBalancerInstanceOperation.Code) == "NULL" {
				return instance, ncloud.StringValue(instance.LoadBalancerInstanceStatus.Code), nil
			}

			return instance, ncloud.StringValue(instance.LoadBalancerInstanceOperation.Code), nil
		},
		Timeout:    conn.DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for LoadBalancerInstanceStatus state to be \"USED\": %s", err)
	}

	return resourceNcloudLoadBalancerRead(d, meta)
}

func resourceNcloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client

	lb, err := GetLoadBalancerInstance(client, d.Id())
	if err != nil {
		return err
	}

	if lb != nil {
		d.Set("virtual_ip", lb.VirtualIp)
		d.Set("name", lb.LoadBalancerName)
		d.Set("description", lb.LoadBalancerDescription)
		d.Set("domain_name", lb.DomainName)
		d.Set("instance_status_name", lb.LoadBalancerInstanceStatusName)
		d.Set("is_http_keep_alive", lb.IsHttpKeepAlive)
		d.Set("connection_timeout", lb.ConnectionTimeout)
		d.Set("certificate_name", lb.CertificateName)

		if algorithmType := FlattenCommonCode(lb.LoadBalancerAlgorithmType); algorithmType["code"] != nil {
			d.Set("algorithm_type", algorithmType["code"])
		}

		if instanceStatus := FlattenCommonCode(lb.LoadBalancerInstanceStatus); instanceStatus["code"] != nil {
			d.Set("instance_status", instanceStatus["code"])
		}

		if instanceOperation := FlattenCommonCode(lb.LoadBalancerInstanceOperation); instanceOperation["code"] != nil {
			d.Set("instance_operation", instanceOperation["code"])
		}

		if networkUsageType := FlattenCommonCode(lb.NetworkUsageType); networkUsageType["code"] != nil {
			d.Set("network_usage_type", networkUsageType["code"])
		}

		if len(lb.LoadBalancerRuleList) != 0 {
			if err := d.Set("rule_list", flattenLoadBalancerRuleList(lb.LoadBalancerRuleList)); err != nil {
				return err
			}
		}

		if len(lb.LoadBalancedServerInstanceList) != 0 {
			if err := d.Set("load_balanced_server_instance_list", flattenLoadBalancedServerInstanceList(lb.LoadBalancedServerInstanceList)); err != nil {
				return err
			}
		} else {
			d.Set("load_balanced_server_instance_list", nil)
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client
	if err := deleteLoadBalancerInstance(client, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client

	// Change Load Balanced Server Instances
	if d.HasChange("server_instance_no_list") {
		if err := changeLoadBalancedServerInstances(client, d); err != nil {
			return err
		}
	}

	reqParams := &loadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
		LoadBalancerInstanceNo:        ncloud.String(d.Id()),
		LoadBalancerAlgorithmTypeCode: ncloud.String(d.Get("algorithm_type").(string)),
	}

	if loadBalancerRuleParams, err := expandLoadBalancerRuleParams(d.Get("rule_list").([]interface{})); err == nil {
		reqParams.LoadBalancerRuleList = loadBalancerRuleParams
	}

	if d.HasChange("description") {
		reqParams.LoadBalancerDescription = ncloud.String(d.Get("description").(string))
	}

	if d.HasChange("algorithm_type") || d.HasChange("description") || d.HasChange("rule_list") {
		LogCommonRequest("ChangeLoadBalancerInstanceConfiguration", reqParams)
		resp, err := client.Loadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(reqParams)
		if err != nil {
			LogErrorResponse("ChangeLoadBalancerInstanceConfiguration", err, reqParams)
			return err
		}
		LogCommonResponse("ChangeLoadBalancerInstanceConfiguration", GetCommonResponse(resp))

		stateConf := &resource.StateChangeConf{
			Pending: []string{"INIT", "USE"},
			Target:  []string{"USED"},
			Refresh: func() (interface{}, string, error) {
				instance, err := GetLoadBalancerInstance(client, d.Id())
				if err != nil {
					return 0, "", err
				}

				if ncloud.StringValue(instance.LoadBalancerInstanceOperation.Code) == "NULL" {
					return instance, ncloud.StringValue(instance.LoadBalancerInstanceStatus.Code), nil
				}

				return instance, ncloud.StringValue(instance.LoadBalancerInstanceOperation.Code), nil
			},
			Timeout:    conn.DefaultUpdateTimeout,
			Delay:      2 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for LoadBalancerInstanceStatus state to be \"USED\": %s", err)
		}
	}

	return resourceNcloudLoadBalancerRead(d, meta)
}

func changeLoadBalancedServerInstances(client *conn.NcloudAPIClient, d *schema.ResourceData) error {
	reqParams := &loadbalancer.ChangeLoadBalancedServerInstancesRequest{
		LoadBalancerInstanceNo: ncloud.String(d.Id()),
		ServerInstanceNoList:   ExpandStringInterfaceList(d.Get("server_instance_no_list").([]interface{})),
	}

	LogCommonRequest("ChangeLoadBalancedServerInstances", reqParams)

	resp, err := client.Loadbalancer.V2Api.ChangeLoadBalancedServerInstances(reqParams)
	if err != nil {
		LogErrorResponse("ChangeLoadBalancedServerInstances", err, reqParams)
		return err
	}
	LogCommonResponse("ChangeLoadBalancedServerInstances", GetCommonResponse(resp))

	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "USE"},
		Target:  []string{"USED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetLoadBalancerInstance(client, d.Id())
			if err != nil {
				return 0, "", err
			}

			if ncloud.StringValue(instance.LoadBalancerInstanceOperation.Code) == "NULL" {
				return instance, ncloud.StringValue(instance.LoadBalancerInstanceStatus.Code), nil
			}

			return instance, ncloud.StringValue(instance.LoadBalancerInstanceOperation.Code), nil
		},
		Timeout:    conn.DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for LoadBalancerInstanceStatus state to be \"USED\": %s", err)
	}

	return nil
}

func buildCreateLoadBalancerInstanceParams(d *schema.ResourceData) (*loadbalancer.CreateLoadBalancerInstanceRequest, error) {
	regionNo, err := conn.ParseRegionNoParameter(d)
	if err != nil {
		return nil, err
	}

	reqParams := &loadbalancer.CreateLoadBalancerInstanceRequest{
		RegionNo: regionNo,
	}

	if loadBalancerName, ok := d.GetOk("name"); ok {
		reqParams.LoadBalancerName = ncloud.String(loadBalancerName.(string))
	}

	if loadBalancerAlgorithmTypeCode, ok := d.GetOk("algorithm_type"); ok {
		reqParams.LoadBalancerAlgorithmTypeCode = ncloud.String(loadBalancerAlgorithmTypeCode.(string))
	}

	if loadBalancerDescription, ok := d.GetOk("description"); ok {
		reqParams.LoadBalancerDescription = ncloud.String(loadBalancerDescription.(string))
	}

	if serverInstanceNoList, ok := d.GetOk("server_instance_no_list"); ok {
		reqParams.ServerInstanceNoList = ExpandStringInterfaceList(serverInstanceNoList.([]interface{}))
	}

	if networkUsageTypeCode, ok := d.GetOk("network_usage_type"); ok {
		reqParams.NetworkUsageTypeCode = ncloud.String(networkUsageTypeCode.(string))
	}

	if loadBalancerRuleParams, err := expandLoadBalancerRuleParams(d.Get("rule_list").([]interface{})); err == nil {
		reqParams.LoadBalancerRuleList = loadBalancerRuleParams
	}

	return reqParams, nil
}

func GetLoadBalancerInstance(client *conn.NcloudAPIClient, loadBalancerInstanceNo string) (*loadbalancer.LoadBalancerInstance, error) {
	reqParams := &loadbalancer.GetLoadBalancerInstanceListRequest{
		LoadBalancerInstanceNoList: []*string{ncloud.String(loadBalancerInstanceNo)},
	}
	LogCommonRequest("GetLoadBalancerInstanceList", reqParams)
	resp, err := client.Loadbalancer.V2Api.GetLoadBalancerInstanceList(reqParams)
	if err != nil {
		LogErrorResponse("GetLoadBalancerInstanceList", err, reqParams)
		return nil, err
	}
	LogCommonResponse("GetLoadBalancerInstanceList", GetCommonResponse(resp))

	for _, inst := range resp.LoadBalancerInstanceList {
		if loadBalancerInstanceNo == ncloud.StringValue(inst.LoadBalancerInstanceNo) {
			return inst, nil
		}
	}
	return nil, nil
}

func deleteLoadBalancerInstance(client *conn.NcloudAPIClient, loadBalancerInstanceNo string) error {
	reqParams := &loadbalancer.DeleteLoadBalancerInstancesRequest{
		LoadBalancerInstanceNoList: []*string{ncloud.String(loadBalancerInstanceNo)},
	}
	LogCommonRequest("DeleteLoadBalancerInstance", reqParams)
	resp, err := client.Loadbalancer.V2Api.DeleteLoadBalancerInstances(reqParams)
	if err != nil {
		LogErrorResponse("DeleteLoadBalancerInstance", err, loadBalancerInstanceNo)
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}
	LogCommonResponse("DeleteLoadBalancerInstance", commonResponse)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"", "USED"},
		Target:  []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetLoadBalancerInstance(client, loadBalancerInstanceNo)
			if err != nil {
				return 0, "", err
			}

			if instance == nil {
				return 0, "OK", err
			}

			return instance, "", nil
		},
		Timeout:    conn.DefaultUpdateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting to delete LoadBalancerInstance: %s", err)
	}

	return nil
}

var loadBalancerRuleSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"protocol_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Protocol type code of load balancer rules. The following codes are available. [HTTP | HTTPS | TCP | SSL]",
		},
		"load_balancer_port": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Load balancer port of load balancer rules",
		},
		"server_port": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Server port of load balancer rules",
		},
		"l7_health_check_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Health check path of load balancer rules. Required when the protocol_type is HTTP/HTTPS.",
		},
		"certificate_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Load balancer SSL certificate. Required when the protocol_type value is SSL/HTTPS.",
		},
		"proxy_protocol_use_yn": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "N",
			Description: "Use 'Y' if you want to check client IP addresses by enabling the proxy protocol while you select TCP or SSL.",
		},
	},
}
