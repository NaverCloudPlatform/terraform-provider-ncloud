package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudServerImage_classic_byCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByCodeConfig("SPSW0LINUX000139"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test1"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_vpc_byCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByCodeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test1"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_classic_byFilterProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByFilterProductCodeConfig("SPSW0LINUX000139"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_vpc_byFilterProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByFilterProductCodeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_classic_byFilterProductName(t *testing.T) {
	testAccDataSourceNcloudServerImageByFilterProductName(t, false)
}

func TestAccDataSourceNcloudServerImage_vpc_byFilterProductName(t *testing.T) {
	testAccDataSourceNcloudServerImageByFilterProductName(t, true)
}

func testAccDataSourceNcloudServerImageByFilterProductName(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
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

func TestAccDataSourceNcloudServerImage_classic_byBlockStorageSize(t *testing.T) {
	testAccDataSourceNcloudServerImageByBlockStorageSize(t, false)
}

func TestAccDataSourceNcloudServerImage_vpc_byBlockStorageSize(t *testing.T) {
	testAccDataSourceNcloudServerImageByBlockStorageSize(t, true)
}

func testAccDataSourceNcloudServerImageByBlockStorageSize(t *testing.T, isVpc bool) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
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

func testAccDataSourceNcloudServerImageByCodeConfig(productCode string) string {
	return fmt.Sprintf(`
data "ncloud_server_image" "test1" {
  product_code = "%s"
}
`, productCode)
}

func testAccDataSourceNcloudServerImageByFilterProductCodeConfig(productCode string) string {
	return fmt.Sprintf(`
data "ncloud_server_image" "test2" {
  filter {
    name = "product_code"
    values = ["%s"]
  }
}
`, productCode)
}

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
