package ncloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNcloudRouteTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudRouteTableAssociationCreate,
		Read:   resourceNcloudRouteTableAssociationRead,
		Update: resourceNcloudRouteTableAssociationUpdate,
		Delete: resourceNcloudRouteTableAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				routeTableNo, subnetNo, err := getNosFromID(d.Id())
				if err != nil {
					return nil, err
				}

				d.Set("route_table_no", routeTableNo)
				d.Set("subnet_no", subnetNo)
				d.SetId(associationID(routeTableNo, subnetNo))
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"subnet_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_no": {
				Type:     schema.TypeString,
				Required: true,
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

	reqParams := &vpc.AddRouteTableSubnetRequest{
		VpcNo:        ncloud.String(*routeTable.VpcNo),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		SubnetNoList: []*string{ncloud.String(d.Get("subnet_no").(string))},
	}

	logCommonRequest("resource_ncloud_route_table_association > AddRouteTableSubnet", reqParams)
	resp, err := client.vpc.V2Api.AddRouteTableSubnet(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_route_table_association > AddRouteTableSubnet", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_route_table_association > AddRouteTableSubnet", resp)

	d.SetId(routeRuleHash(d.Get("route_table_no").(string), d.Get("subnet_no").(string)))

	log.Printf("[INFO] Association ID: %s", d.Id())

	waitForNcloudRouteTableAssociationTableUpdate(client, d.Get("route_table_no").(string))

	return resourceNcloudRouteTableAssociationRead(d, meta)
}

func resourceNcloudRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	routeTable, err := getRouteTableInstance(client, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable == nil {
		return fmt.Errorf("No matching route table: %s", d.Get("route_table_no"))
	}

	instance, err := getRouteTableAssociationInstance(client, d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(associationID(*routeTable.RouteTableNo, *instance.SubnetNo))
	d.Set("route_table_no", routeTable.RouteTableNo)
	d.Set("subnet_no", instance.SubnetNo)

	return nil
}

func resourceNcloudRouteTableAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudRouteTableAssociationRead(d, meta)
}

func resourceNcloudRouteTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	routeTable, err := getRouteTableInstance(client, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable == nil {
		return fmt.Errorf("No matching route table: %s", d.Get("route_table_no"))
	}

	reqParams := &vpc.RemoveRouteTableSubnetRequest{
		VpcNo:        ncloud.String(*routeTable.VpcNo),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		SubnetNoList: []*string{ncloud.String(d.Get("subnet_no").(string))},
	}

	logCommonRequest("resource_ncloud_route_table_association > RemoveRouteTableSubnet", reqParams)
	resp, err := client.vpc.V2Api.RemoveRouteTableSubnet(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_route_table_association > RemoveRouteTableSubnet", err, reqParams)
		return err
	}

	logResponse("resource_ncloud_route_table_association > RemoveRouteTableSubnet", resp)

	waitForNcloudRouteTableAssociationTableUpdate(client, d.Get("route_table_no").(string))

	return nil
}

func waitForNcloudRouteTableAssociationTableUpdate(client *NcloudAPIClient, id string) error {
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

func getRouteTableAssociationInstance(client *NcloudAPIClient, id string) (*vpc.Subnet, error) {
	routeTableNo, subnetNo, err := getNosFromID(id)

	if err != nil {
		return nil, err
	}

	reqParams := &vpc.GetRouteTableSubnetListRequest{
		RouteTableNo: ncloud.String(routeTableNo),
	}

	logCommonRequest("resource_ncloud_route_table_association > GetRouteTableSubnetList", reqParams)
	resp, err := client.vpc.V2Api.GetRouteTableSubnetList(reqParams)
	if err != nil {
		logErrorResponse("resource_ncloud_route_table_association > GetRouteTableSubnetList", err, reqParams)
		return nil, err
	}
	logResponse("resource_ncloud_route_table_association > GetRouteTableSubnetList", resp)

	if resp.SubnetList != nil {
		for _, i := range resp.SubnetList {
			if *i.SubnetNo == subnetNo {
				return i, nil
			}
		}
		return nil, nil
	}

	return nil, nil
}

func associationID(routeTableNo, subnetNo string) string {
	return fmt.Sprintf("%s:%s", routeTableNo, subnetNo)
}

func getNosFromID(id string) (string, string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%q), expected ROUTE_TABLE_NO:SUBNET_NO", id)
	}
	return idParts[0], idParts[1], nil
}
