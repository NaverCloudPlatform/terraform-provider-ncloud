package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudPublicIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPublicIpRead,

		Schema: map[string]*schema.Schema{
			"most_recent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "If more than one result is returned, get the most recent created Public IP.",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
				Description:  "Internet line type code. `PUBLC` (Public), `GLBL` (Global)",
			},
			"is_associated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the public IP address is associated or not.",
			},
			"public_ip_instance_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of public IP instance numbers to get.",
			},
			"public_ip_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of public IP addresses to get.",
			},
			"search_filter_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "`publicIp` (Public IP) | `associatedServerName` (Associated server name)",
			},
			"search_filter_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter value to search",
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
				Description:   "Zone code. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.",
				ConflictsWith: []string{"zone_no"},
			},
			"zone_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone number. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.",
				ConflictsWith: []string{"zone_code"},
			},
			"sorted_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The column based on which you want to sort the list.",
			},
			"sorting_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting order of the list. `ascending` (Ascending) | `descending` (Descending) [case insensitive]. Default: `ascending` Ascending",
			},

			"public_ip_instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance number",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP",
			},
			"public_ip_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP description",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of the public ip",
			},
			"internet_line_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Internet line type",
			},
			"public_ip_instance_status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance status name",
			},
			"public_ip_instance_status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Public IP instance status",
			},
			"public_ip_instance_operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Public IP instance operation",
			},
			"public_ip_kind_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
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
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation date of the server instance",
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudPublicIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := new(server.GetPublicIpInstanceListRequest)
	reqParams.InternetLineTypeCode = ncloud.String(d.Get("internet_line_type_code").(string))
	reqParams.IsAssociated = ncloud.Bool(d.Get("is_associated").(bool))
	reqParams.PublicIpInstanceNoList = expandStringInterfaceList(d.Get("public_ip_instance_no_list").([]interface{}))
	reqParams.PublicIpList = expandStringInterfaceList(d.Get("public_ip_list").([]interface{}))
	reqParams.SearchFilterName = ncloud.String(d.Get("search_filter_name").(string))
	reqParams.SearchFilterValue = ncloud.String(d.Get("search_filter_value").(string))
	reqParams.RegionNo = regionNo
	reqParams.ZoneNo = zoneNo
	reqParams.SortedBy = ncloud.String(d.Get("sorted_by").(string))
	reqParams.SortingOrder = ncloud.String(d.Get("sorting_order").(string))
	// log.Printf("[DEBUG] GetPublicIpInstanceList reqParams: %#v", reqParams)
	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("Get Public IP Instance", err, reqParams)
		return err
	}
	publicIpInstanceList := resp.PublicIpInstanceList
	var publicIpInstance *server.PublicIpInstance

	if len(publicIpInstanceList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	if len(publicIpInstanceList) > 1 && d.Get("most_recent").(bool) {
		// Query returned single result.
		publicIpInstance = mostRecentPublicIp(publicIpInstanceList)
	} else {
		publicIpInstance = publicIpInstanceList[0]
	}

	return publicIPAttributes(d, publicIpInstance)
}

func publicIPAttributes(d *schema.ResourceData, instance *server.PublicIpInstance) error {

	d.SetId(ncloud.StringValue(instance.PublicIpInstanceNo))
	d.Set("public_ip_instance_no", instance.PublicIpInstanceNo)
	d.Set("public_ip", instance.PublicIp)
	d.Set("public_ip_description", instance.PublicIpDescription)
	d.Set("create_date", instance.CreateDate)
	d.Set("public_ip_instance_status_name", instance.PublicIpInstanceStatusName)

	if err := d.Set("internet_line_type", flattenCommonCode(instance.InternetLineType)); err != nil {
		return err
	}
	if err := d.Set("public_ip_instance_status", flattenCommonCode(instance.PublicIpInstanceStatus)); err != nil {
		return err
	}
	if err := d.Set("public_ip_instance_operation", flattenCommonCode(instance.PublicIpInstanceOperation)); err != nil {
		return err
	}
	if err := d.Set("public_ip_kind_type", flattenCommonCode(instance.PublicIpKindType)); err != nil {
		return err
	}

	if serverInstance := instance.ServerInstanceAssociatedWithPublicIp; serverInstance != nil {
		mapping := map[string]interface{}{
			"server_instance_no": ncloud.StringValue(serverInstance.ServerInstanceNo),
			"server_name":        ncloud.StringValue(serverInstance.ServerName),
			"create_date":        ncloud.StringValue(serverInstance.CreateDate),
		}
		d.Set("server_instance", mapping)
	}

	return nil
}
