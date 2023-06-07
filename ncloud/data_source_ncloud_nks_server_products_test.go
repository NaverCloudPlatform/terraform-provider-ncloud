package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudNKSServerProductCodes(t *testing.T) {
	dataName := "data.ncloud_nks_server_products.product_codes"
	zone := "KR-1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNKSServerProductConfig(zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNKSServerProductConfig(zone string) string {
	return fmt.Sprintf(`
data "ncloud_nks_server_images" "images"{
}

data "ncloud_nks_server_products" "product_codes" {
  software_code = data.ncloud_nks_server_images.images.images[0].value
  zone = "%[1]s"
}

`, zone)
}
