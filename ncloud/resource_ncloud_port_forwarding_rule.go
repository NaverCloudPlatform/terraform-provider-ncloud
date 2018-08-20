package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
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
				Type:         schema.TypeInt,
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
	client := meta.(*NcloudAPIClient)

	portForwardingConfigurationNo, err := getPortForwardingConfigurationNo(d, meta)
	if err != nil {
		return err
	}

	var portForwardingExternalPort int32
	if v, ok := d.GetOk("port_forwarding_external_port"); ok {
		portForwardingExternalPort = int32(v.(int))
	}
	var portForwardingInternalPort int32
	if v, ok := d.GetOk("port_forwarding_internal_port"); ok {
		portForwardingInternalPort = int32(v.(int))
	}

	reqParams := &server.AddPortForwardingRulesRequest{
		PortForwardingConfigurationNo: ncloud.String(portForwardingConfigurationNo),
		PortForwardingRuleList: []*server.PortForwardingRuleParameter{
			{
				ServerInstanceNo:           ncloud.String(d.Get("server_instance_no").(string)),
				PortForwardingExternalPort: ncloud.Int32(portForwardingExternalPort),
				PortForwardingInternalPort: ncloud.Int32(portForwardingInternalPort),
			},
		},
	}
	resp, err := client.server.V2Api.AddPortForwardingRules(reqParams)
	if err != nil {
		logErrorResponse("AddPortForwardingRules", err, reqParams)
		return err
	}
	logCommonResponse("AddPortForwardingRules", reqParams, GetCommonResponse(resp))
	d.SetId(PortForwardingRuleId(portForwardingConfigurationNo, *resp.Zone.ZoneNo, portForwardingExternalPort))
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func PortForwardingRuleId(portForwardingConfigurationNo string, zonNo string, portForwardingExternalPort int32) string {
	return fmt.Sprintf("%s:%s:%d", portForwardingConfigurationNo, zonNo, portForwardingExternalPort)
}

func parsePortForwardingRuleId(id string) (portForwardingConfigurationNo string, zoneNo string, portForwardingExternalPort int32) {
	arr := strings.Split(id, ":")

	portForwardingConfigurationNo, zoneNo = arr[0], arr[1]
	tmp, _ := strconv.Atoi(arr[2])
	return portForwardingConfigurationNo, zoneNo, int32(tmp)
}

func resourceNcloudPortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	_, zoneNo, portForwardingExternalPort := parsePortForwardingRuleId(d.Id())
	resp, err := getPortForwardingRuleList(client, zoneNo)
	if err != nil {
		return err
	}

	var portForwardingRule *server.PortForwardingRule
	for _, rule := range resp.PortForwardingRuleList {
		if *rule.PortForwardingExternalPort == portForwardingExternalPort {
			portForwardingRule = rule
			break
		}
	}
	if portForwardingRule != nil {
		d.Set("port_forwarding_public_ip", resp.PortForwardingPublicIp)
		d.Set("server_instance_no", portForwardingRule.ServerInstance.ServerInstanceNo)
		d.Set("port_forwarding_external_port", portForwardingRule.PortForwardingExternalPort)
		d.Set("port_forwarding_internal_port", portForwardingRule.PortForwardingInternalPort)
		d.Set("zone", setZone(resp.Zone))

	}
	return nil
}

func resourceNcloudPortForwardingRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*NcloudAPIClient)

	zoneNo, err := getServerZoneNo(client, d.Get("server_instance_no").(string))
	if err != nil {
		return false, err
	}
	var portForwardingExternalPort int32
	if v, ok := d.GetOk("port_forwarding_external_port"); ok {
		portForwardingExternalPort = int32(v.(int))
	}
	return hasPortForwardingRule(client, zoneNo, portForwardingExternalPort)
}

func resourceNcloudPortForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func resourceNcloudPortForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	portForwardingConfigurationNo, err := getPortForwardingConfigurationNo(d, meta)
	if err != nil {
		return err
	}
	var portForwardingExternalPort int32
	if v, ok := d.GetOk("port_forwarding_external_port"); ok {
		portForwardingExternalPort = int32(v.(int))
	}
	var portForwardingInternalPort int32
	if v, ok := d.GetOk("port_forwarding_internal_port"); ok {
		portForwardingInternalPort = int32(v.(int))
	}

	reqParams := &server.DeletePortForwardingRulesRequest{
		PortForwardingConfigurationNo: ncloud.String(portForwardingConfigurationNo),
		PortForwardingRuleList: []*server.PortForwardingRuleParameter{
			{
				ServerInstanceNo:           ncloud.String(d.Get("server_instance_no").(string)),
				PortForwardingExternalPort: ncloud.Int32(portForwardingExternalPort),
				PortForwardingInternalPort: ncloud.Int32(portForwardingInternalPort),
			},
		},
	}
	_, err = client.server.V2Api.DeletePortForwardingRules(reqParams)
	if err != nil {
		logErrorResponse("DeletePortForwardingRules", err, reqParams)
		return err
	}
	return nil
}

func getPortForwardingConfigurationNo(d *schema.ResourceData, meta interface{}) (string, error) {
	client := meta.(*NcloudAPIClient)
	paramPortForwardingConfigurationNo, ok := d.GetOk("port_forwarding_configuration_no")
	var portForwardingConfigurationNo string
	if ok {
		portForwardingConfigurationNo = paramPortForwardingConfigurationNo.(string)
	} else {
		zoneNo, err := getServerZoneNo(client, d.Get("server_instance_no").(string))
		if err != nil {
			return "", err
		}
		resp, err := getPortForwardingRuleList(client, zoneNo)
		if err != nil {
			return "", err
		}
		portForwardingConfigurationNo = *resp.PortForwardingConfigurationNo
	}
	return portForwardingConfigurationNo, nil
}
func getServerZoneNo(client *NcloudAPIClient, serverInstanceNo string) (string, error) {
	serverInstance, err := getServerInstance(client, serverInstanceNo)
	if err != nil {
		return "", err
	}
	return *serverInstance.Zone.ZoneNo, nil
}

func getPortForwardingRuleList(client *NcloudAPIClient, zoneNo string) (*server.GetPortForwardingRuleListResponse, error) {
	reqParams := &server.GetPortForwardingRuleListRequest{
		ZoneNo: ncloud.String(zoneNo),
	}
	resp, err := client.server.V2Api.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetPortForwardingRuleList", reqParams, GetCommonResponse(resp))

	return resp, nil
}

func getPortForwardingRule(client *NcloudAPIClient, zoneNo string, portForwardingExternalPort int32) (*server.PortForwardingRule, error) {
	resp, err := getPortForwardingRuleList(client, zoneNo)
	if err != nil {
		return nil, err
	}
	for _, rule := range resp.PortForwardingRuleList {
		if *rule.PortForwardingExternalPort == portForwardingExternalPort {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("resource not found (portForwardingExternalPort) : %d", portForwardingExternalPort)
}

func hasPortForwardingRule(client *NcloudAPIClient, zoneNo string, portForwardingExternalPort int32) (bool, error) {
	rule, err := getPortForwardingRule(client, zoneNo, portForwardingExternalPort)
	if err != nil {
		return false, err
	}

	if rule != nil {
		return true, nil
	}
	return false, nil
}
