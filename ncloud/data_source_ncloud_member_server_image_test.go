package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudMemberServerImageBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudMemberServerImageConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
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
