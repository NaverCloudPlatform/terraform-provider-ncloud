package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudNKSServerImages(t *testing.T) {
	dataName := "data.ncloud_nks_server_images.images"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNKSServerImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudNKSServerImagesFilter(t *testing.T) {
	dataName := "data.ncloud_nks_server_images.filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNKSServerImagestWithFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
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
	return fmt.Sprintf(`
data "ncloud_nks_server_images" "filter"{
  filter {
    name = "label"
    values = ["ubuntu-20.04-64-server"]
  }
}
`)
}
