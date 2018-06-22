package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
					testAccCheckDataSourceID("data.ncloud_member_server_image.test"),
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
					testAccCheckDataSourceID("data.ncloud_member_server_image.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudMemberServerImageConfig = `
data "ncloud_member_server_image" "test" {}
`

var testAccDataSourceNcloudMemberServerImageMostRecentConfig = `
data "ncloud_member_server_image" "test" {
  "most_recent" = "true"
}
`
