package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceNcloudPortForwadingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudPortForwardingRuleCreate,
		Read:   resourceNcloudPortForwardingRuleRead,
		Update: resourceNcloudPortForwardingRuleUpdate,
		Delete: resourceNcloudPortForwardingRuleDelete,
		Exists: resourceNcloudPortForwardingRuleExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"port_forwarding_configuration_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port forwarding configuration number.",
			},
			"server_instance_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_forwarding_external_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1024, 65534),
			},
			"port_forwarding_internal_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIncludeValues([]string{"22", "3389"}), // [Linux : 22 |Windows : 3389]
			},
			"port_forwarding_public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudPortForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudPortForwardingRuleCreate")
	conn := meta.(*NcloudSdk).conn

	portForwardingConfigurationNo, err := getPortForwardingConfigurationNo(d, meta)
	if err != nil {
		return err
	}

	portForwardingExternalPort := d.Get("port_forwarding_external_port").(string)

	reqParams := &sdk.RequestAddPortForwardingRules{
		PortForwardingConfigurationNo: portForwardingConfigurationNo,
		PortForwardingRuleList: []sdk.PortForwardingRule{
			{
				ServerInstanceNo:           d.Get("server_instance_no").(string),
				PortForwardingExternalPort: portForwardingExternalPort,
				PortForwardingInternalPort: d.Get("port_forwarding_internal_port").(string),
			},
		},
	}
	resp, err := conn.AddPortForwardingRules(reqParams)
	if err != nil {
		logErrorResponse("AddPortForwardingRules", err, reqParams)
		return err
	}
	logCommonResponse("AddPortForwardingRules", reqParams, resp.CommonResponse)

	d.SetId(portForwardingExternalPort)
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func getPortForwardingConfigurationNo(d *schema.ResourceData, meta interface{}) (string, error) {
	conn := meta.(*NcloudSdk).conn
	portForwardingConfigurationNo, ok := d.GetOk("port_forwarding_configuration_no")
	if !ok {
		getRuleRes, err := getPortForwardingRuleList(conn)
		if err != nil {
			return "", err
		}
		portForwardingConfigurationNo = getRuleRes.PortForwardingConfigurationNo
	}
	return portForwardingConfigurationNo.(string), nil
}

func resourceNcloudPortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudPortForwardingRuleRead")
	conn := meta.(*NcloudSdk).conn
	portForwardingExternalPort := d.Get("port_forwarding_external_port").(string)

	rule, err := getPortForwardingRule(conn, portForwardingExternalPort)
	if err != nil {
		return err
	}

	if rule != nil {
		d.Set("server_instance_no", rule.ServerInstanceNo)
		d.Set("port_forwarding_external_port", rule.PortForwardingExternalPort)
		d.Set("port_forwarding_internal_port", rule.PortForwardingInternalPort)
	}
	return nil
}

func resourceNcloudPortForwardingRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Println("[DEBUG] resourceNcloudPortForwardingRuleExists")
	conn := meta.(*NcloudSdk).conn

	portForwardingExternalPort := d.Get("port_forwarding_external_port").(string)

	return hasPortForwardingRule(conn, portForwardingExternalPort)
}
func hasPortForwardingRule(conn *sdk.Conn, portForwardingExternalPort string) (bool, error) {
	rule, err := getPortForwardingRule(conn, portForwardingExternalPort)
	if err != nil {
		return false, err
	}

	if rule != nil {
		return true, nil
	}
	return false, nil
}

func resourceNcloudPortForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudPortForwardingRuleUpdate")
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func resourceNcloudPortForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] resourceNcloudPortForwardingRuleDelete")
	conn := meta.(*NcloudSdk).conn

	portForwardingConfigurationNo, err := getPortForwardingConfigurationNo(d, meta)
	if err != nil {
		return err
	}
	portForwardingExternalPort := d.Get("port_forwarding_external_port").(string)

	reqParams := &sdk.RequestDeletePortForwardingRules{
		PortForwardingConfigurationNo: portForwardingConfigurationNo,
		PortForwardingRuleList: []sdk.PortForwardingRule{
			{
				ServerInstanceNo:           d.Get("server_instance_no").(string),
				PortForwardingExternalPort: portForwardingExternalPort,
				PortForwardingInternalPort: d.Get("port_forwarding_internal_port").(string),
			},
		},
	}
	_, err = conn.DeletePortForwardingRules(reqParams)
	if err != nil {
		logErrorResponse("DeletePortForwardingRules", err, reqParams)
		return err
	}
	return nil
}

func getPortForwardingRuleList(conn *sdk.Conn) (*sdk.PortForwardingRuleList, error) {
	reqParams := &sdk.RequestPortForwardingRuleList{}
	resp, err := conn.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetPortForwardingRuleList", reqParams, resp.CommonResponse)

	return resp, nil
}

func getPortForwardingRule(conn *sdk.Conn, portForwardingExternalPort string) (*sdk.PortForwardingRule, error) {
	resp, err := getPortForwardingRuleList(conn)
	if err != nil {
		return nil, err
	}
	for _, rule := range resp.PortForwardingRuleList {
		if rule.PortForwardingExternalPort == portForwardingExternalPort {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("resource not found (portForwardingExternalPort) : %s", portForwardingExternalPort)
}
