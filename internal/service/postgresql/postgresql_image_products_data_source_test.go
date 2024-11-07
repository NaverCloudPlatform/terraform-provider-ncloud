package postgresql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudPostgresqlImageProducts_basic(t *testing.T) {
	productType := "LINUX"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPostgresqlImageProductsConfig_basic(productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_postgresql_image_products.all", "image_product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudPostgresqlImageProductsConfig_basic(productType string) string {
	return fmt.Sprintf(`
data "ncloud_postgresql_image_products" "all" {
	filter {
			name = "product_type"
			values = ["%s"]
	}
}
`, productType)
}
