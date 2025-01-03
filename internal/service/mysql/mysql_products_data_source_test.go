package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlProducts_basic(t *testing.T) {
	imageProductCode := "SW.VMYSL.OS.LNX64.ROCKY.0810.MYSQL.B050"
	productType := "STAND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMysqlProductsConfig_basic(imageProductCode, productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_mysql_products.all", "product_list.0.infra_resource_type", "VMYSL"),
					resource.TestCheckResourceAttr("data.ncloud_mysql_products.all", "product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMysqlProductsConfig_basic(imageProductCode string, productType string) string {
	return fmt.Sprintf(`
data "ncloud_mysql_products" "all" {
	image_product_code = "%s"
	
	filter {
		name   = "product_type"
		values = ["%s"]
	}
}
`, imageProductCode, productType)
}
