package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudServerImagesBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_images.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImagesLinux(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesLinuxConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_images.linux"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImagesWindows(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesWindowsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_images.windows"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImagesBareMetal(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesBareMetalConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_images.beremetal"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudServerImagesConfig = `
data "ncloud_server_images" "test" {}
`

var testAccDataSourceNcloudServerImagesLinuxConfig = `
data "ncloud_server_images" "linux" {
	platform_type_code_list = ["LNX32", "LNX64"]
}
`

var testAccDataSourceNcloudServerImagesWindowsConfig = `
data "ncloud_server_images" "windows" {
	platform_type_code_list = ["WND32", "WND64"]
}
`

var testAccDataSourceNcloudServerImagesBareMetalConfig = `
data "ncloud_server_images" "beremetal" {
	infra_resource_detail_type_code = "BM"
}
`
