package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudRouteTablesBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRouteTablesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_route_tables.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudRouteTablesFilter(t *testing.T) {
	dataName := "data.ncloud_route_tables.filter"
	name := fmt.Sprintf("test-rt-data-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRouteTablesConfigFilter(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "route_tables.#", "1"),
					resource.TestCheckResourceAttr(dataName, "route_tables.0.name", name),
					resource.TestCheckResourceAttr(dataName, "route_tables.0.is_default", "true"),
					resource.TestCheckResourceAttr(dataName, "route_tables.0.supported_subnet_type", "PRIVATE"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudRouteTablesName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRouteTablesConfigName("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_route_tables.by_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudRouteTablesVpcNo(t *testing.T) {
	name := fmt.Sprintf("test-table-data-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRouteTablesConfigVpcNo(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_route_tables.by_vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudRouteTablesConfig() string {
	return fmt.Sprintf(`
data "ncloud_route_tables" "all" {}
`)
}

func testAccDataSourceNcloudRouteTablesConfigFilter(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

data "ncloud_route_tables" "filter" {
	filter {
		name = "name"
		values = ["%[1]s"]
	}

	filter {
		name = "is_default"
		values = ["true"]
	}

	filter {
		name = "supported_subnet_type"
		values = ["PRIVATE"]
	}
}
`, name)
}

func testAccDataSourceNcloudRouteTablesConfigName(name string) string {
	return fmt.Sprintf(`
data "ncloud_route_tables" "by_name" {
	name               = "%s"
}
`, name)
}

func testAccDataSourceNcloudRouteTablesConfigVpcNo(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

data "ncloud_route_tables" "by_vpc_no" {
	vpc_no          = ncloud_vpc.vpc.id
	is_default      = true
}
`, name)
}
