package memberserverimage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMemberServerImageBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImageConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return SkipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_member_server_image.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudMemberServerImageFilter(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImageConfigFilter,
				SkipFunc: func() (bool, error) {
					return SkipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_member_server_image.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudMemberServerImageConfig = `
data "ncloud_member_server_image" "test" {}
`

var testAccDataSourceNcloudMemberServerImageConfigFilter = `
data "ncloud_member_server_image" "member_server_images" {
	filter {
		name   = "platform_type"
		values = ["LNX64"]
	}
}
`
