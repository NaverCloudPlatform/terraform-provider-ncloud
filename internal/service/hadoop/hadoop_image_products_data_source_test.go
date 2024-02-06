package hadoop_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudHadoopImages_basic(t *testing.T) {
	dataName := "data.ncloud_hadoop_image_products.images"
	productName := "Cloud Hadoop 2.1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHadoopImagesConfig(productName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "image_product_list.0.product_name", "Cloud Hadoop 2.1"),
					resource.TestCheckResourceAttr(dataName, "image_product_list.0.product_type", "LINUX"),
				),
			},
		},
	})
}

func testAccDataSourceHadoopImagesConfig(productName string) string {
	return fmt.Sprintf(`
data "ncloud_hadoop_image_products" "images" {
	filter {
			name = "product_name"
			values = ["%s"]
	}
}
`, productName)
}
