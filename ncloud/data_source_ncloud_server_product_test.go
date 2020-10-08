package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudServerProduct_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductConfig("SPSW0LINUX000032", "SPSVRSTAND000056"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test1"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050", "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test1"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_classic_FilterByProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductCodeConfig("SPSW0LINUX000032", "SPSVRSTAND000056"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_vpc_FilterByProductCode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductCodeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050", "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test2"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_classic_FilterByProductNameProductType(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig("SPSW0LINUX000032"),
				// ignore check: `generation_code` is will be added classic env.
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test3"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProduct_vpc_FilterByProductNameProductType(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_product.test3"),
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

func testAccDataSourceNcloudServerProductFilterByProductNameProductTypeConfig(imageProductCode string) string {
	return fmt.Sprintf(`
data "ncloud_server_product" "test3" {
	server_image_product_code = "%s"
	filter {
		name = "product_name"
		values = ["vCPU 2EA, Memory 8GB, Disk 50GB"]
	}

	filter {
		name = "product_type"
		values = ["STAND"]
	}

	filter {
		name = "generation_code"
		values = ["G2"]
	}
}`, imageProductCode)
}
