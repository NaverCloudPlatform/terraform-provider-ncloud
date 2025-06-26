package vpc

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func ResourceNcloudRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudRouteCreate,
		Read:   resourceNcloudRouteRead,
		Update: resourceNcloudRouteUpdate,
		Delete: resourceNcloudRouteDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ":")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected ROUTE_TABLE_NO:DESTINATION_CIDR_BLOCK", d.Id())
				}
				routeTableNo := idParts[0]
				destinationCidrBlock := idParts[1]

				d.Set("route_table_no", routeTableNo)
				d.Set("destination_cidr_block", destinationCidrBlock)
				d.SetId(routeRuleHash(routeTableNo, destinationCidrBlock))
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"route_table_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_cidr_block": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDRNetwork(0, 32)),
			},
			"target_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NATGW", "VPCPEERING", "VGW"}, false)),
			},
			"target_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNcloudRouteCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	routeTable, err := GetRouteTableInstance(config, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable == nil {
		return fmt.Errorf("No matching route table: %s", d.Get("route_table_no"))
	}

	routeParams := &vpc.RouteParameter{
		DestinationCidrBlock: ncloud.String(d.Get("destination_cidr_block").(string)),
		TargetTypeCode:       ncloud.String(d.Get("target_type").(string)),
		TargetName:           ncloud.String(d.Get("target_name").(string)),
		TargetNo:             ncloud.String(d.Get("target_no").(string)),
	}

	reqParams := &vpc.AddRouteRequest{
		RegionCode:   &config.RegionCode,
		VpcNo:        ncloud.String(*routeTable.VpcNo),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		RouteList:    []*vpc.RouteParameter{routeParams},
	}

	var resp *vpc.AddRouteResponse
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		LogCommonRequest("AddRoute", reqParams)
		resp, err = config.Client.Vpc.V2Api.AddRoute(reqParams)

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == "1017013" {
				LogErrorResponse("retry add Route", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		LogErrorResponse("AddRoute", err, reqParams)
		return err
	}

	LogResponse("AddRoute", resp)

	d.SetId(routeRuleHash(d.Get("route_table_no").(string), d.Get("destination_cidr_block").(string)))

	log.Printf("[INFO] Route ID: %s", d.Id())

	if err := WaitForNcloudRouteTableUpdate(config, d.Get("route_table_no").(string)); err != nil {
		return err
	}

	return resourceNcloudRouteRead(d, meta)
}

func resourceNcloudRouteRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	routeTable, err := GetRouteTableInstance(config, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable != nil {
		d.Set("vpc_no", routeTable.VpcNo)
	} else {
		d.SetId("")
		return nil
	}

	instance, err := getRouteInstance(config, d)
	if err != nil {
		errBody, _ := GetCommonErrorBody(err)
		if errBody.ReturnCode == "1017007" { // Route Table was not found
			d.SetId("")
		}
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(routeRuleHash(*instance.RouteTableNo, *instance.DestinationCidrBlock))
	d.Set("route_table_no", instance.RouteTableNo)
	d.Set("destination_cidr_block", instance.DestinationCidrBlock)
	d.Set("target_type", instance.TargetType.Code)
	d.Set("target_name", instance.TargetName)
	d.Set("target_no", instance.TargetNo)
	d.Set("is_default", instance.IsDefault)

	return nil
}

func resourceNcloudRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudRouteRead(d, meta)
}

func resourceNcloudRouteDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*conn.ProviderConfig)

	routeParams := &vpc.RouteParameter{
		DestinationCidrBlock: ncloud.String(d.Get("destination_cidr_block").(string)),
		TargetTypeCode:       ncloud.String(d.Get("target_type").(string)),
		TargetName:           ncloud.String(d.Get("target_name").(string)),
		TargetNo:             ncloud.String(d.Get("target_no").(string)),
	}

	reqParams := &vpc.RemoveRouteRequest{
		RegionCode:   &config.RegionCode,
		VpcNo:        ncloud.String(d.Get("vpc_no").(string)),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		RouteList:    []*vpc.RouteParameter{routeParams},
	}

	var resp *vpc.RemoveRouteResponse
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error

		LogCommonRequest("RemoveRoute", reqParams)
		resp, err = config.Client.Vpc.V2Api.RemoveRoute(reqParams)

		if err != nil {
			errBody, _ := GetCommonErrorBody(err)
			if errBody.ReturnCode == "1017013" {
				LogErrorResponse("retry remove Route", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		LogErrorResponse("RemoveRoute", err, reqParams)
		return err
	}

	LogResponse("RemoveRoute", resp)

	if err := WaitForNcloudRouteTableUpdate(config, d.Get("route_table_no").(string)); err != nil {
		return err
	}

	return nil
}

func WaitForNcloudRouteTableUpdate(config *conn.ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := GetRouteTableInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "RouteTableStatus")
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Route Table (%s) to become running: %s", id, err)
	}

	return nil
}

func getRouteInstance(config *conn.ProviderConfig, d *schema.ResourceData) (*vpc.Route, error) {
	reqParams := &vpc.GetRouteListRequest{
		RegionCode:   &config.RegionCode,
		VpcNo:        ncloud.String(d.Get("vpc_no").(string)),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
	}

	LogCommonRequest("GetRouteList", reqParams)
	resp, err := config.Client.Vpc.V2Api.GetRouteList(reqParams)
	if err != nil {
		LogErrorResponse("GetRouteList", err, reqParams)
		return nil, err
	}
	LogResponse("GetRouteList", resp)

	if resp.RouteList != nil {
		for _, i := range resp.RouteList {
			if *i.DestinationCidrBlock == d.Get("destination_cidr_block").(string) {
				return i, nil
			}
		}
		return nil, nil
	}

	return nil, nil
}

func routeRuleHash(routeTableNo, destinationCidrBlock string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routeTableNo))
	buf.WriteString(fmt.Sprintf("%s-", destinationCidrBlock))
	return fmt.Sprintf("route-%d", Hashcode(buf.String()))
}
