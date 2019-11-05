package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
	server_image_product_code = "SPSW0LINUX000032"
	product_name_regex = "vCPU 1EA, Memory 1GB, Disk 50GB"
}
`
