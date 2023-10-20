package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlImageProduct_basic(t *testing.T) {
	dataName := "data.ncloud_mysql_image_product.test1"
	productCode := "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMysqlImageProductConfig_basic(productCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "product_name", "mysql(5.7.32)"),
					resource.TestCheckResourceAttr(dataName, "product_type", "LINUX"),
					resource.TestCheckResourceAttr(dataName, "product_description", "CentOS 7.8 with MySQL 5.7.32"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "SW"),
					resource.TestCheckResourceAttr(dataName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(dataName, "os_information", "CentOS 7.8 with MySQL 5.7.32 (64-bit)"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMysqlImageProductConfig_basic(imageProductCode string) string {
	return fmt.Sprintf(`
data "ncloud_mysql_image_product" "test1" {
	product_code = "%s"
}
`, imageProductCode)
}
