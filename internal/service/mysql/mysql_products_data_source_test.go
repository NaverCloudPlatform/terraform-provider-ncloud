package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlProducts_basic(t *testing.T) {
	imageProductCode := "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"
	productCode := "SVR.VDBAS.STAND.C002.M008.NET.HDD.B050.G002"
	excludeProductCode := "SVR.VDBAS.HICPU.C004.M008.NET.HDD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMysqlProductsConfig_basic(imageProductCode, productCode, excludeProductCode),
				Check: resource.ComposeTestCheckFunc(
					acctest.TestAccCheckDataSourceID("data.ncloud_mysql_products.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMysqlProductsConfig_basic(imageProductCode string, productCode string, excludeProductCode string) string {
	return fmt.Sprintf(`
data "ncloud_mysql_products" "all" {
	cloud_mysql_image_product_code = "%s"
	product_code = "%s"
	exclusion_product_code = "%s"
	
	filter {
		name   = "product_type"
		values = ["STAND"]
	}
}
`, imageProductCode, productCode, excludeProductCode)
}
