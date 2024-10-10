package postgresql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudPostgresqlProducts_basic(t *testing.T) {
	imageProductCode := "SW.VPGSL.OS.LNX64.CNTOS.0708.PGSQL.133.B050"
	productType := "STAND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPostgresqlProductsConfig_basic(imageProductCode, productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_postgresql_products.all", "product_list.0.infra_resource_type", "VPGSL"),
					resource.TestCheckResourceAttr("data.ncloud_postgresql_products.all", "product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudPostgresqlProductsConfig_basic(imageProductCode string, productType string) string {
	return fmt.Sprintf(`
data "ncloud_postgresql_products" "all" {
	image_product_code = "%s"
	
	filter {
		name   = "product_type"
		values = ["%s"]
	}
}
`, imageProductCode, productType)
}
