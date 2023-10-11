package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func ResourceNcloudPortForwadingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudPortForwardingRuleCreate,
		Read:   resourceNcloudPortForwardingRuleRead,
		Update: resourceNcloudPortForwardingRuleUpdate,
		Delete: resourceNcloudPortForwardingRuleDelete,
		Exists: resourceNcloudPortForwardingRuleExists,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(conn.DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(conn.DefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"server_instance_no": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Server instance number for which port forwarding is set",
			},
			"port_forwarding_external_port": {
				Type:             schema.TypeInt,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1024, 65534)),
				Description:      "External port for port forwarding",
			},
			"port_forwarding_internal_port": {
				Type:             schema.TypeInt,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{22, 3389})), // [Linux : 22 |Windows : 3389]
				Description:      "Internal port for port forwarding. Only the following ports are available. [Linux: `22` | Windows: `3389`]",
			},
			"port_forwarding_configuration_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Port forwarding configuration number.",
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
	config := meta.(*conn.ProviderConfig)

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
	if err != nil {
		return err
	}

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
		LogCommonRequest("AddPortForwardingRules", reqParams)
		resp, err = config.Client.Server.V2Api.AddPortForwardingRules(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if ContainsInStringList(errBody.ReturnCode, []string{ApiErrorUnknown, ApiErrorPortForwardingObjectInOperation}) {
				LogErrorResponse("retry AddPortForwardingRules", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		LogResponse("AddPortForwardingRules", resp)
		return nil
	})

	if err != nil {
		LogErrorResponse("AddPortForwardingRules", err, reqParams)
		return err
	}
	d.SetId(newPortForwardingRuleId)
	return resourceNcloudPortForwardingRuleRead(d, meta)
}

func resourceNcloudPortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conn.ProviderConfig).Client

	_, zoneNo, portForwardingExternalPort := ParsePortForwardingRuleId(d.Id())
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
		d.Set("port_forwarding_configuration_no", portForwardingRule.PortForwardingConfigurationNo)

		if zone := zone.FlattenZone(portForwardingRule.Zone); zone["zone_code"] != nil {
			d.Set("zone", zone["zone_code"])
		}
	} else {
		log.Printf("unable to find resource: %s", d.Id())
		d.SetId("") // resource not found
	}

	return nil
}

func resourceNcloudPortForwardingRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*conn.ProviderConfig)

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
	client := meta.(*conn.ProviderConfig).Client

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
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		LogCommonRequest("DeletePortForwardingRules", reqParams)
		resp, err = client.Server.V2Api.DeletePortForwardingRules(reqParams)
		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if ContainsInStringList(errBody.ReturnCode, []string{ApiErrorUnknown, ApiErrorPortForwardingObjectInOperation}) {
				LogErrorResponse("retry DeletePortForwardingRules", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		LogResponse("DeletePortForwardingRules", resp)
		return nil
	})

	if err != nil {
		LogErrorResponse("DeletePortForwardingRules", err, reqParams)
		return err
	}
	d.SetId("")
	return nil
}

func PortForwardingRuleId(portForwardingConfigurationNo string, zonNo string, portForwardingExternalPort int32) string {
	return fmt.Sprintf("%s:%s:%d", portForwardingConfigurationNo, zonNo, portForwardingExternalPort)
}

func ParsePortForwardingRuleId(id string) (portForwardingConfigurationNo string, zoneNo string, portForwardingExternalPort int32) {
	arr := strings.Split(id, ":")

	portForwardingConfigurationNo, zoneNo = arr[0], arr[1]
	tmp, _ := strconv.Atoi(arr[2])
	return portForwardingConfigurationNo, zoneNo, int32(tmp)
}

func getPortForwardingConfigurationNo(d *schema.ResourceData, meta interface{}) (string, error) {
	config := meta.(*conn.ProviderConfig)

	paramPortForwardingConfigurationNo, ok := d.GetOk("port_forwarding_configuration_no")
	var portForwardingConfigurationNo string
	if ok {
		portForwardingConfigurationNo = paramPortForwardingConfigurationNo.(string)
	} else {
		resp, err := getPortForwardingConfigurationList(d, config)
		if err != nil {
			return "", err
		}
		portForwardingConfigurationNo = ncloud.StringValue(resp.PortForwardingConfigurationList[0].PortForwardingConfigurationNo)
	}
	return portForwardingConfigurationNo, nil
}

func getPortForwardingConfigurationList(d *schema.ResourceData, config *conn.ProviderConfig) (*server.GetPortForwardingConfigurationListResponse, error) {
	reqParams := &server.GetPortForwardingConfigurationListRequest{
		RegionNo:             ncloud.String(config.RegionNo),
		ServerInstanceNoList: []*string{ncloud.String(d.Get("server_instance_no").(string))},
	}
	LogCommonRequest("GetPortForwardingConfigurationList", reqParams)
	resp, err := config.Client.Server.V2Api.GetPortForwardingConfigurationList(reqParams)
	if err != nil {
		LogErrorResponse("GetPortForwardingConfigurationList", err, reqParams)
		return nil, err
	}
	LogCommonResponse("GetPortForwardingConfigurationList", GetCommonResponse(resp))

	return resp, nil
}

func getPortForwardingRuleList(client *conn.NcloudAPIClient, zoneNo string) (*server.GetPortForwardingRuleListResponse, error) {
	reqParams := &server.GetPortForwardingRuleListRequest{
		ZoneNo: ncloud.String(zoneNo),
	}
	LogCommonRequest("GetPortForwardingRuleList", reqParams)
	resp, err := client.Server.V2Api.GetPortForwardingRuleList(reqParams)
	if err != nil {
		LogErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return nil, err
	}
	LogCommonResponse("GetPortForwardingRuleList", GetCommonResponse(resp))

	return resp, nil
}

func GetPortForwardingRule(client *conn.NcloudAPIClient, zoneNo string, portForwardingExternalPort int32) (*server.PortForwardingRule, error) {
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

func hasPortForwardingRule(client *conn.NcloudAPIClient, zoneNo string, portForwardingExternalPort int32) (bool, error) {
	rule, _ := GetPortForwardingRule(client, zoneNo, portForwardingExternalPort)
	if rule != nil {
		return true, nil
	}
	return false, nil
}
