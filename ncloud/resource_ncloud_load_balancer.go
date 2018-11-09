package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/loadbalancer"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudLoadBalancerCreate,
		Read:   resourceNcloudLoadBalancerRead,
		Delete: resourceNcloudLoadBalancerDelete,
		Update: resourceNcloudLoadBalancerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Update: schema.DefaultTimeout(DefaultUpdateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(3, 30),
				Description:  "Name of a load balancer to create. Default: Automatically specified by Ncloud.",
			},
			"load_balancer_algorithm_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIncludeValues([]string{"RR", "LC", "SIPHS"}),
				Description:  "Load balancer algorithm type code. The available algorithms are as follows: [ROUND ROBIN (RR) | LEAST_CONNECTION (LC)]. Default: ROUND ROBIN (RR)",
			},
			"load_balancer_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1000),
				Description:  "Description of a load balancer to create",
			},
			"load_balancer_rule_list": {
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
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIncludeValues([]string{"PUBLC", "GLBL"}),
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"network_usage_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIncludeValues([]string{"PBLIP", "PRVT"}),
				Description:  "Network usage identification code. PBLIP(PublicIp), PRVT(PrivateIP). default : PBLIP(PublicIp)",
			},
			"region_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region code. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_no"},
			},
			"region_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region number. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_code"},
			},
			"load_balancer_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_algorithm_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_line_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"load_balancer_instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"load_balancer_instance_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"network_usage_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
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
	client := meta.(*NcloudAPIClient)

	reqParams, err := buildCreateLoadBalancerInstanceParams(client, d)
	if err != nil {
		return err
	}
	logCommonRequest("CreateLoadBalancerInstance", reqParams)
	resp, err := client.loadbalancer.V2Api.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateLoadBalancerInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateLoadBalancerInstance", GetCommonResponse(resp))

	loadBalancerInstance := resp.LoadBalancerInstanceList[0]
	d.SetId(*loadBalancerInstance.LoadBalancerInstanceNo)

	if err := waitForLoadBalancerInstance(client, ncloud.StringValue(loadBalancerInstance.LoadBalancerInstanceNo), "USED", DefaultCreateTimeout); err != nil {
		return err
	}
	return resourceNcloudLoadBalancerRead(d, meta)
}

func resourceNcloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	lb, err := getLoadBalancerInstance(client, d.Id())
	if err != nil {
		return err
	}
	if lb != nil {
		d.Set("virtual_ip", lb.VirtualIp)
		d.Set("load_balancer_name", lb.LoadBalancerName)
		d.Set("load_balancer_algorithm_type", flattenCommonCode(lb.LoadBalancerAlgorithmType))
		d.Set("load_balancer_description", lb.LoadBalancerDescription)
		d.Set("create_date", lb.CreateDate)
		d.Set("domain_name", lb.DomainName)
		d.Set("internet_line_type", flattenCommonCode(lb.InternetLineType))
		d.Set("load_balancer_instance_status_name", lb.LoadBalancerInstanceStatusName)
		d.Set("load_balancer_instance_status", flattenCommonCode(lb.LoadBalancerInstanceStatus))
		d.Set("load_balancer_instance_operation", flattenCommonCode(lb.LoadBalancerInstanceOperation))
		d.Set("network_usage_type", flattenCommonCode(lb.NetworkUsageType))
		d.Set("is_http_keep_alive", lb.IsHttpKeepAlive)
		d.Set("connection_timeout", lb.ConnectionTimeout)
		d.Set("certificate_name", lb.CertificateName)

		if len(lb.LoadBalancerRuleList) != 0 {
			d.Set("load_balancer_rule_list", getLoadBalancerRuleList(lb.LoadBalancerRuleList))
		}
		if len(lb.LoadBalancedServerInstanceList) != 0 {
			d.Set("load_balanced_server_instance_list", getLoadBalancedServerInstanceList(lb.LoadBalancedServerInstanceList))
		} else {
			d.Set("load_balanced_server_instance_list", nil)
		}
	}

	return nil
}

func getLoadBalancerRuleList(lbRuleList []*loadbalancer.LoadBalancerRule) []interface{} {
	list := make([]interface{}, 0, len(lbRuleList))

	for _, r := range lbRuleList {
		rule := map[string]interface{}{
			"protocol_type_code":    flattenCommonCode(r.ProtocolType),
			"load_balancer_port":    r.LoadBalancerPort,
			"server_port":           r.ServerPort,
			"l7_health_check_path":  r.L7HealthCheckPath,
			"certificate_name":      r.CertificateName,
			"proxy_protocol_use_yn": r.ProxyProtocolUseYn,
		}
		list = append(list, rule)
	}

	return list
}

