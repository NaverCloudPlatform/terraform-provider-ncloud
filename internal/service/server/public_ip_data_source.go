package server

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/zone"
)

func DataSourceNcloudPublicIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPublicIpRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Deprecated
			"internet_line_type": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ValidateDiagFunc: verify.ToDiagFunc(validation.StringInSlice([]string{"PUBLC", "GLBL"}, false)),
				Deprecated:       "This parameter is no longer used.",
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
			"filter": DataSourceFiltersSchema(),

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
	config := meta.(*conn.ProviderConfig)

	var resources []map[string]interface{}
	var err error

	if config.SupportVPC {
		resources, err = getVpcPublicIpList(d, meta.(*conn.ProviderConfig))
	} else {
		resources, err = getClassicPublicIpList(d, meta.(*conn.ProviderConfig))
	}

	if err != nil {
		return err
	}

	if f, ok := d.GetOk("filter"); ok {
		resources = ApplyFilters(f.(*schema.Set), resources, DataSourceNcloudPublicIp().Schema)
	}

	if err := verify.ValidateOneResult(len(resources)); err != nil {
		return err
	}

	SetSingularResourceDataFromMapSchema(DataSourceNcloudPublicIp(), d, resources[0])

	return nil
}

func getClassicPublicIpList(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
	client := config.Client
	regionNo := config.RegionNo

	reqParams := &server.GetPublicIpInstanceListRequest{
		RegionNo: &regionNo,
		ZoneNo:   StringPtrOrNil(d.GetOk("zone")),
	}

	if isAssociated, ok := d.GetOk("is_associated"); ok {
		reqParams.IsAssociated = ncloud.Bool(isAssociated.(bool))
	}

	if v, ok := d.GetOk("id"); ok {
		reqParams.PublicIpInstanceNoList = []*string{ncloud.String(v.(string))}
	}

	LogCommonRequest("getClassicPublicIpList", reqParams)
	resp, err := client.Server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		LogErrorResponse("getClassicPublicIpList", err, reqParams)
		return nil, err
	}
	LogCommonResponse("getClassicPublicIpList", GetCommonResponse(resp))

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

		if m := FlattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
			instance["status"] = m["code"]
		}

		if m := FlattenCommonCode(r.PublicIpKindType); m["code"] != nil {
			instance["kind_type"] = m["code"]
		}

		if m := zone.FlattenZone(r.Zone); m["zone_code"] != nil {
			instance["zone"] = m["zone_code"]
		}

		if serverInstance := r.ServerInstanceAssociatedWithPublicIp; serverInstance != nil {
			SetStringIfNotNilAndEmpty(instance, "server_instance_no", serverInstance.ServerInstanceNo)
			SetStringIfNotNilAndEmpty(instance, "server_name", serverInstance.ServerName)
		}

		resources = append(resources, instance)
	}

	return resources, nil
}

func getVpcPublicIpList(d *schema.ResourceData, config *conn.ProviderConfig) ([]map[string]interface{}, error) {
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

	LogCommonRequest("getVpcPublicIpList", reqParams)
	resp, err := client.Vserver.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		LogErrorResponse("getVpcPublicIpList", err, reqParams)
		return nil, err
	}
	LogCommonResponse("getVpcPublicIpList", GetCommonResponse(resp))

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

		if m := FlattenCommonCode(r.PublicIpInstanceStatus); m["code"] != nil {
			instance["status"] = m["code"]
		}

		resources = append(resources, instance)
	}

	return resources, nil
}
