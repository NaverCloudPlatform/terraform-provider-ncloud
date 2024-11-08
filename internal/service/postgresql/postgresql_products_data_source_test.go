package postgresql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudPostgresqlProducts_basic(t *testing.T) {
	productType := "STAND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPostgresqlProductsConfig_basic(productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_postgresql_products.all", "product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudPostgresqlProductsConfig_basic(productType string) string {
	return fmt.Sprintf(`
data "ncloud_postgresql_image_products" "all" { }

data "ncloud_postgresql_products" "all" {
    image_product_code = data.ncloud_postgresql_image_products.all.image_product_list.0.product_code

	filter {
		name   = "product_type"
		values = ["%s"]
	}
}
`, productType)
}
