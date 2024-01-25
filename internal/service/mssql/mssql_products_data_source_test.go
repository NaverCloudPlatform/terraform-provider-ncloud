package mssql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMssqlProducts_basic(t *testing.T) {
	imageProductCode := "SW.VMSSL.OS.WND64.WINNT.SVR2016.MSSQL.15042981.SE.B100"
	productType := "STAND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMssqlProductsConfig_basic(imageProductCode, productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_mssql_products.all", "product_list.0.infra_resource_type", "VMSSL"),
					resource.TestCheckResourceAttr("data.ncloud_mssql_products.all", "product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMssqlProductsConfig_basic(imageProductCode string, productType string) string {
	return fmt.Sprintf(`
data "ncloud_mssql_products" "all" {
	image_product_code = "%s"
	
	filter {
		name   = "product_type"
		values = ["%s"]
	}
}
`, imageProductCode, productType)
}
