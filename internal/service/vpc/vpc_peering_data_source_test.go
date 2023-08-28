package vpc_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudVpcPeering_basic(t *testing.T) {
	name := fmt.Sprintf("test-peering-data-%s", sdkacctest.RandString(5))
	resourceName := "ncloud_vpc_peering.foo"
	dataNameById := "data.ncloud_vpc_peering.by_id"
	dataNameByName := "data.ncloud_vpc_peering.by_name"
	dataNameBySourceName := "data.ncloud_vpc_peering.by_source_name"
	dataNameByTargetName := "data.ncloud_vpc_peering.by_target_name"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcPeeringConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testVpcPeeringCheckResourceAttrPair(dataNameById, resourceName),
					testVpcPeeringCheckResourceAttrPair(dataNameByName, resourceName),
					testVpcPeeringCheckResourceAttrPair(dataNameBySourceName, resourceName),
					testVpcPeeringCheckResourceAttrPair(dataNameByTargetName, resourceName),
				),
			},
		},
	})
}

func testVpcPeeringCheckResourceAttrPair(dataName string, resourceName string) func(*terraform.State) error {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
		resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
		resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
		resource.TestCheckResourceAttrPair(dataName, "source_vpc_no", resourceName, "source_vpc_no"),
		resource.TestCheckResourceAttrPair(dataName, "target_vpc_no", resourceName, "target_vpc_no"),
		resource.TestCheckResourceAttrPair(dataName, "target_vpc_name", resourceName, "target_vpc_name"),
		resource.TestCheckResourceAttrPair(dataName, "target_vpc_login_id", resourceName, "target_vpc_login_id"),
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

data "ncloud_vpc_peering" "by_id" {
	id             = ncloud_vpc_peering.foo.id
	depends_on     = [ncloud_vpc_peering.foo]
}

data "ncloud_vpc_peering" "by_name" {
	name        = ncloud_vpc_peering.foo.name
	depends_on  = [ncloud_vpc_peering.foo]
}

data "ncloud_vpc_peering" "by_source_name" {
	source_vpc_name = ncloud_vpc.main.name
	depends_on      = [ncloud_vpc_peering.foo]
}

data "ncloud_vpc_peering" "by_target_name" {
	target_vpc_name = ncloud_vpc.peer.name
	depends_on      = [ncloud_vpc_peering.foo]
}
`, name)
}
