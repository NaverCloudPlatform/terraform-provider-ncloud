package ncloud

import (
	"fmt"
	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudPortForwardingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPortForwardingRulesRead,

		Schema: map[string]*schema.Schema{
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInternetLineTypeCode,
				Description:  "Internet line code. PUBLC(Public), GLBL(Global)",
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
			"zone_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone code",
				ConflictsWith: []string{"zone_no"},
			},
			"zone_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone number",
				ConflictsWith: []string{"zone_code"},
			},
			"port_forwarding_internal_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port forwarding internal port.",
			},
			"port_forwarding_configuration_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port forwarding configuration number.",
			},
			"port_forwarding_rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Port forwarding rule list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_instance_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server instance number",
						},
						"port_forwarding_external_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port forwarding external port.",
						},
						"port_forwarding_internal_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port forwarding internal port.",
						},
						"port_forwarding_public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port forwarding public ip",
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNcloudPortForwardingRulesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	regionNo, err := parseRegionNoParameter(conn, d)
	if err != nil {
		return err
	}
	zoneNo, err := parseZoneNoParameter(conn, d)
	if err != nil {
		return err
	}
	reqParams := &sdk.RequestPortForwardingRuleList{
		InternetLineTypeCode: d.Get("internet_line_type_code").(string),
		RegionNo:             regionNo,
		ZoneNo:               zoneNo,
	}

	resp, err := conn.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return err
	}
	logCommonResponse("GetPortForwardingRuleList", reqParams, resp.CommonResponse)

	allPortForwardingRules := resp.PortForwardingRuleList
	var filteredPortForwardingRuleList []sdk.PortForwardingRule

	filterInternalPort, filterInternalPortOk := d.GetOk("port_forwarding_internal_port")
	if filterInternalPortOk {
		for _, portForwardingRule := range allPortForwardingRules {
			if portForwardingRule.PortForwardingInternalPort == filterInternalPort {
				filteredPortForwardingRuleList = append(filteredPortForwardingRuleList, portForwardingRule)
			}
		}
	} else {
		filteredPortForwardingRuleList = allPortForwardingRules[:]
	}

	if len(filteredPortForwardingRuleList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}
	return portForwardingRulesAttributes(d, resp.PortForwardingConfigurationNo, filteredPortForwardingRuleList)
}

func portForwardingRulesAttributes(d *schema.ResourceData, portForwardingConfigurationNo int, portForwardingRuleList []sdk.PortForwardingRule) error {
	var s []map[string]interface{}

	d.SetId(strconv.Itoa(portForwardingConfigurationNo))
	d.Set("port_forwarding_configuration_no", portForwardingConfigurationNo)

	for _, rule := range portForwardingRuleList {
		mapping := map[string]interface{}{
			"server_instance_no":            rule.ServerInstanceNo,
			"port_forwarding_external_port": rule.PortForwardingExternalPort,
			"port_forwarding_internal_port": rule.PortForwardingInternalPort,
			"port_forwarding_public_ip":     rule.PortForwardingPublicIp,
		}
		s = append(s, mapping)
	}

	if err := d.Set("port_forwarding_rule_list", s); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
