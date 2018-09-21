package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudNasVolumesBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNasVolumesConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_nas_volumes.volumes"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudNasVolumesConfig = `
data "ncloud_nas_volumes" "volumes" {}
`
