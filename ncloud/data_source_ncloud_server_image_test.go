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
					testAccCheckDataSourceID("data.ncloud_server_image.test1"),
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
					testAccCheckDataSourceID("data.ncloud_server_image.test2"),
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
					testAccCheckDataSourceID("data.ncloud_server_image.test3"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImageByBlockStorageSize(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByBlockStorageSizeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test4"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudServerImageByCodeConfig = `
data "ncloud_server_image" "test1" {
  product_code = "SPSW0LINUX000139"
}
`

var testAccDataSourceNcloudServerImageByFilterProductCodeConfig = `
data "ncloud_server_image" "test2" {
  filter {
    name = "product_code"
    values = ["SPSW0LINUX000139"]
  }
}
`

var testAccDataSourceNcloudServerImageByFilterProductNameConfig = `
data "ncloud_server_image" "test3" {
  filter {
    name = "product_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}
`

var testAccDataSourceNcloudServerImageByBlockStorageSizeConfig = `
data "ncloud_server_image" "test4" {
	block_storage_size = 50
	filter {
    name = "product_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}
`
