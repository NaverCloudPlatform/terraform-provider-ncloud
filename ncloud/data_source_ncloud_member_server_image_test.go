package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceNcloudMemberServerImageBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudMemberServerImageDataSourceID("data.ncloud_member_server_image.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudMemberServerImageMostRecent(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImageMostRecentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudMemberServerImageDataSourceID("data.ncloud_member_server_image.test"),
				),
			},
		},
	})
}

func testAccCheckNcloudMemberServerImageDataSourceID(n string) resource.TestCheckFunc {
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

var testAccDataSourceNcloudMemberServerImageConfig = `
data "ncloud_member_server_image" "test" {}
`

var testAccDataSourceNcloudMemberServerImageMostRecentConfig = `
data "ncloud_member_server_image" "test" {
  "most_recent" = "true"
}
`
