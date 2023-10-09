package cloudmongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMongoDbImageProduct_basic(t *testing.T) {
	dataName := "data.ncloud_mongodb_image_product.test"
	productCode := "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.4212.CE.B050"
	generationCode := "G2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMongoDbImageProductConfig_basic(productCode, generationCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "product_name", "MongoDB 4.2.12 Community Edition"),
					resource.TestCheckResourceAttr(dataName, "product_type", "LINUX"),
					resource.TestCheckResourceAttr(dataName, "product_description", "CentOS 7.8 with MongoDB 4.2.12 Community Edition"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "SW"),
					resource.TestCheckResourceAttr(dataName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(dataName, "os_information", "CentOS 7.8 with MongoDB 4.2.12 Community Edition (64-bit)"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMongoDbImageProductConfig_basic(productCode string, generationCode string) string {
	return fmt.Sprintf(`
data "ncloud_mongodb_image_product" "test" {
	product_code = "%s"
	generation_code = "%s"
}
`, productCode, generationCode)
}
