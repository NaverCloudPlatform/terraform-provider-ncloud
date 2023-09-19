package vpc_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	vpcservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
)

func TestAccResourceNcloudRouteTable_basic(t *testing.T) {
	var routeTable vpc.RouteTable
	resourceName := "ncloud_route_table.foo"
	name := fmt.Sprintf("test-table-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteTableConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "route_table_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudRouteTable_disappears(t *testing.T) {
	var routeTable vpc.RouteTable
	resourceName := "ncloud_route_table.foo"
	name := fmt.Sprintf("test-table-disappear-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteTableConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
					testAccCheckRouteTableDisappears(&routeTable),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudRouteTable_onlyRequiredParam(t *testing.T) {
	var routeTable vpc.RouteTable
	resourceName := "ncloud_route_table.foo"
	name := fmt.Sprintf("test-table-required-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteTableConfigOnlyRequired(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "route_table_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^[a-z0-9]+$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudRouteTable_updateName(t *testing.T) {
	var routeTable vpc.RouteTable
	resourceName := "ncloud_route_table.foo"
	name := fmt.Sprintf("test-table-update-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteTableConfigOnlyRequired(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudRouteTable_description(t *testing.T) {
	var routeTable vpc.RouteTable
	resourceName := "ncloud_route_table.foo"
	name := fmt.Sprintf("test-table-desc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudRouteTableConfigDescription(name, "foo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
					resource.TestCheckResourceAttr(resourceName, "description", "foo"),
				),
			},
			{
				Config: testAccResourceNcloudRouteTableConfigDescription(name, "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(resourceName, &routeTable),
					resource.TestCheckResourceAttr(resourceName, "description", "bar"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceNcloudRouteTableConfig(name string) string {
	return testAccResourceNcloudRouteTableConfigDescription(name, "for acc test")
}

func testAccResourceNcloudRouteTableConfigDescription(name, description string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_route_table" "foo" {
	vpc_no                = ncloud_vpc.vpc.vpc_no
	name                  = "%[1]s"
	description           = "%[2]s"
	supported_subnet_type = "PUBLIC"
}
`, name, description)
}

func testAccResourceNcloudRouteTableConfigOnlyRequired(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_route_table" "foo" {
	vpc_no                = ncloud_vpc.vpc.vpc_no
	supported_subnet_type = "PUBLIC"
}
`, name)
}

func testAccCheckRouteTableExists(n string, routeTable *vpc.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Table id is set: %s", n)
		}

		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		instance, err := vpcservice.GetRouteTableInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*routeTable = *instance

		return nil
	}
}

func testAccCheckRouteTableDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_route_table" {
			continue
		}

		instance, err := vpcservice.GetRouteTableInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Route Table still exists")
		}
	}

	return nil
}

func testAccCheckRouteTableDisappears(instance *vpc.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

		reqParams := &vpc.DeleteRouteTableRequest{
			RegionCode:   &config.RegionCode,
			RouteTableNo: instance.RouteTableNo,
		}

		_, err := config.Client.Vpc.V2Api.DeleteRouteTable(reqParams)

		if err := vpcservice.WaitForNcloudRouteTableDeletion(config, *instance.RouteTableNo); err != nil {
			return err
		}

		return err
	}
}
