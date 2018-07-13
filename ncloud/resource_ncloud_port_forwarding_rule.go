package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"strings"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server instance number for which port forwarding is set",
			},
			"port_forwarding_external_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1024, 65534),
				Description:  "External port for port forwarding",
			},
			"port_forwarding_internal_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIncludeValues([]string{"22", "3389"}), // [Linux : 22 |Windows : 3389]
				Description:  "Internal port for port forwarding. Only the following ports are available. [Linux: `22` | Windows: `3389`]",
			},
			"port_forwarding_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port forwarding Public IP",
			},
			"zone": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        zoneSchemaResource,
				Description: "Zone info",
			},
		},
	}
}

func resourceNcloudPortForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
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
	d.SetId(PortForwardingRuleId(portForwardingConfigurationNo, resp.Zone.ZoneNo, portForwardingExternalPort))
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func PortForwardingRuleId(portForwardingConfigurationNo string, zonNo string, portForwardingExternalPort string) string {
	return fmt.Sprintf("%s:%s:%s", portForwardingConfigurationNo, zonNo, portForwardingExternalPort)
}

func parsePortForwardingRuleId(id string) (portForwardingConfigurationNo string, zoneNo string, portForwardingExternalPort string) {
	arr := strings.Split(id, ":")
	portForwardingConfigurationNo, zoneNo, portForwardingExternalPort = arr[0], arr[1], arr[2]
	return portForwardingConfigurationNo, zoneNo, portForwardingExternalPort
}

func resourceNcloudPortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	_, zoneNo, portForwardingExternalPort := parsePortForwardingRuleId(d.Id())
	resp, err := getPortForwardingRuleList(conn, zoneNo)
	if err != nil {
		return err
	}

	var portForwardingRule *sdk.PortForwardingRule
	for _, rule := range resp.PortForwardingRuleList {
		if rule.PortForwardingExternalPort == portForwardingExternalPort {
			portForwardingRule = &rule
			break
		}
	}
	if portForwardingRule != nil {
		d.Set("port_forwarding_public_ip", resp.PortForwardingPublicIp)
		d.Set("server_instance_no", portForwardingRule.ServerInstanceNo)
		d.Set("port_forwarding_external_port", portForwardingRule.PortForwardingExternalPort)
		d.Set("port_forwarding_internal_port", portForwardingRule.PortForwardingInternalPort)
		d.Set("zone", setZone(resp.Zone))

	}
	return nil
}

func resourceNcloudPortForwardingRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*NcloudSdk).conn

	zoneNo, err := getServerZoneNo(conn, d.Get("server_instance_no").(string))
	if err != nil {
		return false, err
	}
	portForwardingExternalPort := d.Get("port_forwarding_external_port").(string)

	return hasPortForwardingRule(conn, zoneNo, portForwardingExternalPort)
}

func resourceNcloudPortForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func resourceNcloudPortForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
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

func getPortForwardingConfigurationNo(d *schema.ResourceData, meta interface{}) (string, error) {
	conn := meta.(*NcloudSdk).conn
	paramPortForwardingConfigurationNo, ok := d.GetOk("port_forwarding_configuration_no")
	var portForwardingConfigurationNo string
	if ok {
		portForwardingConfigurationNo = paramPortForwardingConfigurationNo.(string)
	} else {
		zoneNo, err := getServerZoneNo(conn, d.Get("server_instance_no").(string))
		if err != nil {
			return "", err
		}
		resp, err := getPortForwardingRuleList(conn, zoneNo)
		if err != nil {
			return "", err
		}
		portForwardingConfigurationNo = strconv.Itoa(resp.PortForwardingConfigurationNo)
	}
	return portForwardingConfigurationNo, nil
}
func getServerZoneNo(conn *sdk.Conn, serverInstanceNo string) (string, error) {
	serverInstance, err := getServerInstance(conn, serverInstanceNo)
	if err != nil {
		return "", err
	}
	return serverInstance.Zone.ZoneNo, nil
}

func getPortForwardingRuleList(conn *sdk.Conn, zoneNo string) (*sdk.PortForwardingRuleList, error) {
	reqParams := &sdk.RequestPortForwardingRuleList{
		ZoneNo: zoneNo,
	}
	resp, err := conn.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetPortForwardingRuleList", reqParams, resp.CommonResponse)

	return resp, nil
}

func getPortForwardingRule(conn *sdk.Conn, zoneNo string, portForwardingExternalPort string) (*sdk.PortForwardingRule, error) {
	resp, err := getPortForwardingRuleList(conn, zoneNo)
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

func hasPortForwardingRule(conn *sdk.Conn, zoneNo string, portForwardingExternalPort string) (bool, error) {
	rule, err := getPortForwardingRule(conn, zoneNo, portForwardingExternalPort)
	if err != nil {
		return false, err
	}

	if rule != nil {
		return true, nil
	}
	return false, nil
}
