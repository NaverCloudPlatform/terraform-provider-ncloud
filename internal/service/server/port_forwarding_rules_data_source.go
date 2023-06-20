package server

import (
	"fmt"
	"strconv"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func init() {
	RegisterDataSource("ncloud_port_forwarding_rules", dataSourceNcloudPortForwardingRules())
}

func dataSourceNcloudPortForwardingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPortForwardingRulesRead,

		Schema: map[string]*schema.Schema{
			// Deprecated
			"internet_line_type_code": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
				Description:      "Internet line code. PUBLC(Public), GLBL(Global)",
				Deprecated:       "This parameter is no longer used.",
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
	config := meta.(*ProviderConfig)
	client := config.Client

	if config.SupportVPC {
		return NotSupportVpc("data source `ncloud_port_forwarding_rules`")
	}

	regionNo, err := ParseRegionNoParameter(d)
	if err != nil {
		return err
	}
	zoneNo, err := zone.ParseZoneNoParameter(config, d)
	if err != nil {
		return err
	}
	reqParams := &server.GetPortForwardingRuleListRequest{
		RegionNo: regionNo,
		ZoneNo:   zoneNo,
	}

	LogCommonRequest("GetPortForwardingRuleList", reqParams)
	resp, err := client.Server.V2Api.GetPortForwardingRuleList(reqParams)
	if err != nil {
		LogErrorResponse("GetPortForwardingRuleList", err, reqParams)
		return err
	}
	LogCommonResponse("GetPortForwardingRuleList", GetCommonResponse(resp))

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
	return portForwardingRulesAttributes(d, resp.PortForwardingRuleList[0].PortForwardingConfigurationNo, filteredPortForwardingRuleList)
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
		return WriteToFile(output.(string), s)
	}

	return nil
}
