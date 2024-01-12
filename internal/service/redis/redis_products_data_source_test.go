package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudRedisProducts_vpc_basic(t *testing.T) {
	productType := "STAND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRedisProductsConfig(productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_redis_products.all", "product_list.0.infra_resource_type", "VRDS"),
					resource.TestCheckResourceAttr("data.ncloud_redis_products.all", "product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudRedisProductsConfig(productType string) string {
	return fmt.Sprintf(`
data "ncloud_redis_products" "all" {
    redis_image_product_code = "SW.VDBAS.VRDS.LNX64.CNTOS.0708.REDIS.7013.B050"
	filter {
			name = "product_type"
			values = ["%[1]s"]
	}
}
`, productType)
}
