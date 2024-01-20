package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMysqlImageProducts_basic(t *testing.T) {
	productType := "LINUX"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMysqlImageProductsConfig_basic(productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_mysql_image_products.all", "image_product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMysqlImageProductsConfig_basic(productType string) string {
	return fmt.Sprintf(`
data "ncloud_mysql_image_products" "all" {
	filter {
			name = "product_type"
			values = ["%s"]
	}
}
`, productType)
}
