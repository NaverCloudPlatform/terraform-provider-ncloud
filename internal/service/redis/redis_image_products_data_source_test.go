package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudRedisImageProducts_vpc_basic(t *testing.T) {
	productType := "LINUX"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudRedisImageProductsConfig(productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_redis_image_products.all", "image_product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudRedisImageProductsConfig(productType string) string {
	return fmt.Sprintf(`
data "ncloud_redis_image_products" "all" {
	filter {
			name = "product_type"
			values = ["%[1]s"]
	}
}
`, productType)
}
