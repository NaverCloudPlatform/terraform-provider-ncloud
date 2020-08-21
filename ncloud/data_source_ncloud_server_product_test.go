package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudServerProductBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test1"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProductFilterByProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductCodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProductFilterByProductNameProductType(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test3"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudServerProductConfig = `
data "ncloud_server_product" "test1" {
	server_image_product_code = "SPSW0LINUX000032"
  product_code = "SPSVRSTAND000056"
}
`

var testAccDataSourceNcloudServerProductFilterByProductCodeConfig = `
data "ncloud_server_product" "test2" {
	server_image_product_code = "SPSW0LINUX000032"
  filter {
    name = "product_code"
    values = ["SPSVRSTAND000056"]
	}
}
`

var testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig = `
data "ncloud_server_product" "test3" {
	server_image_product_code = "SPSW0LINUX000032"
  filter {
    name = "product_name"
    values = ["vCPU 1EA, Memory 1GB, Disk 50GB"]
  }

  filter {
    name = "product_type"
    values = ["MICRO"]
  }
}
`
