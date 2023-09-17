package ses_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSESNodeOsImages(t *testing.T) {
	dataName := "data.ncloud_ses_node_os_images.images"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSESNodeOsImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSESNodeOsImagesFilter(t *testing.T) {
	dataName := "data.ncloud_ses_node_os_images.filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSESNodeOsImagestWithFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "versions.0.id", "SW.VELST.OS.LNX64.CNTOS.0708.B050"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudSESNodeOsImagesConfig = `
data "ncloud_ses_node_os_images" "images" {}
`

func testAccDataSourceNcloudSESNodeOsImagestWithFilterConfig() string {
	return fmt.Sprintf(`
data "ncloud_ses_node_os_images" "filter" {
	filter {
		name = "id"
		values = ["SW.VELST.OS.LNX64.CNTOS.0708.B050"]
	}
}
`)
}