func getLoadBalancedServerInstanceList(loadBalancedServerInstanceList []*loadbalancer.LoadBalancedServerInstance) []string {
	list := make([]string, 0, len(loadBalancedServerInstanceList))

	for _, instance := range loadBalancedServerInstanceList {
		list = append(list, ncloud.StringValue(instance.ServerInstance.ServerInstanceNo))
	}

	return list
}

func resourceNcloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	return deleteLoadBalancerInstance(client, d.Id())
}

func resourceNcloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	// Change Load Balanced Server Instances
	if d.HasChange("server_instance_no_list") {
		if err := changeLoadBalancedServerInstances(client, d); err != nil {
			return err
		}
	}

	reqParams := &loadbalancer.ChangeLoadBalancerInstanceConfigurationRequest{
		LoadBalancerInstanceNo:        ncloud.String(d.Id()),
		LoadBalancerAlgorithmTypeCode: ncloud.String(d.Get("load_balancer_algorithm_type_code").(string)),
		LoadBalancerRuleList:          buildLoadBalancerRuleParams(d),
	}

	if d.HasChange("load_balancer_description") {
		reqParams.LoadBalancerDescription = ncloud.String(d.Get("load_balancer_description").(string))
	}

	if d.HasChange("load_balancer_algorithm_type_code") || d.HasChange("load_balancer_description") || d.HasChange("load_balancer_rule_list") {
		logCommonRequest("ChangeLoadBalancerInstanceConfiguration", reqParams)
		resp, err := client.loadbalancer.V2Api.ChangeLoadBalancerInstanceConfiguration(reqParams)
		if err != nil {
			logErrorResponse("ChangeLoadBalancerInstanceConfiguration", err, reqParams)
			return err
		}
		logCommonResponse("ChangeLoadBalancerInstanceConfiguration", GetCommonResponse(resp))

		if err := waitForLoadBalancerInstance(client, d.Id(), "USED", DefaultUpdateTimeout); err != nil {
			return err
		}
	}

	return resourceNcloudLoadBalancerRead(d, meta)
}

func changeLoadBalancedServerInstances(client *NcloudAPIClient, d *schema.ResourceData) error {
	reqParams := &loadbalancer.ChangeLoadBalancedServerInstancesRequest{
		LoadBalancerInstanceNo: ncloud.String(d.Id()),
		ServerInstanceNoList:   expandStringInterfaceList(d.Get("server_instance_no_list").([]interface{})),
	}

	logCommonRequest("ChangeLoadBalancedServerInstances", reqParams)

	resp, err := client.loadbalancer.V2Api.ChangeLoadBalancedServerInstances(reqParams)
	if err != nil {
		logErrorResponse("ChangeLoadBalancedServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("ChangeLoadBalancedServerInstances", GetCommonResponse(resp))

	if err := waitForLoadBalancerInstance(client, d.Id(), "USED", DefaultUpdateTimeout); err != nil {
		return err
	}

	return nil
}

func buildLoadBalancerRuleParams(d *schema.ResourceData) []*loadbalancer.LoadBalancerRuleParameter {
	lbRuleList := make([]*loadbalancer.LoadBalancerRuleParameter, 0, len(d.Get("load_balancer_rule_list").([]interface{})))

	for _, v := range d.Get("load_balancer_rule_list").([]interface{}) {
		lbRule := new(loadbalancer.LoadBalancerRuleParameter)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "protocol_type_code":
				lbRule.ProtocolTypeCode = ncloud.String(value.(string))
			case "load_balancer_port":
				lbRule.LoadBalancerPort = ncloud.Int32(int32(value.(int)))
			case "server_port":
				lbRule.ServerPort = ncloud.Int32(int32(value.(int)))
			case "l7_health_check_path":
				lbRule.L7HealthCheckPath = ncloud.String(value.(string))
			case "certificate_name":
				lbRule.CertificateName = ncloud.String(value.(string))
			case "proxy_protocol_use_yn":
				lbRule.ProxyProtocolUseYn = ncloud.String(value.(string))
			}
		}
		lbRuleList = append(lbRuleList, lbRule)
	}

	return lbRuleList
}

