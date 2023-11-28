package hadoop_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoopIamges_basic(t *testing.T) {
	dataName := "data.ncloud_hadoop_images.images"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHadoopImagesConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "images.0.product_name", "Cloud Hadoop 1.5"),
					resource.TestCheckResourceAttr(dataName, "images.0.product_type", "LINUX"),
					resource.TestCheckResourceAttr(dataName, "images.0.product_description", "CentOS 7.8 with Cloud Hadoop 1.5"),
				),
			},
		},
	})
}

func testAccDataSourceHadoopImagesConfig() string {
	return fmt.Sprintf(`
data "ncloud_hadoop_images" "images" {
	product_code = "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
}
`)
}
