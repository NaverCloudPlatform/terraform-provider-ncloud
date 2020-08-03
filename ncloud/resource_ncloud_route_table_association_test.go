package ncloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudRouteTableAssociation_basic(t *testing.T) {
	var route vpc.Route
	name := fmt.Sprintf("test-table-ass-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_route_table_association.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableAssocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteTableAssociationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableAssociationExists(resourceName, &route),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccNcloudRouteTableAssociationImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudRouteTableAssociation_disappears(t *testing.T) {
	var route vpc.Route
	name := fmt.Sprintf("test-table-ass-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_route_table_association.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableAssocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteTableAssociationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableAssociationExists(resourceName, &route),
					testAccCheckRouteTableAssocationDisappears(&route),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNcloudRouteTableAssociationConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_route_table" "route_table" {
	vpc_no                = ncloud_vpc.vpc.id
	name                  = "%[1]s"
	description           = "for test"
	supported_subnet_type = "PRIVATE"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.id
  zone        = "KR-1"
}

resource "ncloud_route_table_association" "foo" {
	route_table_no         = ncloud_route_table.route_table.id
	destination_cidr_block = "0.0.0.0/0"
	target_type            = "NATGW"
	target_name            = ncloud_nat_gateway.nat_gateway.name
	target_no              = ncloud_nat_gateway.nat_gateway.id
}
`, name)
}

func testAccCheckRouteTableAssociationExists(n string, route *vpc.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL Rule id is set: %s", n)
		}

		client := testAccProvider.Meta().(*NcloudAPIClient)

		reqParams := &vpc.GetRouteListRequest{
			VpcNo:        ncloud.String(rs.Primary.Attributes["vpc_no"]),
			RouteTableNo: ncloud.String(rs.Primary.Attributes["route_table_no"]),
		}

		logCommonRequest("resource_ncloud_route_table_association_test > GetRouteList", reqParams)
		resp, err := client.vpc.V2Api.GetRouteList(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_route_table_association_test > GetRouteList", err, reqParams)
			return fmt.Errorf("Not found: %s", n)
		}
		logResponse("resource_ncloud_route_table_association_test > GetRouteList", resp)

		if resp.RouteList != nil {
			for _, i := range resp.RouteList {
				if *i.DestinationCidrBlock == rs.Primary.Attributes["destination_cidr_block"] {
					*route = *i
				}
			}
			return nil
		}

		return fmt.Errorf("Entry not found: %v", resp.RouteList)
	}
}

func testAccNcloudRouteTableAssociationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		routeTableNo := rs.Primary.Attributes["route_table_no"]
		destinationCidrBlock := rs.Primary.Attributes["destination_cidr_block"]

		return fmt.Sprintf("%s:%s", routeTableNo, destinationCidrBlock), nil
	}
}

func testAccCheckRouteTableAssocationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*NcloudAPIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_route_table_association" {
			continue
		}

		instance, err := getRouteTableInstance(client, rs.Primary.Attributes["route_table_no"])

		if err != nil {
			return err
		}

		if instance == nil {
			return nil
		}

		reqParams := &vpc.GetRouteListRequest{
			VpcNo:        ncloud.String(rs.Primary.Attributes["vpc_no"]),
			RouteTableNo: ncloud.String(rs.Primary.Attributes["route_table_no"]),
		}

		resp, err := client.vpc.V2Api.GetRouteList(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_route_table_association_test > GetRouteList", err, reqParams)
			return err
		}

		if resp.RouteList != nil {
			for _, i := range resp.RouteList {
				if *i.DestinationCidrBlock == rs.Primary.Attributes["destination_cidr_block"] {
					return errors.New("Route Table Rule still exists")
				}
			}
			return nil
		}
	}

	return nil
}

func testAccCheckRouteTableAssocationDisappears(instance *vpc.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*NcloudAPIClient)
		routeTable, err := getRouteTableInstance(client, *instance.RouteTableNo)
		if err != nil {
			return err
		}

		if routeTable == nil {
			return fmt.Errorf("No matching route table: %s", *instance.RouteTableNo)
		}

		routeParams := &vpc.RouteParameter{
			DestinationCidrBlock: instance.DestinationCidrBlock,
			TargetTypeCode:       instance.TargetType.Code,
			TargetName:           instance.TargetName,
			TargetNo:             instance.TargetNo,
		}

		reqParams := &vpc.RemoveRouteRequest{
			VpcNo:        routeTable.VpcNo,
			RouteTableNo: instance.RouteTableNo,
			RouteList:    []*vpc.RouteParameter{routeParams},
		}

		_, err = client.vpc.V2Api.RemoveRoute(reqParams)

		waitForNcloudRouteTableUpdate(client, *instance.RouteTableNo)

		return err
	}
}
