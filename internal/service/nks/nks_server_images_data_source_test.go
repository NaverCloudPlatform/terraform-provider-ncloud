package nks_test

import (
	"regexp"
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
					resource.TestMatchResourceAttr(dataName, "images.0.value", regexp.MustCompile("SW.VSVR.OS.LNX64.UBNTU.SVR24.WRKND.G003")),
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
  hypervisor_code = "KVM"
  filter {
    name = "label"
    values = ["ubuntu-24.04"]
    regex = true
  }
}
`
}
