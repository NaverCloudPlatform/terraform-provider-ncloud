package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudNatGateway_basic(t *testing.T) {
	resourceName := "ncloud_nat_gateway.nat_gateway"
	dataName := "data.ncloud_nat_gateway.by_id"
	name := fmt.Sprintf("tf-data-testacc-nat-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNatGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					TestAccCheckDataSourceID("data.ncloud_nat_gateway.by_filter"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "nat_gateway_no", resourceName, "nat_gateway_no"),
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_name", resourceName, "subnet_name"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip_no", resourceName, "public_ip_no"),
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

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = cidrsubnet(ncloud_vpc.vpc.ipv4_cidr_block, 8, 1)
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
  usage_type     = "NATGW"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.vpc_no
  subnet_no   = ncloud_subnet.subnet.id
  zone        = "KR-1"
  name        = "%[1]s"
  description = "description"
}

data "ncloud_nat_gateway" "by_id" {
  id = ncloud_nat_gateway.nat_gateway.id
}

data "ncloud_nat_gateway" "by_filter" {
  filter {
		name   = "nat_gateway_no"
		values = [ncloud_nat_gateway.nat_gateway.id]
	}
}
`, name)
}
