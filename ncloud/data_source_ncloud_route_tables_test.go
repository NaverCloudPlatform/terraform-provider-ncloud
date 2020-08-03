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
