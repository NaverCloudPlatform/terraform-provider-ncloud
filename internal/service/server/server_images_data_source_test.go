package server_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServerImages_classic_basic(t *testing.T) {
	// Stop testing classic images
	t.Skip()

	testAccDataSourceNcloudServerImagesBasic(t, false)
}

func TestAccDataSourceNcloudServerImages_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudServerImagesBasic(t, false)
}

func testAccDataSourceNcloudServerImagesBasic(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_images.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImages_classic_linux(t *testing.T) {
	// Stop testing classic images
	t.Skip()

	testAccDataSourceNcloudServerImagesLinux(t, false)
}

func TestAccDataSourceNcloudServerImages_vpc_linux(t *testing.T) {
	testAccDataSourceNcloudServerImagesLinux(t, false)
}

func testAccDataSourceNcloudServerImagesLinux(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesLinuxConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_images.linux"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImages_classic_windows(t *testing.T) {
	// Stop testing classic images
	t.Skip()

	testAccDataSourceNcloudServerImagesWindows(t, false)
}

func TestAccDataSourceNcloudServerImages_vpc_windows(t *testing.T) {
	testAccDataSourceNcloudServerImagesWindows(t, false)
}

func testAccDataSourceNcloudServerImagesWindows(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesWindowsConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_images.windows"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImages_classic_bareMetal(t *testing.T) {
	// Stop testing classic images
	t.Skip()

	testAccDataSourceNcloudServerImagesBareMetal(t, false)
}

func TestAccDataSourceNcloudServerImages_vpc_bareMetal(t *testing.T) {
	testAccDataSourceNcloudServerImagesBareMetal(t, false)
}

func testAccDataSourceNcloudServerImagesBareMetal(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesBareMetalConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_images.beremetal"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImages_classic_blockStorageSize(t *testing.T) {
	// Stop testing classic images
	t.Skip()

	testAccDataSourceNcloudServerImagesBlockStorageSize(t, false)
}

func TestAccDataSourceNcloudServerImages_vpc_blockStorageSize(t *testing.T) {
	testAccDataSourceNcloudServerImagesBlockStorageSize(t, false)
}

func testAccDataSourceNcloudServerImagesBlockStorageSize(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesBlockStorageSizeConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_images.blockstorage"),
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

var testAccDataSourceNcloudServerImagesBlockStorageSizeConfig = `
data "ncloud_server_images" "blockstorage" {
	block_storage_size = 50
}
`
