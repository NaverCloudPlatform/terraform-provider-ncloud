package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func dataSourceNcloudPortForwardingRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPortForwardingRuleRead,

		Schema: map[string]*schema.Schema{
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInternetLineTypeCode,
				Description:  "Internet line code. PUBLC(Public), GLBL(Global)",
			},
			"region_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region number.",
			},
			"zone_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone number.",
			},
			"server_instance_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Server instance number",
			},
			"port_forwarding_internal_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Port forwarding internal port.",
			},
			"port_forwarding_external_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Port forwarding external port.",
			},

			"port_forwarding_configuration_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port forwarding configuration number.",
			},
			"port_forwarding_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port forwarding public ip",
			},
		},
	}
}

func dataSourceNcloudPortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := &sdk.RequestPortForwardingRuleList{
		InternetLineTypeCode: d.Get("internet_line_type_code").(string),
		RegionNo:             parseRegionNoParameter(conn, d),
		ZoneNo:               d.Get("zone_no").(string),
	}

	resp, err := conn.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return err
	}
	logCommonResponse("GetPortForwardingRuleList", reqParams, resp.CommonResponse)

	allPortForwardingRules := resp.PortForwardingRuleList
	var filteredPortForwardingRuleList []sdk.PortForwardingRule
	var portForwardingRule sdk.PortForwardingRule

	filterServerInstanceNo, filterServerInstanceNoOk := d.GetOk("server_instance_no")
	filterInternalPort, filterInternalPortOk := d.GetOk("port_forwarding_internal_port")
	filterExternalPort, filterExternalPortOk := d.GetOk("port_forwarding_external_port")
	if filterServerInstanceNoOk || filterInternalPortOk || filterExternalPortOk {
		for _, portForwardingRule := range allPortForwardingRules {
			if filterServerInstanceNoOk && portForwardingRule.ServerInstanceNo == filterServerInstanceNo {
				filteredPortForwardingRuleList = append(filteredPortForwardingRuleList, portForwardingRule)
			} else if filterInternalPortOk && portForwardingRule.PortForwardingInternalPort == filterInternalPort {
				filteredPortForwardingRuleList = append(filteredPortForwardingRuleList, portForwardingRule)
			} else if filterExternalPortOk && portForwardingRule.PortForwardingExternalPort == filterExternalPort {
				filteredPortForwardingRuleList = append(filteredPortForwardingRuleList, portForwardingRule)
			}
		}
	} else {
		filteredPortForwardingRuleList = allPortForwardingRules[:]
	}

	if len(filteredPortForwardingRuleList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	portForwardingRule = filteredPortForwardingRuleList[0]

	return portForwardingRuleAttributes(d, resp.PortForwardingConfigurationNo, resp.PortForwardingPublicIp, portForwardingRule)
}

func portForwardingRuleAttributes(d *schema.ResourceData, portForwardingConfigurationNo int, portForwardingPublicIp string, rule sdk.PortForwardingRule) error {

	d.SetId(strconv.Itoa(portForwardingConfigurationNo))
	d.Set("port_forwarding_configuration_no", portForwardingConfigurationNo)
	d.Set("port_forwarding_public_ip", portForwardingPublicIp)
	d.Set("server_instance_no", rule.ServerInstanceNo)
	d.Set("port_forwarding_external_port", rule.PortForwardingExternalPort)
	d.Set("port_forwarding_internal_port", rule.PortForwardingInternalPort)

	return nil
}
