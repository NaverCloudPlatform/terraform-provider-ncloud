package hadoop_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoopProducts_basic(t *testing.T) {
	imageProductCode := "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
	productCode := "SVR.VCHDP.EDGND.HIMEM.C004.M032.NET.HDD.B050.G002"
	infraResourceDetailTypeCode := "EDGND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHadoopProductsConfig(imageProductCode, productCode, infraResourceDetailTypeCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_hadoop_products.products", "product_list.0.infra_resource_type", "VCHDP"),
					resource.TestCheckResourceAttr("data.ncloud_hadoop_products.products", "product_list.0.product_code", productCode),
					resource.TestCheckResourceAttr("data.ncloud_hadoop_products.products", "product_list.0.infra_resource_detail_type", infraResourceDetailTypeCode),
				),
			},
		},
	})
}

func testAccDataSourceHadoopProductsConfig(imageProductCode, productCode, infraResourceDetailTypeCode string) string {
	return fmt.Sprintf(`
data "ncloud_hadoop_products" "products" {
	image_product_code = "%[1]s"
	product_code = "%[2]s"
	infra_resource_detail_type_code = "%[3]s"
}
`, imageProductCode, productCode, infraResourceDetailTypeCode)
}
