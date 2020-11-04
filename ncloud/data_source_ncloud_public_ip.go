package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func init() {
	RegisterDataSource("ncloud_public_ip", dataSourceNcloudPublicIp())
}

func dataSourceNcloudPublicIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPublicIpRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internet_line_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLC", "GLBL"}, false),
			},
			"is_associated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter": dataSourceFiltersSchema(),

			"public_ip_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_instance": {
				Type:       schema.TypeMap,
				Computed:   true,
				Deprecated: "use 'server_instance_no' instead",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_instance_no": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "use 'server_instance_no' instead",
						},
						"server_name": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This field no longer support",
						},
					},
				},
			},
			"search_filter_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use 'filter' instead",
			},
			"search_filter_value": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use 'filter' instead",
			},
			"region": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use region attribute of provider instead",
			},
			"sorted_by": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This attribute no longer support",
			},
			"sorting_order": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This attribute no longer support",
			},
			"instance_no": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use 'id' instead",
			},
			"list": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "use 'filter' instead",
			},
			"instance_no_list": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "use 'filter' instead",
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

	if v, ok := d.GetOk("id"); ok {
		reqParams.PublicIpInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getClassicPublicIpList", reqParams)
	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getClassicPublicIpList", err, reqParams)
		return nil, err
	}
	logCommonResponse("getClassicPublicIpList", GetCommonResponse(resp))

	var resources []map[string]interface{}
	for _, r := range resp.PublicIpInstanceList {
		instance := map[string]interface{}{
			"id":                 *r.PublicIpInstanceNo,
			"instance_no":        *r.PublicIpInstanceNo,
			"public_ip_no":       *r.PublicIpInstanceNo,
			"public_ip":          *r.PublicIp,
			"description":        *r.PublicIpDescription,
			"server_instance_no": nil,
			"server_name":        nil,
		}

		if m := flattenCommonCode(r.InternetLineType); m["code"] != nil {
			instance["internet_line_type"] = m["code"]
		}

		if m := flattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
			instance["status"] = m["code"]
		}

		if m := flattenCommonCode(r.PublicIpKindType); m["code"] != nil {
			instance["kind_type"] = m["code"]
		}

		if m := flattenZone(r.Zone); m["zone_code"] != nil {
			instance["zone"] = m["zone_code"]
		}

		if serverInstance := r.ServerInstanceAssociatedWithPublicIp; serverInstance != nil {
			SetStringIfNotNilAndEmpty(instance, "server_instance_no", serverInstance.ServerInstanceNo)
			SetStringIfNotNilAndEmpty(instance, "server_name", serverInstance.ServerName)

			// Deprecated
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

	if v, ok := d.GetOk("is_associated"); ok {
		reqParams.IsAssociated = ncloud.Bool(v.(bool))
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.PublicIpInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	logCommonRequest("getVpcPublicIpList", reqParams)
	resp, err := client.vserver.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("getVpcPublicIpList", err, reqParams)
		return nil, err
	}
	logCommonResponse("getVpcPublicIpList", GetCommonResponse(resp))

	var resources []map[string]interface{}
	for _, r := range resp.PublicIpInstanceList {
		instance := map[string]interface{}{
			"id":                 *r.PublicIpInstanceNo,
			"public_ip_no":       *r.PublicIpInstanceNo,
			"public_ip":          *r.PublicIp,
			"description":        *r.PublicIpDescription,
			"server_instance_no": nil,
			"server_name":        nil,
		}

		SetStringIfNotNilAndEmpty(instance, "server_instance_no", r.ServerInstanceNo)
		SetStringIfNotNilAndEmpty(instance, "server_name", r.ServerName)

		if m := flattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
			instance["status"] = m["code"]
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
