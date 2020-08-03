package ncloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNcloudRouteTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudRouteTableAssociationCreate,
		Read:   resourceNcloudRouteTableAssociationRead,
		Update: resourceNcloudRouteTableAssociationUpdate,
		Delete: resourceNcloudRouteTableAssociationDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 32),
			},
			"target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NATGW", "VPCPEERING", "VGW"}, false),
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

func resourceNcloudRouteTableAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	routeTable, err := getRouteTableInstance(client, d.Get("route_table_no").(string))
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
		VpcNo:        ncloud.String(*routeTable.VpcNo),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		RouteList:    []*vpc.RouteParameter{routeParams},
	}

	logCommonRequest("resource_ncloud_route_table_association > AddRoute", reqParams)
	resp, err := client.vpc.V2Api.AddRoute(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_route_table_association > AddRoute", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_route_table_association > AddRoute", resp)

	d.SetId(routeRuleHash(d.Get("route_table_no").(string), d.Get("destination_cidr_block").(string)))

	log.Printf("[INFO] Route Table Rule ID: %s", d.Id())

	waitForNcloudRouteTableUpdate(client, d.Get("route_table_no").(string))

	return resourceNcloudRouteTableAssociationRead(d, meta)
}

func resourceNcloudRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	routeTable, err := getRouteTableInstance(client, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable != nil {
		d.Set("vpc_no", routeTable.VpcNo)
	}

	instance, err := getRouteTableRouteInstance(client, d)
	if err != nil {
		d.SetId("")
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

func resourceNcloudRouteTableAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudRouteTableAssociationRead(d, meta)
}

func resourceNcloudRouteTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	routeParams := &vpc.RouteParameter{
		DestinationCidrBlock: ncloud.String(d.Get("destination_cidr_block").(string)),
		TargetTypeCode:       ncloud.String(d.Get("target_type").(string)),
		TargetName:           ncloud.String(d.Get("target_name").(string)),
		TargetNo:             ncloud.String(d.Get("target_no").(string)),
	}

	reqParams := &vpc.RemoveRouteRequest{
		VpcNo:        ncloud.String(d.Get("vpc_no").(string)),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		RouteList:    []*vpc.RouteParameter{routeParams},
	}

	logCommonRequest("resource_ncloud_route_table_association > RemoveRoute", reqParams)
	resp, err := client.vpc.V2Api.RemoveRoute(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_route_table_association > RemoveRoute", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_route_table_association > RemoveRoute", resp)

	waitForNcloudRouteTableUpdate(client, d.Get("route_table_no").(string))

	return nil
}

func waitForNcloudRouteTableUpdate(client *NcloudAPIClient, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getRouteTableInstance(client, id)
			return VpcCommonStateRefreshFunc(instance, err, "RouteTableStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Route Table (%s) to become running: %s", id, err)
	}

	return nil
}

func getRouteTableRouteInstance(client *NcloudAPIClient, d *schema.ResourceData) (*vpc.Route, error) {
	reqParams := &vpc.GetRouteListRequest{
		VpcNo:        ncloud.String(d.Get("vpc_no").(string)),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
	}

	logCommonRequest("resource_ncloud_route_table_association > GetRouteList", reqParams)
	resp, err := client.vpc.V2Api.GetRouteList(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_route_table_association > GetRouteList", err, reqParams)
		return nil, err
	}
	logResponse("resource_ncloud_route_table_association > GetRouteList", resp)

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
	return fmt.Sprintf("route-%d", hashcode.String(buf.String()))
}
