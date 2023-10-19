package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServerImage_classic_byCode(t *testing.T) {
	dataName := "data.ncloud_server_image.test1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByCodeConfig("SPSW0LINUX000046"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "product_code", "SPSW0LINUX000046"),
					resource.TestCheckResourceAttr(dataName, "product_name", "centos-7.3-64"),
					resource.TestCheckResourceAttr(dataName, "product_description", "CentOS 7.3 (64-bit)"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "SW"),
					resource.TestCheckResourceAttr(dataName, "os_information", "CentOS 7.3 (64-bit)"),
					resource.TestCheckResourceAttr(dataName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(dataName, "base_block_storage_size", "50GB"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_vpc_byCode(t *testing.T) {
	dataName := "data.ncloud_server_image.test1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByCodeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "product_code", "SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
					resource.TestCheckResourceAttr(dataName, "product_name", "centos-7.3-64"),
					resource.TestCheckResourceAttr(dataName, "product_description", "CentOS 7.3 (64-bit)"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "SW"),
					resource.TestCheckResourceAttr(dataName, "os_information", "CentOS 7.3 (64-bit)"),
					resource.TestCheckResourceAttr(dataName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(dataName, "base_block_storage_size", "50GB"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_classic_byFilterProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByFilterProductCodeConfig("SPSW0LINUX000139"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_image.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_vpc_byFilterProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByFilterProductCodeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_image.test2"),
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByFilterProductNameConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_image.test3"),
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByBlockStorageSizeConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_image.test4"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImage_vpc_byPlatformType(t *testing.T) {
	dataName := "data.ncloud_server_image.test5"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageByPlatformTypeConfig("LNX64"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "product_code", "SW.VSVR.APP.LNX64.CNTOS.0708.PINPT.LATEST.B050"),
					resource.TestCheckResourceAttr(dataName, "product_name", "Pinpoint-centos-7.8-64"),
					resource.TestCheckResourceAttr(dataName, "product_description", "CentOS 7.8 with Pinpoint"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "SW"),
					resource.TestCheckResourceAttr(dataName, "os_information", "CentOS 7.8 with Pinpoint"),
					resource.TestCheckResourceAttr(dataName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(dataName, "base_block_storage_size", "50GB"),
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
    values = ["centos-7.3-64"]
  }
}
`

var testAccDataSourceNcloudServerImageByBlockStorageSizeConfig = `
data "ncloud_server_image" "test4" {
	filter {
		name = "product_name"
		values = ["centos-7.3-64"]
	}

	filter {
		name = "base_block_storage_size"
		values = ["50GB"]
	}
}
`

func testAccDataSourceNcloudServerImageByPlatformTypeConfig(platformType string) string {
	return fmt.Sprintf(`
data "ncloud_server_image" "test5" {
  product_code = "SW.VSVR.APP.LNX64.CNTOS.0708.PINPT.LATEST.B050"
  platform_type = "%s"
}
`, platformType)
}