func buildCreateLoadBalancerInstanceParams(client *NcloudAPIClient, d *schema.ResourceData) (*loadbalancer.CreateLoadBalancerInstanceRequest, error) {
	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return nil, err
	}
	reqParams := &loadbalancer.CreateLoadBalancerInstanceRequest{
		LoadBalancerName:              ncloud.String(d.Get("load_balancer_name").(string)),
		LoadBalancerAlgorithmTypeCode: ncloud.String(d.Get("load_balancer_algorithm_type_code").(string)),
		LoadBalancerDescription:       ncloud.String(d.Get("load_balancer_description").(string)),
		LoadBalancerRuleList:          buildLoadBalancerRuleParams(d),
		ServerInstanceNoList:          expandStringInterfaceList(d.Get("server_instance_no_list").([]interface{})),
		InternetLineTypeCode:          StringPtrOrNil(d.GetOk("internet_line_type_code")),
		NetworkUsageTypeCode:          ncloud.String(d.Get("network_usage_type_code").(string)),
		RegionNo:                      regionNo,
	}
	return reqParams, nil
}

func getLoadBalancerInstance(client *NcloudAPIClient, loadBalancerInstanceNo string) (*loadbalancer.LoadBalancerInstance, error) {
	reqParams := &loadbalancer.GetLoadBalancerInstanceListRequest{
		LoadBalancerInstanceNoList: []*string{ncloud.String(loadBalancerInstanceNo)},
	}
	logCommonRequest("GetLoadBalancerInstanceList", reqParams)
	resp, err := client.loadbalancer.V2Api.GetLoadBalancerInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetLoadBalancerInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetLoadBalancerInstanceList", GetCommonResponse(resp))

	for _, inst := range resp.LoadBalancerInstanceList {
		if loadBalancerInstanceNo == ncloud.StringValue(inst.LoadBalancerInstanceNo) {
			return inst, nil
		}
	}
	return nil, nil
}

func deleteLoadBalancerInstance(client *NcloudAPIClient, loadBalancerInstanceNo string) error {
	reqParams := &loadbalancer.DeleteLoadBalancerInstancesRequest{
		LoadBalancerInstanceNoList: []*string{ncloud.String(loadBalancerInstanceNo)},
	}
	logCommonRequest("DeleteLoadBalancerInstance", reqParams)
	resp, err := client.loadbalancer.V2Api.DeleteLoadBalancerInstances(reqParams)
	if err != nil {
		logErrorResponse("DeleteLoadBalancerInstance", err, loadBalancerInstanceNo)
		return err
	}
	var commonResponse = &CommonResponse{}
	if resp != nil {
		commonResponse = GetCommonResponse(resp)
	}
	logCommonResponse("DeleteLoadBalancerInstance", commonResponse)

	return waitForDeleteLoadBalancerInstance(client, loadBalancerInstanceNo)
}

func waitForLoadBalancerInstance(client *NcloudAPIClient, id string, status string, timeout time.Duration) error {
	c1 := make(chan error, 1)

	go func() {
		for {
			instance, err := getLoadBalancerInstance(client, id)

			if err != nil {
				c1 <- err
				return
			}

			if instance == nil || (ncloud.StringValue(instance.LoadBalancerInstanceStatus.Code) == status && ncloud.StringValue(instance.LoadBalancerInstanceOperation.Code) == "NULL") {
				c1 <- nil
				return
			}

			log.Printf("[DEBUG] Wait get load balancer instance [%s] status [%s] to be [%s]", id, *instance.LoadBalancerInstanceStatus.Code, status)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(timeout):
		return fmt.Errorf("TIMEOUT : delete load balancer instance [%s] ", id)
	}
}

func waitForDeleteLoadBalancerInstance(client *NcloudAPIClient, id string) error {
	c1 := make(chan error, 1)

	go func() {
		for {
			instance, err := getLoadBalancerInstance(client, id)

			if err != nil {
				c1 <- err
				return
			}

			if instance == nil {
				c1 <- nil
				return
			}

			log.Printf("[DEBUG] Wait delete load balancer instance [%s] ", id)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultTimeout):
		return fmt.Errorf("TIMEOUT : delete load balancer instance [%s] ", id)
	}
}

var loadBalancerRuleSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"protocol_type_code": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Protocol type code of load balancer rules. The following codes are available. [HTTP | HTTPS | TCP | SSL]",
		},
		"protocol_type": {
			Type:        schema.TypeMap,
			Computed:    true,
			Elem:        commonCodeSchemaResource,
			Description: "",
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
			Description: "Health check path of load balancer rules. Required when the protocol_type_code is HTTP/HTTPS.",
		},
		"certificate_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Load balancer SSL certificate. Required when the protocol_type_code value is SSL/HTTPS.",
		},
		"proxy_protocol_use_yn": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Use 'Y' if you want to check client IP addresses by enabling the proxy protocol while you select TCP or SSL.",
		},
	},
}
