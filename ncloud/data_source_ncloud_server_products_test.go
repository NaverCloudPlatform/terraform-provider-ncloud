package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDataSourceNcloudServerProductsBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudServerProductsDataSourceID("data.ncloud_server_products.all"),
				),
			},
		},
	})
}

func testAccCheckNcloudServerProductsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find server product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server product data source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudServerProductsConfig = `
data "ncloud_server_products" "all" {
	"server_image_product_code" = "SPSW0LINUX000032"
}
`
