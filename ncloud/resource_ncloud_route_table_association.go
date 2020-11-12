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

func init() {
	RegisterResource("ncloud_route_table_association", resourceNcloudRouteTableAssociation())
}

func resourceNcloudRouteTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudRouteTableAssociationCreate,
		Read:   resourceNcloudRouteTableAssociationRead,
		Update: resourceNcloudRouteTableAssociationUpdate,
		Delete: resourceNcloudRouteTableAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				routeTableNo, subnetNo, err := convInstanceID(d.Id())
				if err != nil {
					return nil, err
				}

				d.Set("route_table_no", routeTableNo)
				d.Set("subnet_no", subnetNo)
				d.SetId(convAssociationID(routeTableNo, subnetNo))
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
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `ncloud_route_table_association`")
	}

	routeTable, err := getRouteTableInstance(config, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable == nil {
		return fmt.Errorf("No matching route table: %s", d.Get("route_table_no"))
	}

	reqParams := &vpc.AddRouteTableSubnetRequest{
		RegionCode:   &config.RegionCode,
		VpcNo:        ncloud.String(*routeTable.VpcNo),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		SubnetNoList: []*string{ncloud.String(d.Get("subnet_no").(string))},
	}

	logCommonRequest("AddRouteTableSubnet", reqParams)
	resp, err := config.Client.vpc.V2Api.AddRouteTableSubnet(reqParams)
	if err != nil {
		logErrorResponse("AddRouteTableSubnet", err, reqParams)
		return err
	}

	logResponse("AddRouteTableSubnet", resp)

	d.SetId(convAssociationID(d.Get("route_table_no").(string), d.Get("subnet_no").(string)))

	log.Printf("[INFO] Association ID: %s", d.Id())

	if err := waitForNcloudRouteTableAssociationTableUpdate(config, d.Get("route_table_no").(string)); err != nil {
		return err
	}

	return resourceNcloudRouteTableAssociationRead(d, meta)
}

func resourceNcloudRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	routeTable, err := getRouteTableInstance(config, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable == nil {
		return fmt.Errorf("No matching route table: %s", d.Get("route_table_no"))
	}

	instance, err := getRouteTableAssociationInstance(config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(convAssociationID(*routeTable.RouteTableNo, *instance.SubnetNo))
	d.Set("route_table_no", routeTable.RouteTableNo)
	d.Set("subnet_no", instance.SubnetNo)

	return nil
}

func resourceNcloudRouteTableAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceNcloudRouteTableAssociationRead(d, meta)
}

func resourceNcloudRouteTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	routeTable, err := getRouteTableInstance(config, d.Get("route_table_no").(string))
	if err != nil {
		return err
	}

	if routeTable == nil {
		return fmt.Errorf("No matching route table: %s", d.Get("route_table_no"))
	}

	reqParams := &vpc.RemoveRouteTableSubnetRequest{
		RegionCode:   &config.RegionCode,
		VpcNo:        ncloud.String(*routeTable.VpcNo),
		RouteTableNo: ncloud.String(d.Get("route_table_no").(string)),
		SubnetNoList: []*string{ncloud.String(d.Get("subnet_no").(string))},
	}

	logCommonRequest("RemoveRouteTableSubnet", reqParams)
	resp, err := config.Client.vpc.V2Api.RemoveRouteTableSubnet(reqParams)
	if err != nil {
		logErrorResponse("RemoveRouteTableSubnet", err, reqParams)
		return err
	}

	logResponse("RemoveRouteTableSubnet", resp)

	if err := waitForNcloudRouteTableAssociationTableUpdate(config, d.Get("route_table_no").(string)); err != nil {
		return err
	}

	return nil
}

func waitForNcloudRouteTableAssociationTableUpdate(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"SET"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getRouteTableInstance(config, id)
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

func getRouteTableAssociationInstance(config *ProviderConfig, id string) (*vpc.Subnet, error) {
	routeTableNo, subnetNo, err := convInstanceID(id)

	if err != nil {
		return nil, err
	}

	reqParams := &vpc.GetRouteTableSubnetListRequest{
		RegionCode:   &config.RegionCode,
		RouteTableNo: ncloud.String(routeTableNo),
	}

	logCommonRequest("GetRouteTableSubnetList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetRouteTableSubnetList(reqParams)
	if err != nil {
		logErrorResponse("GetRouteTableSubnetList", err, reqParams)
		return nil, err
	}
	logResponse("GetRouteTableSubnetList", resp)

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

func convAssociationID(routeTableNo, subnetNo string) string {
	return fmt.Sprintf("%s:%s", routeTableNo, subnetNo)
}

func convInstanceID(id string) (string, string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%q), expected ROUTE_TABLE_NO:SUBNET_NO", id)
	}
	return idParts[0], idParts[1], nil
}
