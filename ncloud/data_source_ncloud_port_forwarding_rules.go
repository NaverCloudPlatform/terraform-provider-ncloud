package ncloud

import (
	"fmt"
	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceNcloudPortForwardingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPortForwardingRulesRead,

		Schema: map[string]*schema.Schema{
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line code. PUBLC(Public), GLBL(Global)",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone code. Get available values using the `data ncloud_zones`.",
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

	logCommonRequest("GetPortForwardingRuleList", reqParams)
	resp, err := client.server.V2Api.GetPortForwardingRuleList(reqParams)
	if err != nil {
		logErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return err
	}
	logCommonResponse("GetPortForwardingRuleList", GetCommonResponse(resp))

	allPortForwardingRules := resp.PortForwardingRuleList
	var filteredPortForwardingRuleList []*server.PortForwardingRule

	filterInternalPort, filterInternalPortOk := d.GetOk("port_forwarding_internal_port")
	if filterInternalPortOk {
		for _, portForwardingRule := range allPortForwardingRules {
			if filterInternalPort == strconv.Itoa(int(ncloud.Int32Value(portForwardingRule.PortForwardingInternalPort))) {
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

	d.SetId(ncloud.StringValue(portForwardingConfigurationNo))
	d.Set("port_forwarding_configuration_no", portForwardingConfigurationNo)

	for _, rule := range portForwardingRuleList {
		mapping := map[string]interface{}{
			"server_instance_no":            ncloud.StringValue(rule.ServerInstance.ServerInstanceNo),
			"port_forwarding_external_port": strconv.Itoa(int(ncloud.Int32Value(rule.PortForwardingExternalPort))),
			"port_forwarding_internal_port": strconv.Itoa(int(ncloud.Int32Value(rule.PortForwardingInternalPort))),
			"port_forwarding_public_ip":     ncloud.StringValue(rule.ServerInstance.PortForwardingPublicIp),
		}
		s = append(s, mapping)
	}

	if err := d.Set("port_forwarding_rule_list", s); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), s)
	}

	return nil
}
