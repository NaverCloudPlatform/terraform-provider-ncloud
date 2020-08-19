package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceNcloudVpcPeering_basic(t *testing.T) {
	name := fmt.Sprintf("test-perring-data-%s", acctest.RandString(5))
	resourceName := "ncloud_vpc_peering.foo"
	dataNameByName := "data.ncloud_vpc_peering.by_name"
	dataNameBySourceName := "data.ncloud_vpc_peering.by_source_name"
	dataNameByTargetName := "data.ncloud_vpc_peering.by_target_name"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcPeeringConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testVpcPeeringCheckResourceAttrPair(dataNameByName, resourceName),
					testVpcPeeringCheckResourceAttrPair(dataNameBySourceName, resourceName),
					testVpcPeeringCheckResourceAttrPair(dataNameByTargetName, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testVpcPeeringCheckResourceAttrPair(dataName string, resourceName string) func(*terraform.State) error {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(dataName, "vpc_peering_no", resourceName, "vpc_peering_no"),
		resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
		resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
		resource.TestCheckResourceAttrPair(dataName, "source_vpc_no", resourceName, "source_vpc_no"),
		resource.TestCheckResourceAttrPair(dataName, "target_vpc_no", resourceName, "target_vpc_no"),
		resource.TestCheckResourceAttrPair(dataName, "target_vpc_name", resourceName, "target_vpc_name"),
		resource.TestCheckResourceAttrPair(dataName, "target_vpc_login_id", resourceName, "target_vpc_login_id"),
		resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
		resource.TestCheckResourceAttrPair(dataName, "is_between_accounts", resourceName, "is_between_accounts"),
		resource.TestCheckResourceAttrPair(dataName, "has_reverse_vpc_peering", resourceName, "has_reverse_vpc_peering"),
	)
}

func testAccDataSourceNcloudVpcPeeringConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "main" {
	name               = "%[1]s-main"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_vpc" "peer" {
	name               = "%[1]s-peer"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_vpc_peering" "foo" {
	name           = "%[1]s-foo"
	source_vpc_no  = ncloud_vpc.main.id
	target_vpc_no  = ncloud_vpc.peer.id
}

data "ncloud_vpc_peering" "by_name" {
	name           = ncloud_vpc_peering.foo.name

	depends_on  = [ncloud_vpc_peering.foo]
}

data "ncloud_vpc_peering" "by_source_name" {
	source_vpc_name = ncloud_vpc.main.name

	depends_on  = [ncloud_vpc_peering.foo]
}

data "ncloud_vpc_peering" "by_target_name" {
	target_vpc_name = ncloud_vpc.peer.name

	depends_on  = [ncloud_vpc_peering.foo]
}
`, name)
}
