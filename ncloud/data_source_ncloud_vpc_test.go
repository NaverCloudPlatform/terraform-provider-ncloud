package ncloud

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudVpc(t *testing.T) {
	rInt := rand.Intn(16)
	cidr := fmt.Sprintf("10.%d.0.0/16", rInt)
	name := fmt.Sprintf("testacc-vpc-basic-%d", rInt)
	resourceName := "ncloud_vpc.test"
	dataName := "data.ncloud_vpc.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttr(dataName, "ipv4_cidr_block", cidr),
					resource.TestCheckResourceAttr(dataName, "name", name),
					testAccCheckDataSourceID("data.ncloud_vpc.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudVpcConfig(name, cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
  name               = "%s"
  ipv4_cidr_block    = "%s"
}

data "ncloud_vpc" "by_id" {
  id = ncloud_vpc.test.id
}

data "ncloud_vpc" "by_filter" {
  filter {
    name = "vpc_no"
    values = [ncloud_vpc.test.id]
  }
}
`, name, cidr)
}
