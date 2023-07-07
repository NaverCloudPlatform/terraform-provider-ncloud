package memberserverimage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMemberServerImagesBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImagesConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return SkipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_member_server_images.member_server_images"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudMemberServerImagesFilter(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImagesConfigFilter,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return SkipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_member_server_images.member_server_images"),
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
