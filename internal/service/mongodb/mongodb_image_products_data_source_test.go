package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMongoDbImageProducts_basic(t *testing.T) {
	productCode := "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.4212.CE.B050"
	generationCode := "G2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMongoDbImageProductsConfig_basic(productCode, generationCode),
				Check: resource.ComposeTestCheckFunc(
					acctest.TestAccCheckDataSourceID("data.ncloud_mongodb_image_product_list.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMongoDbImageProductsConfig_basic(productCode string, generationCode string) string {
	return fmt.Sprintf(`
data "ncloud_mongodb_image_product_list" "all" {
	product_code = "%s"
	generation_code = "%s"
	
	filter {
		name = "product_code"
		values = ["%s"]
	}
}
`, productCode, generationCode, productCode)
}
