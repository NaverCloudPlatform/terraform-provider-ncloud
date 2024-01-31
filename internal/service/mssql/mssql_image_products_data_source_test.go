package mssql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMssqlImageProducts_basic(t *testing.T) {
	productType := "WINNT"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMssqlImageProductsConfig_basic(productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_mssql_image_products.all", "image_product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMssqlImageProductsConfig_basic(productType string) string {
	return fmt.Sprintf(`
data "ncloud_mssql_image_products" "all" {
	filter {
			name = "product_type"
			values = ["%s"]
	}
}
`, productType)
}
