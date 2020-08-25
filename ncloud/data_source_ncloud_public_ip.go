package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceNcloudPublicIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPublicIpRead,

		Schema: map[string]*schema.Schema{
			"instance_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of public IP instance numbers to get.",
			},
			"internet_line_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
				Description:  "Internet line type code. `PUBLC` (Public), `GLBL` (Global)",
			},
			"is_associated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the public IP address is associated or not.",
			},
			"list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of public IP addresses to get.",
			},
			"search_filter_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "use filter instead",
				Description: "`publicIp` (Public IP) | `associatedServerName` (Associated server name)",
			},
			"search_filter_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "use filter instead",
				Description: "Filter value to search",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "use region attribute of provider instead",
				Description: "Region code. Get available values using the `data ncloud_regions`.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone code. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.",
			},
			"sorted_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This attribute no longer support",
				Description: "The column based on which you want to sort the list.",
			},
			"sorting_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This attribute no longer support",
				Description: "Sorting order of the list. `ascending` (Ascending) | `descending` (Descending) [case insensitive]. Default: `ascending` Ascending",
			},
			"filter": dataSourceFiltersSchema(),

			"instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance number",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP description",
			},
			"instance_status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance status name",
			},
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance status",
			},
			"instance_operation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance operation",
			},
			"kind_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP kind type",
			},
			"server_instance": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Associated server instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_instance_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated server instance number",
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated server name",
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudPublicIpRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcPublicIpList(d, meta.(*ProviderConfig))
	} else {
		resources, err = getClassicPublicIpList(d, meta.(*ProviderConfig))
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceNcloudPublicIp().Schema)
	}

	if err := validateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMap(d, resources[0])

	return nil
}

func getClassicPublicIpList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	reqParams := &server.GetPublicIpInstanceListRequest{
		RegionNo:             &regionNo,
		ZoneNo:               StringPtrOrNil(d.GetOk("zone")),
		InternetLineTypeCode: StringPtrOrNil(d.GetOk("internet_line_type")),
	}

	if isAssociated, ok := d.GetOk("is_associated"); ok {
		reqParams.IsAssociated = ncloud.Bool(isAssociated.(bool))
	}

	if instanceNoList, ok := d.GetOk("instance_no_list"); ok {
		reqParams.PublicIpInstanceNoList = expandStringInterfaceList(instanceNoList.([]interface{}))
	}

	if publicIPList, ok := d.GetOk("list"); ok {
		reqParams.PublicIpList = expandStringInterfaceList(publicIPList.([]interface{}))
	}

	logCommonRequest("getClassicPublicIpList", reqParams)
	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getClassicPublicIpList", err, reqParams)
		return nil, err
	}
	logCommonResponse("getClassicPublicIpList", GetCommonResponse(resp))

	resources := []map[string]interface{}{}

	for _, r := range resp.PublicIpInstanceList {
		instance := map[string]interface{}{
			"id":                   *r.PublicIpInstanceNo,
			"instance_no":          *r.PublicIpInstanceNo,
			"public_ip":            *r.PublicIp,
			"description":          *r.PublicIpDescription,
			"instance_status_name": *r.PublicIpInstanceStatusName,
		}

		if m := flattenCommonCode(r.InternetLineType); m["code"] != nil {
			instance["internet_line_type"] = m["code"]
		}

		if m := flattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
			instance["instance_status"] = m["code"]
		}

		if m := flattenCommonCode(r.PublicIpInstanceOperation); m["code"] != nil {
			instance["instance_operation"] = m["code"]
		}

		if m := flattenCommonCode(r.PublicIpKindType); m["code"] != nil {
			instance["kind_type"] = m["code"]
		}

		if m := flattenCommonCode(r.Zone); m["code"] != nil {
			instance["zone"] = m["code"]
		}

		if serverInstance := r.ServerInstanceAssociatedWithPublicIp; serverInstance != nil {
			mapping := map[string]interface{}{
				"server_instance_no": ncloud.StringValue(serverInstance.ServerInstanceNo),
				"server_name":        ncloud.StringValue(serverInstance.ServerName),
			}

			instance["server_instance"] = mapping
		}

		resources = append(resources, instance)
	}

	return resources, nil
}

func getVpcPublicIpList(d *schema.ResourceData, config *ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionCode := config.RegionCode

	reqParams := &vserver.GetPublicIpInstanceListRequest{
		RegionCode: &regionCode,
	}

	if isAssociated, ok := d.GetOk("is_associated"); ok {
		reqParams.IsAssociated = ncloud.Bool(isAssociated.(bool))
	}

	if instanceNoList, ok := d.GetOk("instance_no_list"); ok {
		reqParams.PublicIpInstanceNoList = expandStringInterfaceList(instanceNoList.([]interface{}))
	}

	logCommonRequest("getVpcPublicIpList", reqParams)
	resp, err := client.vserver.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getVpcPublicIpList", err, reqParams)
		return nil, err
	}
	logCommonResponse("getVpcPublicIpList", GetCommonResponse(resp))

	resources := []map[string]interface{}{}

	for _, r := range resp.PublicIpInstanceList {
		instance := map[string]interface{}{
			"id":                   *r.PublicIpInstanceNo,
			"instance_no":          *r.PublicIpInstanceNo,
			"public_ip":            *r.PublicIp,
			"description":          *r.PublicIpDescription,
			"instance_status_name": *r.PublicIpInstanceStatusName,
		}

		if m := flattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
			instance["instance_status"] = m["code"]
		}

		if m := flattenCommonCode(r.PublicIpInstanceOperation); m["code"] != nil {
			instance["instance_operation"] = m["code"]
		}

		if r.ServerInstanceNo != nil && r.ServerName != nil {
			mapping := map[string]interface{}{
				"server_instance_no": ncloud.StringValue(r.ServerInstanceNo),
				"server_name":        ncloud.StringValue(r.ServerName),
			}

			instance["server_instance"] = mapping
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
