package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudMemberServerImagesBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImagesConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_member_server_images.member_server_images"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudMemberServerImagesFilter(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImagesConfigFilter,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_member_server_images.member_server_images"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudMemberServerImagesConfig = `
data "ncloud_member_server_images" "member_server_images" {}
`

var testAccDataSourceNcloudMemberServerImagesConfigFilter = `
data "ncloud_member_server_images" "member_server_images" {
	filter {
		name   = "platform_type"
		values = ["LNX64"]
	}
}
`
