package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudServerImageByCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByCodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}
func TestAccDataSourceNcloudServerImageByFilterProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByFilterProductCodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImageByFilterProductName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByFilterProductNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudServerImageByCodeConfig = `
data "ncloud_server_image" "test" {
  product_code = "SPSW0LINUX000139"
}
`

var testAccDataSourceNcloudServerImageByFilterProductCodeConfig = `
data "ncloud_server_image" "test" {
  filter {
    name = "product_code"
    values = ["SPSW0LINUX000139"]
  }
}
`

var testAccDataSourceNcloudServerImageByFilterProductNameConfig = `
data "ncloud_server_image" "test" {
  filter {
    name = "product_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}
`
