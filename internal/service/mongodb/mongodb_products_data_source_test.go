package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMongoDbProducts_basic(t *testing.T) {
	imageProductCode := "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.4223.CE.B050"
	infraResourceDetailTypeCode := "MNGOD"
	productType := "STAND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMongoDbProductsConfig_basic(imageProductCode, infraResourceDetailTypeCode, productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_mongodb_products.all", "product_list.0.infra_resource_type", "VMGDB"),
					resource.TestCheckResourceAttr("data.ncloud_mongodb_products.all", "product_list.0.infra_resource_detail_type", infraResourceDetailTypeCode),
					resource.TestCheckResourceAttr("data.ncloud_mongodb_products.all", "product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMongoDbProductsConfig_basic(imageProductCode string, infraResourceDetailTypeCode string, productType string) string {
	return fmt.Sprintf(`
data "ncloud_mongodb_products" "all" {
	image_product_code = "%s"
	infra_resource_detail_type_code = "%s"

	filter {
		name = "product_type"
		values = ["%s"]
	}
}
`, imageProductCode, infraResourceDetailTypeCode, productType)
}
