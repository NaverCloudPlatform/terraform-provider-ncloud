package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudServerProducts_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { testAccPreCheck(t) },
		IsUnitTest: false,
		Providers:  testAccClassicProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductsConfig("SPSW0LINUX000032"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_products.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProducts_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { testAccPreCheck(t) },
		IsUnitTest: false,
		Providers:  testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductsConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_products.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudServerProductsConfig(productCode string) string {
	return fmt.Sprintf(`
data "ncloud_server_products" "all" {
	server_image_product_code = "%s"
}
`, productCode)
}
