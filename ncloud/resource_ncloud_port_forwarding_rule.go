package ncloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterResource("ncloud_port_forwarding_rule", resourceNcloudPortForwadingRule())
}

func resourceNcloudPortForwadingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudPortForwardingRuleCreate,
		Read:   resourceNcloudPortForwardingRuleRead,
		Update: resourceNcloudPortForwardingRuleUpdate,
		Delete: resourceNcloudPortForwardingRuleDelete,
		Exists: resourceNcloudPortForwardingRuleExists,

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
				ValidateFunc: validation.IntBetween(1024, 65534),
				Description:  "External port for port forwarding",
			},
			"port_forwarding_internal_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{22, 3389}), // [Linux : 22 |Windows : 3389]
				Description:  "Internal port for port forwarding. Only the following ports are available. [Linux: `22` | Windows: `3389`]",
			},
			"port_forwarding_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port forwarding Public IP",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone info",
			},
		},
	}
}

func resourceNcloudPortForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

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

	serverInstanceNo := d.Get("server_instance_no").(string)
	zoneNo, err := getServerZoneNo(config, serverInstanceNo)
	newPortForwardingRuleId := PortForwardingRuleId(portForwardingConfigurationNo, zoneNo, portForwardingExternalPort)
	log.Printf("[DEBUG] AddPortForwardingRules newPortForwardingRuleId: %s", newPortForwardingRuleId)

	reqParams := &server.AddPortForwardingRulesRequest{
		PortForwardingConfigurationNo: ncloud.String(portForwardingConfigurationNo),
		PortForwardingRuleList: []*server.PortForwardingRuleParameter{
			{
				ServerInstanceNo:           ncloud.String(serverInstanceNo),
				PortForwardingExternalPort: ncloud.Int32(portForwardingExternalPort),
				PortForwardingInternalPort: ncloud.Int32(portForwardingInternalPort),
			},
		},
	}

	var resp *server.AddPortForwardingRulesResponse
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		logCommonRequest("AddPortForwardingRules", reqParams)
		resp, err = config.Client.server.V2Api.AddPortForwardingRules(reqParams)

		if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorPortForwardingObjectInOperation}) {
			logErrorResponse("retry AddPortForwardingRules", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}
		logCommonResponse("AddPortForwardingRules", GetCommonResponse(resp))

		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("AddPortForwardingRules", err, reqParams)
		return err
	}
	d.SetId(newPortForwardingRuleId)
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func resourceNcloudPortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	_, zoneNo, portForwardingExternalPort := parsePortForwardingRuleId(d.Id())
	resp, err := getPortForwardingRuleList(client, zoneNo)
	if err != nil {
		return err
	}

	var portForwardingRule *server.PortForwardingRule
	for _, rule := range resp.PortForwardingRuleList {
		if ncloud.Int32Value(rule.PortForwardingExternalPort) == portForwardingExternalPort {
			portForwardingRule = rule
			break
		}
	}
	if portForwardingRule != nil {
		d.Set("port_forwarding_public_ip", portForwardingRule.ServerInstance.PortForwardingPublicIp)
		d.Set("server_instance_no", portForwardingRule.ServerInstance.ServerInstanceNo)
		d.Set("port_forwarding_external_port", portForwardingRule.PortForwardingExternalPort)
		d.Set("port_forwarding_internal_port", portForwardingRule.PortForwardingInternalPort)

		if zone := flattenZone(resp.Zone); zone["zone_code"] != nil {
			d.Set("zone", zone["zone_code"])
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudPortForwardingRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*ProviderConfig)

	zoneNo, err := getServerZoneNo(config, d.Get("server_instance_no").(string))
	if err != nil {
		return false, err
	}
	var portForwardingExternalPort int32
	if v, ok := d.GetOk("port_forwarding_external_port"); ok {
		portForwardingExternalPort = int32(v.(int))
	}
	return hasPortForwardingRule(config.Client, zoneNo, portForwardingExternalPort)
}

func resourceNcloudPortForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func resourceNcloudPortForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

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

	serverInstanceNo := d.Get("server_instance_no").(string)
	reqParams := &server.DeletePortForwardingRulesRequest{
		PortForwardingConfigurationNo: ncloud.String(portForwardingConfigurationNo),
		PortForwardingRuleList: []*server.PortForwardingRuleParameter{
			{
				ServerInstanceNo:           ncloud.String(serverInstanceNo),
				PortForwardingExternalPort: ncloud.Int32(portForwardingExternalPort),
				PortForwardingInternalPort: ncloud.Int32(portForwardingInternalPort),
			},
		},
	}

	var resp *server.DeletePortForwardingRulesResponse
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error

		logCommonRequest("DeletePortForwardingRules", reqParams)

		resp, err = client.server.V2Api.DeletePortForwardingRules(reqParams)
		log.Printf("=================> DeletePortForwardingRules resp: %#v, err: %#v", resp, err)
		if err == nil && resp == nil {
			return resource.NonRetryableError(err)
		}
		if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorPortForwardingObjectInOperation}) {
			logErrorResponse("DeletePortForwardingRules Retry", err, reqParams)
			time.Sleep(time.Second * 5)
			return resource.RetryableError(err)
		}
		logCommonResponse("DeletePortForwardingRules", GetCommonResponse(resp))
		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("DeletePortForwardingRules", err, reqParams)
		return err
	}
	d.SetId("")
	return nil
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

func getPortForwardingConfigurationNo(d *schema.ResourceData, meta interface{}) (string, error) {
	config := meta.(*ProviderConfig)

	paramPortForwardingConfigurationNo, ok := d.GetOk("port_forwarding_configuration_no")
	var portForwardingConfigurationNo string
	if ok {
		portForwardingConfigurationNo = paramPortForwardingConfigurationNo.(string)
	} else {
		zoneNo, err := getServerZoneNo(config, d.Get("server_instance_no").(string))
		if err != nil {
			return "", err
		}
		resp, err := getPortForwardingRuleList(config.Client, zoneNo)
		if err != nil {
			return "", err
		}
		portForwardingConfigurationNo = ncloud.StringValue(resp.PortForwardingConfigurationNo)
	}
	return portForwardingConfigurationNo, nil
}

func getPortForwardingRuleList(client *NcloudAPIClient, zoneNo string) (*server.GetPortForwardingRuleListResponse, error) {
	reqParams := &server.GetPortForwardingRuleListRequest{
		ZoneNo: ncloud.String(zoneNo),
	}
	logCommonRequest("GetPortForwardingRuleList", reqParams)
	resp, err := client.server.V2Api.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetPortForwardingRuleList", GetCommonResponse(resp))

	return resp, nil
}

func getPortForwardingRule(client *NcloudAPIClient, zoneNo string, portForwardingExternalPort int32) (*server.PortForwardingRule, error) {
	resp, err := getPortForwardingRuleList(client, zoneNo)
	if err != nil {
		return nil, err
	}
	for _, rule := range resp.PortForwardingRuleList {
		if portForwardingExternalPort == ncloud.Int32Value(rule.PortForwardingExternalPort) {
			return rule, nil
		}
	}
	return nil, nil
}

func hasPortForwardingRule(client *NcloudAPIClient, zoneNo string, portForwardingExternalPort int32) (bool, error) {
	rule, _ := getPortForwardingRule(client, zoneNo, portForwardingExternalPort)
	if rule != nil {
		return true, nil
	}
	return false, nil
}
