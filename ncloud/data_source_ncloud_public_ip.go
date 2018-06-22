package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudPublicIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPublicIPRead,

		Schema: map[string]*schema.Schema{
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
			},
			"is_associated": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip_instance_no_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"search_filter_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_filter_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sorted_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sorting_order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_no": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"public_ip_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_line_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"public_ip_instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"public_ip_instance_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"public_ip_kind_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"server_instance": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_instance_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NcloudSdk).conn

	reqParams := new(sdk.RequestPublicIPInstanceList)
	reqParams.InternetLineTypeCode = d.Get("internet_line_type_code").(string)
	reqParams.IsAssociated = d.Get("is_associated").(string)
	reqParams.PublicIPInstanceNoList = StringList(d.Get("public_ip_instance_no_list").([]interface{}))
	reqParams.PublicIPList = StringList(d.Get("public_ip_list").([]interface{}))
	reqParams.SearchFilterName = d.Get("search_filter_name").(string)
	reqParams.SearchFilterValue = d.Get("search_filter_value").(string)
	reqParams.RegionNo = parseRegionNoParameter(conn, d)
	reqParams.ZoneNo = d.Get("zone_no").(string)
	reqParams.SortedBy = d.Get("sorted_by").(string)
	reqParams.SortingOrder = d.Get("sorting_order").(string)
	reqParams.PageNo = d.Get("page_no").(int)
	reqParams.PageSize = d.Get("page_size").(int)
	resp, err := conn.GetPublicIPInstanceList(reqParams)

	if err != nil {
		logErrorResponse("Get Public IP Instance", err, reqParams)
		return err
	}
	publicIPInstanceList := resp.PublicIPInstanceList
	var publicIPInstance sdk.PublicIPInstance

	if len(publicIPInstanceList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	if len(publicIPInstanceList) > 1 && d.Get("most_recent").(bool) {
		// Query returned single result.
		publicIPInstance = mostRecentPublicIP(publicIPInstanceList)
	} else {
		publicIPInstance = publicIPInstanceList[0]
	}

	return publicIPAttributes(d, publicIPInstance)
}

func publicIPAttributes(d *schema.ResourceData, instance sdk.PublicIPInstance) error {

	d.SetId(instance.PublicIPInstanceNo)
	d.Set("public_ip_instance_no", instance.PublicIPInstanceNo)
	d.Set("public_ip", instance.PublicIP)
	d.Set("public_ip_description", instance.PublicIPDescription)
	d.Set("create_date", instance.CreateDate)
	d.Set("internet_line_type", setCommonCode(instance.InternetLineType))
	d.Set("public_ip_instance_status_name", instance.PublicIPInstanceStatusName)
	d.Set("public_ip_instance_status", setCommonCode(instance.PublicIPInstanceStatus))
	d.Set("public_ip_instance_operation", setCommonCode(instance.PublicIPInstanceOperation))
	d.Set("public_ip_kind_type", setCommonCode(instance.PublicIPKindType))

	if instance.ServerInstance.ServerInstanceNo != "" {
		serverInstance := instance.ServerInstance
		mapping := map[string]interface{}{
			"server_instance_no": serverInstance.ServerInstanceNo,
			"server_name":        serverInstance.ServerName,
			"create_date":        serverInstance.CreateDate,
		}
		d.Set("server_instance", mapping)
	}
	return nil
}
