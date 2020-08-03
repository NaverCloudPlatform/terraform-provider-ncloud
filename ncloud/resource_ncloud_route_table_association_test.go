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

func TestAccresourceNcloudRouteTableAssociation_basic(t *testing.T) {
	var association vpc.Subnet
	var routeTableNo string

	name := fmt.Sprintf("test-assoc-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_route_table_association.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccresourceNcloudRouteTableAssociationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableAssociationExists(resourceName, &association, &routeTableNo),
					testAccCheckDataSourceID(resourceName),
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

func TestAccresourceNcloudRouteTableAssociation_disappears(t *testing.T) {
	var association vpc.Subnet
	var routeTableNo string

	name := fmt.Sprintf("test-route-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_route_table_association.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccresourceNcloudRouteTableAssociationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableAssociationExists(resourceName, &association, &routeTableNo),
					testAccCheckRouteTableAssociationDisappears(&association, &routeTableNo),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccresourceNcloudRouteTableAssociationConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_subnet" "subnet_a" {
	vpc_no             = ncloud_vpc.vpc.id
	name               = "%[1]s"
	subnet             = "10.3.1.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_route_table" "route_table" {
	vpc_no                = ncloud_vpc.vpc.id
	name                  = "%[1]s"
	description           = "for test"
	supported_subnet_type = "PRIVATE"
}

resource "ncloud_route_table_association" "test" {
	route_table_no        = ncloud_route_table.route_table.id
	subnet_no             = ncloud_subnet.subnet_a.id
}
`, name)
}

func testAccCheckRouteTableAssociationExists(n string, subnet *vpc.Subnet, routeTableNo *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route table association id is set: %s", n)
		}

		*routeTableNo = rs.Primary.Attributes["route_table_no"]
		config := testAccProvider.Meta().(*ProviderConfig)

		reqParams := &vpc.GetRouteTableSubnetListRequest{
			RegionCode:   &config.RegionCode,
			RouteTableNo: ncloud.String(rs.Primary.Attributes["route_table_no"]),
		}

		logCommonRequest("resource_ncloud_route_test > GetRouteTableSubnetList", reqParams)
		resp, err := config.Client.vpc.V2Api.GetRouteTableSubnetList(reqParams)
		if err != nil {
			logErrorResponse("resource_ncloud_route_test > GetRouteTableSubnetList", err, reqParams)
			return fmt.Errorf("Not found: %s", n)
		}
		logResponse("resource_ncloud_route_test > GetRouteTableSubnetList", resp)

		if resp.SubnetList != nil {
			for _, i := range resp.SubnetList {
				if *i.SubnetNo == rs.Primary.Attributes["subnet_no"] {
					*subnet = *i
				}
			}
			return nil
		}

		return fmt.Errorf("Entry not found: %v", resp.SubnetList)
	}
}

func testAccNcloudRouteTableAssociationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		routeTableNo := rs.Primary.Attributes["route_table_no"]
		subnetNo := rs.Primary.Attributes["subnet_no"]

		return fmt.Sprintf("%s:%s", routeTableNo, subnetNo), nil
	}
}

func testAccCheckRouteTableAssociationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_route_table_association" {
			continue
		}

		routeTable, err := getRouteTableInstance(config, rs.Primary.Attributes["route_table_no"])

		if err != nil {
			return err
		}

		if routeTable == nil {
			return nil
		}

		instance, err := getRouteTableAssociationInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Route table association still exists")
		}
	}

	return nil
}

func testAccCheckRouteTableAssociationDisappears(instance *vpc.Subnet, routeTableNo *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		routeTable, err := getRouteTableInstance(config, *routeTableNo)
		if err != nil {
			return err
		}

		if routeTable == nil {
			return fmt.Errorf("No matching route table: %s", *routeTableNo)
		}

		reqParams := &vpc.RemoveRouteTableSubnetRequest{
			RegionCode:   &config.RegionCode,
			VpcNo:        routeTable.VpcNo,
			RouteTableNo: routeTable.RouteTableNo,
			SubnetNoList: []*string{instance.SubnetNo},
		}

		_, err = config.Client.vpc.V2Api.RemoveRouteTableSubnet(reqParams)

		waitForNcloudRouteTableAssociationTableUpdate(config, *routeTableNo)

		return err
	}
}
