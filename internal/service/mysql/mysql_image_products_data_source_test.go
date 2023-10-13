package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlImageProducts_basic(t *testing.T) {
	productCode := "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMysqlImageProductsConfig_basic(productCode),
				Check: resource.ComposeTestCheckFunc(
					acctest.TestAccCheckDataSourceID("data.ncloud_mysql_image_products.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMysqlImageProductsConfig_basic(imageProductCode string) string {
	return fmt.Sprintf(`
data "ncloud_mysql_image_products" "all" {
	product_code = "%s"
	filter {
			name = "product_code"
			values = ["%s"]
	}
}
`, imageProductCode, imageProductCode)
}
