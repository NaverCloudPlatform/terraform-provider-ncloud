package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudNasVolumeBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNasVolumeConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_nas_volume.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudNasVolumeConfig = `
data "ncloud_nas_volume" "test" {}
`
