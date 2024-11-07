package nks_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudNKSServerImages(t *testing.T) {
	dataName := "data.ncloud_nks_server_images.images"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNKSServerImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNKSServerImagesFilter(t *testing.T) {
	dataName := "data.ncloud_nks_server_images.filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNKSServerImagestWithFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "images.0.value", "SW.VSVR.OS.LNX64.UBNTU.SVR2004.WRKND.B050"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudNKSServerImagesConfig = `
data "ncloud_nks_server_images" "images" {}
`

func testAccDataSourceNcloudNKSServerImagestWithFilterConfig() string {
	return `
data "ncloud_nks_server_images" "filter"{
  filter {
    name = "label"
    values = ["ubuntu-22.04"]
  }
}
`
}
