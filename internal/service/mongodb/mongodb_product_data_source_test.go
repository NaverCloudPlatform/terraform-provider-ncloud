package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMongoDbProduct_basic(t *testing.T) {
	dataName := "data.ncloud_mongodb_product.test"
	imageProductCode := "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.4212.CE.B050"
	productCode := "SVR.VMGDB.MNGOS.STAND.C002.M008.NET.SSD.B050.G002"
	infraResourceDetailTypeCode := "MNGOD"
	exclusionProductCode := "SVR.VMGDB.MNGOS.HICPU.C004.M008.NET.SSD.B050.G00"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMongoDbProductConfig(imageProductCode, productCode, infraResourceDetailTypeCode, exclusionProductCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "product_name", "vCPU 2EA, Memory 8GB"),
					resource.TestCheckResourceAttr(dataName, "product_type", "STAND"),
					resource.TestCheckResourceAttr(dataName, "product_description", "vCPU 2개, 메모리 8GB, [SSD]디스크 50GB"),
					resource.TestCheckResourceAttr(dataName, "infra_resource_type", "VMGDB"),
					resource.TestCheckResourceAttr(dataName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(dataName, "memory_size", "8589934592"),
					resource.TestCheckResourceAttr(dataName, "disk_type", "NET"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudMongoDbProductConfig(imageProductCode string, productCode string, infraResourceDetailTypeCode string, exclusionProductCode string) string {
	return fmt.Sprintf(`
data "ncloud_mongodb_product" "test" {
	image_product_code = "%s"
	product_code = "%s"
	infra_resource_detail_type_code = "%s"
	exclusion_product_code = "%s"
}
`, imageProductCode, productCode, infraResourceDetailTypeCode, exclusionProductCode)
}
