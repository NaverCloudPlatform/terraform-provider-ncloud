package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudNatGateway_basic(t *testing.T) {
	resourceName := "ncloud_nat_gateway.nat_gateway"
	dataName := "data.ncloud_nat_gateway.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNatGatewayConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),
					resource.TestCheckResourceAttrPair(dataName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNatGatewayConfig() string {
	return `
resource "ncloud_vpc" "vpc" {
	name            = "tf-data-testacc-nat-gateway"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  zone        = "KR-1"
  name        = "tf-data-testacc-nat-gateway"
  description = "description"
}

data "ncloud_nat_gateway" "by_id" {
  instance_no = ncloud_nat_gateway.nat_gateway.instance_no
}
`
}
