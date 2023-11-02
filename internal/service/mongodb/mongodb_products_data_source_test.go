package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMongoDbProducts_basic(t *testing.T) {
	imageProductCode := "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.4212.CE.B050"
	productCode := "SVR.VMGDB.MNGOS.STAND.C002.M008.NET.SSD.B050.G002"
	infraResourceDetailTypeCode := "MNGOD"
	exclusionProductCode := "SVR.VMGDB.MNGOS.HICPU.C004.M008.NET.SSD.B050.G00"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMongoDbProductsConfig_basic(imageProductCode, productCode, infraResourceDetailTypeCode, exclusionProductCode),
				Check: resource.ComposeTestCheckFunc(
					acctest.TestAccCheckDataSourceID("data.ncloud_mongodb_products.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMongoDbProductsConfig_basic(imageProductCode string, productCode string, infraResourceDetailTypeCode string, exclusionProductCode string) string {
	return fmt.Sprintf(`
data "ncloud_mongodb_products" "all" {
	image_product_code = "%s"
	product_code = "%s"
	infra_resource_detail_type_code = "%s"
	exclusion_product_code = "%s"

	filter {
		name = "product_type"
		values = ["LINUX"]
	}
}
`, imageProductCode, productCode, infraResourceDetailTypeCode, exclusionProductCode)
}
