package hadoop_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoopProducts_basic(t *testing.T) {
	imageProductCode := "SW.VCHDP.LNX64.CNTOS.0708.HDP.21.B050"
	infraResourceDetailTypeCode := "EDGND"
	productType := "STAND"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHadoopProductsConfig(imageProductCode, infraResourceDetailTypeCode, productType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ncloud_hadoop_products.products", "product_list.0.infra_resource_type", "VCHDP"),
					resource.TestCheckResourceAttr("data.ncloud_hadoop_products.products", "product_list.0.infra_resource_detail_type", infraResourceDetailTypeCode),
					resource.TestCheckResourceAttr("data.ncloud_hadoop_products.products", "product_list.0.product_type", productType),
				),
			},
		},
	})
}

func testAccDataSourceHadoopProductsConfig(imageProductCode, infraResourceDetailTypeCode, productType string) string {
	return fmt.Sprintf(`
data "ncloud_hadoop_products" "products" {
	image_product_code = "%[1]s"
	infra_resource_detail_type_code = "%[2]s"
	
	filter {
		name = "product_type"
		values = ["%[3]s"]
	}
}
`, imageProductCode, infraResourceDetailTypeCode, productType)
}
