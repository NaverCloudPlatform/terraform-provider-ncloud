package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceNcloudServerImageFilterByName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageFilterByNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImageFilterByType(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageFilterByTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudServerImageFilterByNameConfig = `
data "ncloud_server_image" "test" {
  product_name_regex = "Server.*2016"
}
`

var testAccDataSourceNcloudServerImageFilterByTypeConfig = `
data "ncloud_server_image" "test" {
  product_type = "LINUX"
}
`
