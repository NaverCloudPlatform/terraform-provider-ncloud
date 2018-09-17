package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
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
	client := meta.(*NcloudAPIClient)

	_, zoneNoOk := d.GetOk("zone_no")
	_, zoneCodeOk := d.GetOk("zone_code")
	if !zoneNoOk && !zoneCodeOk {
		return fmt.Errorf("required to select one among two parameters: `zone_no` and `zone_code`")
	}

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := &server.GetPortForwardingRuleListRequest{
		InternetLineTypeCode: StringPtrOrNil(d.GetOk("internet_line_type_code")),
		RegionNo:             regionNo,
		ZoneNo:               zoneNo,
	}

	resp, err := client.server.V2Api.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return err
	}
	logCommonResponse("GetPortForwardingRuleList", reqParams, GetCommonResponse(resp))

	allPortForwardingRules := resp.PortForwardingRuleList
	var filteredPortForwardingRuleList []*server.PortForwardingRule

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

func portForwardingRulesAttributes(d *schema.ResourceData, portForwardingConfigurationNo *string, portForwardingRuleList []*server.PortForwardingRule) error {
	var s []map[string]interface{}

	d.SetId(*portForwardingConfigurationNo)
	d.Set("port_forwarding_configuration_no", portForwardingConfigurationNo)

	for _, rule := range portForwardingRuleList {
		mapping := map[string]interface{}{
			"server_instance_no":            *rule.ServerInstance.ServerInstanceNo,
			"port_forwarding_external_port": strconv.Itoa(int(*rule.PortForwardingExternalPort)),
			"port_forwarding_internal_port": strconv.Itoa(int(*rule.PortForwardingInternalPort)),
			"port_forwarding_public_ip":     *rule.ServerInstance.PortForwardingPublicIp,
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
