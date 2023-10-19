package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSubnetsBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSubnetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_subnets.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSubnetsName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSubnetsConfigSubnet("10.2.1.0"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_subnets.by_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSubnetsVpcNo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSubnetsConfigVpcNo("502"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_subnets.by_vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSubnetsConfig() string {
	return `
data "ncloud_subnets" "all" {}
`
}

func testAccDataSourceNcloudSubnetsConfigSubnet(cidr string) string {
	return fmt.Sprintf(`
data "ncloud_subnets" "by_cidr" {
	subnet = "%s"
}
`, cidr)
}

func testAccDataSourceNcloudSubnetsConfigVpcNo(vpcNo string) string {
	return fmt.Sprintf(`
data "ncloud_subnets" "by_vpc_no" {
	vpc_no = "%s"
}
`, vpcNo)
}
