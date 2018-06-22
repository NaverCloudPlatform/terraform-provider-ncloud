package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
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
					testAccCheckDataSourceID("data.ncloud_server_products.all"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudServerProductsConfig = `
data "ncloud_server_products" "all" {
	"server_image_product_code" = "SPSW0LINUX000032"
}
`
