package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDataSourceServerImages_basic(t *testing.T) {
	// t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudServerImagesDataSourceID("data.ncloud_server_images.server_images"),
				),
			},
		},
	})
}

func testAccCheckNcloudServerImagesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Server Image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server Image data source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudServerImagesConfig = `
data "ncloud_server_images" "server_images" {}
`
