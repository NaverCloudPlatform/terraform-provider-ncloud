package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudServerProductBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudServerProductConfig = `
data "ncloud_server_product" "test" {
	"server_image_product_code" = "SPSW0LINUX000032"
}
`
