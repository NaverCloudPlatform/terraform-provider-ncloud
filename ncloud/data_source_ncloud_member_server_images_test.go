package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceNcloudMemberServerImagesBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudMemberServerImagesDataSourceID("data.ncloud_member_server_images.member_server_images"),
				),
			},
		},
	})
}

func testAccCheckNcloudMemberServerImagesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find member server image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("member server Image data source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudMemberServerImagesConfig = `
data "ncloud_member_server_images" "member_server_images" {}
`
