package postgresql_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudPostgresqlImageProducts_basic(t *testing.T) {
	dataName := "data.ncloud_postgresql_image_products.all"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePostgresqlImageProductsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataName, "image_product_list.0.product_code", regexp.MustCompile(`^[A-Z0-9]+[A-Z0-9-.]+[A-Z0-9]$`)),
					resource.TestMatchResourceAttr(dataName, "image_product_list.1.product_code", regexp.MustCompile(`^[A-Z0-9]+[A-Z0-9-.]+[A-Z0-9]$`)),
				),
			},
		},
	})
}

var testAccDataSourcePostgresqlImageProductsConfig = `
data "ncloud_postgresql_image_products" "all" { }
`
