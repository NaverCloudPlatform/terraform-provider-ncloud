package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlProduct_basic(t *testing.T) {
	dataName := "data.ncloud_mysql_product.test1"
	imageProductCode := "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"
	productCode := "SVR.VDBAS.STAND.C002.M008.NET.HDD.B050.G002"
	excludeProductCode := "SVR.VDBAS.HICPU.C004.M008.NET.HDD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMysqlProductConfig(imageProductCode, productCode, excludeProductCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "product_name", "vCPU 2EA, Memory 8GB"),
					resource.TestCheckResourceAttr(dataName, "product_type", "STAND"),
					resource.TestCheckResourceAttr(dataName, "product_description", "vCPU 2개, 메모리 8GB"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "VMYSL"),
					resource.TestCheckResourceAttr(dataName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(dataName, "memory_size", "8589934592"),
					resource.TestCheckResourceAttr(dataName, "disk_type", "NET"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMysqlProductConfig(imageProductCode string, productCode string, excludeProductCode string) string {
	return fmt.Sprintf(`
data "ncloud_mysql_product" "test1" {
	cloud_mysql_image_product_code = "%s"
	product_code = "%s"
	exclusion_product_code = "%s"
}
`, imageProductCode, productCode, excludeProductCode)
}
