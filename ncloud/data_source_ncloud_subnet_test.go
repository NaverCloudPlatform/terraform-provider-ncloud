package ncloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudSubnet(t *testing.T) {
	cidr := "10.2.2.0/24"
	name := "testacc-data-subnet-basic"
	resourceName := "ncloud_subnet.bar"
	dataName := "data.ncloud_subnet.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccDataSourceNcloudSubnetConfig(name, cidr),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no", resourceName, "subnet_no"),
					resource.TestCheckResourceAttr(dataName, "subnet", cidr),
					resource.TestCheckResourceAttr(dataName, "name", name),
					resource.TestCheckResourceAttr(dataName, "zone", "KR-1"),
					resource.TestMatchResourceAttr(dataName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(dataName, "subnet_type", "PUBLIC"),
					resource.TestCheckResourceAttr(dataName, "usage_type", "GEN"),
					testAccCheckDataSourceID("data.ncloud_subnet.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSubnetConfig(name, cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "foo" {
	name               = "testacc-data-subnet-basic"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "bar" {
	vpc_no             = ncloud_vpc.foo.vpc_no
	name               = "%s"
	subnet             = "%s"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.foo.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

data "ncloud_subnet" "by_id" {
  subnet_no = "${ncloud_subnet.bar.id}"
}

data "ncloud_subnet" "by_filter" {
	filter {
		name   = "subnet_no"
		values = [ncloud_subnet.bar.id]
	}
}
`, name, cidr)
}
