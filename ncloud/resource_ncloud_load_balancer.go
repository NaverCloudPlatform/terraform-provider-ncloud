package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
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
				Description:  "Network usage identification code. PBLIP(PublicIP), PRVT(PrivateIP). default : PBLIP(PublicIP)",
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
	conn := meta.(*NcloudSdk).conn

	reqParams := buildCreateLoadBalancerInstanceParams(conn, d)
	resp, err := conn.CreateLoadBalancerInstance(reqParams)
	if err != nil {
		logErrorResponse("CreateLoadBalancerInstance", err, reqParams)
		return err
	}
	logCommonResponse("CreateLoadBalancerInstance", reqParams, resp.CommonResponse)

	LoadBalancerInstance := &resp.LoadBalancerInstanceList[0]
	d.SetId(LoadBalancerInstance.LoadBalancerInstanceNo)

	if err := waitForLoadBalancerInstance(conn, LoadBalancerInstance.LoadBalancerInstanceNo, "USED", DefaultCreateTimeout); err != nil {
		return err
	}
	return resourceNcloudLoadBalancerRead(d, meta)
}

func resourceNcloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	lb, err := getLoadBalancerInstance(conn, d.Id())
	if err != nil {
		return err
	}
	if lb != nil {
		d.Set("virtual_ip", lb.VirtualIP)
		d.Set("load_balancer_name", lb.LoadBalancerName)
		d.Set("load_balancer_algorithm_type", setCommonCode(lb.LoadBalancerAlgorithmType))
		d.Set("load_balancer_description", lb.LoadBalancerDescription)
		d.Set("create_date", lb.CreateDate)
		d.Set("domain_name", lb.DomainName)
		d.Set("internet_line_type", setCommonCode(lb.InternetLineType))
		d.Set("load_balancer_instance_status_name", lb.LoadBalancerInstanceStatusName)
		d.Set("load_balancer_instance_status", setCommonCode(lb.LoadBalancerInstanceStatus))
		d.Set("load_balancer_instance_operation", setCommonCode(lb.LoadBalancerInstanceOperation))
		d.Set("network_usage_type", setCommonCode(lb.NetworkUsageType))
		d.Set("is_http_keep_alive", lb.IsHTTPKeepAlive)
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

func getLoadBalancerRuleList(lbRuleList []sdk.LoadBalancerRule) []interface{} {
	list := make([]interface{}, 0, len(lbRuleList))

	for _, r := range lbRuleList {
		rule := map[string]interface{}{
			"protocol_type_code":    setCommonCode(r.ProtocolType),
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

func getLoadBalancedServerInstanceList(loadBalancedServerInstanceList []sdk.LoadBalancedServerInstance) []string {
	list := make([]string, 0, len(loadBalancedServerInstanceList))

	for _, instance := range loadBalancedServerInstanceList {
		list = append(list, instance.ServerInstanceList[0].ServerInstanceNo)
	}

	return list
}

func resourceNcloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn
	return deleteLoadBalancerInstance(conn, d.Id())
}

func resourceNcloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	// Change Load Balanced Server Instances
	if d.HasChange("server_instance_no_list") {
		if err := changeLoadBalancedServerInstances(conn, d); err != nil {
			return err
		}
	}

	reqParams := &sdk.RequestChangeLoadBalancerInstanceConfiguration{
		LoadBalancerInstanceNo:        d.Id(),
		LoadBalancerAlgorithmTypeCode: d.Get("load_balancer_algorithm_type_code").(string),
		LoadBalancerRuleList:          buildLoadBalancerRuleParams(d),
	}

	if d.HasChange("load_balancer_description") {
		reqParams.LoadBalancerDescription = d.Get("load_balancer_description").(string)
	}

	if d.HasChange("load_balancer_algorithm_type_code") || d.HasChange("load_balancer_description") || d.HasChange("load_balancer_rule_list") {
		resp, err := conn.ChangeLoadBalancerInstanceConfiguration(reqParams)
		if err != nil {
			logErrorResponse("ChangeLoadBalancerInstanceConfiguration", err, reqParams)
			return err
		}
		logCommonResponse("ChangeLoadBalancerInstanceConfiguration", reqParams, resp.CommonResponse)

		if err := waitForLoadBalancerInstance(conn, d.Id(), "USED", DefaultUpdateTimeout); err != nil {
			return err
		}
	}

	return resourceNcloudLoadBalancerRead(d, meta)
}

func changeLoadBalancedServerInstances(conn *sdk.Conn, d *schema.ResourceData) error {
	reqParams := &sdk.RequestChangeLoadBalancedServerInstances{
		LoadBalancerInstanceNo: d.Id(),
		ServerInstanceNoList:   StringList(d.Get("server_instance_no_list").([]interface{})),
	}

	resp, err := conn.ChangeLoadBalancedServerInstances(reqParams)
	if err != nil {
		logErrorResponse("ChangeLoadBalancedServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("ChangeLoadBalancedServerInstances", reqParams, resp.CommonResponse)

	if err := waitForLoadBalancerInstance(conn, d.Id(), "USED", DefaultUpdateTimeout); err != nil {
		return err
	}

	return nil
}

func buildLoadBalancerRuleParams(d *schema.ResourceData) []sdk.RequestLoadBalancerRule {
	lbRuleList := make([]sdk.RequestLoadBalancerRule, 0, len(d.Get("load_balancer_rule_list").([]interface{})))

	for _, v := range d.Get("load_balancer_rule_list").([]interface{}) {
		lbRule := new(sdk.RequestLoadBalancerRule)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "protocol_type_code":
				lbRule.ProtocolTypeCode = value.(string)
			case "load_balancer_port":
				lbRule.LoadBalancerPort = value.(int)
			case "server_port":
				lbRule.ServerPort = value.(int)
			case "l7_health_check_path":
				lbRule.L7HealthCheckPath = value.(string)
			case "certificate_name":
				lbRule.CertificateName = value.(string)
			case "proxy_protocol_use_yn":
				lbRule.ProxyProtocolUseYn = value.(string)
			}
		}
		lbRuleList = append(lbRuleList, *lbRule)
	}

	return lbRuleList
}

func buildCreateLoadBalancerInstanceParams(conn *sdk.Conn, d *schema.ResourceData) *sdk.RequestCreateLoadBalancerInstance {
	reqParams := &sdk.RequestCreateLoadBalancerInstance{
		LoadBalancerName:              d.Get("load_balancer_name").(string),
		LoadBalancerAlgorithmTypeCode: d.Get("load_balancer_algorithm_type_code").(string),
		LoadBalancerDescription:       d.Get("load_balancer_description").(string),
		LoadBalancerRuleList:          buildLoadBalancerRuleParams(d),
		ServerInstanceNoList:          StringList(d.Get("server_instance_no_list").([]interface{})),
		InternetLineTypeCode:          d.Get("internet_line_type_code").(string),
		NetworkUsageTypeCode:          d.Get("network_usage_type_code").(string),
		RegionNo:                      parseRegionNoParameter(conn, d),
	}
	return reqParams
}

func getLoadBalancerInstance(conn *sdk.Conn, LoadBalancerInstanceNo string) (*sdk.LoadBalancerInstance, error) {
	reqParams := &sdk.RequestLoadBalancerInstanceList{
		LoadBalancerInstanceNoList: []string{LoadBalancerInstanceNo},
	}
	resp, err := conn.GetLoadBalancerInstanceList(reqParams)
	if err != nil {
		logErrorResponse("GetLoadBalancerInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetLoadBalancerInstanceList", reqParams, resp.CommonResponse)

	for _, inst := range resp.LoadBalancerInstanceList {
		if LoadBalancerInstanceNo == inst.LoadBalancerInstanceNo {
			return &inst, nil
		}
	}
	return nil, nil
}

func deleteLoadBalancerInstance(conn *sdk.Conn, LoadBalancerInstanceNo string) error {
	reqParams := &sdk.RequestDeleteLoadBalancerInstances{
		LoadBalancerInstanceNoList: []string{LoadBalancerInstanceNo},
	}
	resp, err := conn.DeleteLoadBalancerInstances(reqParams)
	if err != nil {
		logErrorResponse("DeleteLoadBalancerInstance", err, LoadBalancerInstanceNo)
		return err
	}
	var commonResponse = common.CommonResponse{}
	if resp != nil {
		commonResponse = resp.CommonResponse
	}
	logCommonResponse("DeleteLoadBalancerInstance", LoadBalancerInstanceNo, commonResponse)

	return waitForDeleteLoadBalancerInstance(conn, LoadBalancerInstanceNo)
}

func waitForLoadBalancerInstance(conn *sdk.Conn, id string, status string, timeout time.Duration) error {
	c1 := make(chan error, 1)

	go func() {
		for {
			instance, err := getLoadBalancerInstance(conn, id)

			if err != nil {
				c1 <- err
				return
			}

			if instance == nil || (instance.LoadBalancerInstanceStatus.Code == status && instance.LoadBalancerInstanceOperation.Code == "NULL") {
				c1 <- nil
				return
			}

			log.Printf("[DEBUG] Wait get load balancer instance [%s] status [%s] to be [%s]", id, instance.LoadBalancerInstanceStatus.Code, status)
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

func waitForDeleteLoadBalancerInstance(conn *sdk.Conn, id string) error {
	c1 := make(chan error, 1)

	go func() {
		for {
			instance, err := getLoadBalancerInstance(conn, id)

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
