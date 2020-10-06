package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudNatGateway_basic(t *testing.T) {
	resourceName := "ncloud_nat_gateway.nat_gateway"
	dataName := "data.ncloud_nat_gateway.by_id"
	name := fmt.Sprintf("tf-data-testacc-nat-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudNatGatewayConfig(name),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					testAccCheckDataSourceID("data.ncloud_nat_gateway.by_filter"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),
					resource.TestCheckResourceAttrPair(dataName, "nat_gateway_no", resourceName, "nat_gateway_no"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNatGatewayConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  zone        = "KR-1"
  name        = "%[1]s"
  description = "description"
}

data "ncloud_nat_gateway" "by_id" {
  nat_gateway_no = ncloud_nat_gateway.nat_gateway.nat_gateway_no
}

data "ncloud_nat_gateway" "by_filter" {
  filter {
		name   = "nat_gateway_no"
		values = [ncloud_nat_gateway.nat_gateway.nat_gateway_no]
	}
}
`, name)
}
