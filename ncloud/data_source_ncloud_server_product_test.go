package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					testAccCheckNcloudServerProductDataSourceID("data.ncloud_server_product.test"),
				),
			},
		},
	})
}

func testAccCheckNcloudServerProductDataSourceID(n string) resource.TestCheckFunc {
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

var testAccDataSourceNcloudServerProductConfig = `
data "ncloud_server_product" "test" {
	"server_image_product_code" = "SPSW0LINUX000032"
}
`
