package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudVpcsBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfig(),
				Check: resource.ComposeTestCheckFunc(
					acctest.TestAccCheckDataSourceID("data.ncloud_vpcs.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigName("test"),
				Check: resource.ComposeTestCheckFunc(
					acctest.TestAccCheckDataSourceID("data.ncloud_vpcs.by_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVpcsVpcNo(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcsConfigVpcNo("446"),
				Check: resource.ComposeTestCheckFunc(
					acctest.TestAccCheckDataSourceID("data.ncloud_vpcs.by_vpc_no"),
					acctest.TestAccCheckDataSourceID("data.ncloud_vpcs.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudVpcsConfig() string {
	return `
data "ncloud_vpcs" "all" {}
`
}

func testAccDataSourceNcloudVpcsConfigName(name string) string {
	return fmt.Sprintf(`
data "ncloud_vpcs" "by_name" {
	name               = "%s"
}
`, name)
}

func testAccDataSourceNcloudVpcsConfigVpcNo(vpcNo string) string {
	return fmt.Sprintf(`
data "ncloud_vpcs" "by_vpc_no" {
	vpc_no = "%[1]s"
}

data "ncloud_vpcs" "by_filter" {
	filter {
		name   = "vpc_no"
		values = ["%[1]s"]
	}
}
`, vpcNo)
}
