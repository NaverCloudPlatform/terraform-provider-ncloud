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

func TestAccresourceNcloudRoute_basic(t *testing.T) {
	var route vpc.Route
	name := fmt.Sprintf("test-route-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_route.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteExists(resourceName, &route),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccNcloudRouteImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccresourceNcloudRoute_disappears(t *testing.T) {
	var route vpc.Route
	name := fmt.Sprintf("test-route-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_route.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteExists(resourceName, &route),
					testAccCheckRouteDisappears(&route),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNcloudRouteConfig(name string) string {
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

resource "ncloud_route" "foo" {
	route_table_no         = ncloud_route_table.route_table.id
	destination_cidr_block = "0.0.0.0/0"
	target_type            = "NATGW"
	target_name            = ncloud_nat_gateway.nat_gateway.name
	target_no              = ncloud_nat_gateway.nat_gateway.id
}
`, name)
}

func testAccCheckRouteExists(n string, route *vpc.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL Rule id is set: %s", n)
		}

		config := testAccProvider.Meta().(*ProviderConfig)

		reqParams := &vpc.GetRouteListRequest{
			VpcNo:        ncloud.String(rs.Primary.Attributes["vpc_no"]),
			RouteTableNo: ncloud.String(rs.Primary.Attributes["route_table_no"]),
		}

		logCommonRequest("GetRouteList", reqParams)
		resp, err := config.Client.vpc.V2Api.GetRouteList(reqParams)
		if err != nil {
			logErrorResponse("GetRouteList", err, reqParams)
			return fmt.Errorf("Not found: %s", n)
		}
		logResponse("GetRouteList", resp)

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

func testAccNcloudRouteImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
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

func testAccCheckRouteDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_route" {
			continue
		}

		instance, err := getRouteTableInstance(config, rs.Primary.Attributes["route_table_no"])

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

		resp, err := config.Client.vpc.V2Api.GetRouteList(reqParams)
		if err != nil {
			logErrorResponse("GetRouteList", err, reqParams)
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

func testAccCheckRouteDisappears(instance *vpc.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		routeTable, err := getRouteTableInstance(config, *instance.RouteTableNo)
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

		_, err = config.Client.vpc.V2Api.RemoveRoute(reqParams)

		if err := waitForNcloudRouteTableUpdate(config, *instance.RouteTableNo); err != nil {
			return err
		}

		return err
	}
}
