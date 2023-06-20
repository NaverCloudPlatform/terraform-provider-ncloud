package server

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServerProduct_classic_basic(t *testing.T) {
	dataName := "data.ncloud_server_product.test1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductConfig("SPSW0LINUX000045", "SPSVRSTAND000004"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "server_image_product_code", "SPSW0LINUX000045"),
					resource.TestCheckResourceAttr(dataName, "product_code", "SPSVRSTAND000004"),
					resource.TestCheckResourceAttr(dataName, "product_name", "vCPU 2EA, Memory 4GB, Disk 50GB"),
					resource.TestCheckResourceAttr(dataName, "product_description", "vCPU 2개, 메모리 4GB, 디스크 50GB"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "SVR"),
					resource.TestCheckResourceAttr(dataName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(dataName, "memory_size", "4GB"),
					resource.TestCheckResourceAttr(dataName, "disk_type", "NET"),
					resource.TestCheckResourceAttr(dataName, "generation_code", "G1"),
					resource.TestCheckResourceAttr(dataName, "base_block_storage_size", "50GB"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_server_product.test1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050", "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "server_image_product_code", "SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
					resource.TestCheckResourceAttr(dataName, "product_code", "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"),
					resource.TestCheckResourceAttr(dataName, "product_name", "vCPU 2EA, Memory 8GB, Disk 50GB"),
					resource.TestCheckResourceAttr(dataName, "product_description", "vCPU 2EA, Memory 8GB, Disk 50GB"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "VSVR"),
					resource.TestCheckResourceAttr(dataName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(dataName, "memory_size", "8GB"),
					resource.TestCheckResourceAttr(dataName, "disk_type", "NET"),
					resource.TestCheckResourceAttr(dataName, "generation_code", "G2"),
					resource.TestCheckResourceAttr(dataName, "base_block_storage_size", "50GB"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_classic_FilterByProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductCodeConfig("SPSW0LINUX000045", "SPSVRSTAND000056"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_product.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_vpc_FilterByProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductCodeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050", "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_product.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_classic_FilterByProductNameProductType(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig("SPSW0LINUX000045", "G1"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_product.test3"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_vpc_FilterByProductNameProductType(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050", "G2"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_product.test3"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudServerProductConfig(imageProductCode, productCode string) string {
	return fmt.Sprintf(`
data "ncloud_server_product" "test1" {
	server_image_product_code = "%s"
	product_code = "%s"
}
`, imageProductCode, productCode)
}

func testAccDataSourceNcloudServerProductFilterByProductCodeConfig(imageProductCode, productCode string) string {
	return fmt.Sprintf(`
 data "ncloud_server_product" "test2" {
	server_image_product_code = "%s"
	filter {
		name = "product_code"
		values = ["%s"]
	}
}
`, imageProductCode, productCode)
}

func testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig(imageProductCode, generation string) string {
	return fmt.Sprintf(`
data "ncloud_server_product" "test3" {
	server_image_product_code = "%[1]s"

	filter {
		name = "product_code"
		values = ["SSD"]
		regex = true
	}

	filter {
		name = "cpu_count"
		values = ["2"]
	}

	filter {
		name = "memory_size"
		values = ["8GB"]
	}

	filter {
		name = "base_block_storage_size"
		values = ["50GB"]
	}

	filter {
		name = "product_type"
		values = ["STAND"]
	}

	filter {
		name = "generation_code"
		values = ["%[2]s"]
	}
}`, imageProductCode, generation)
}
