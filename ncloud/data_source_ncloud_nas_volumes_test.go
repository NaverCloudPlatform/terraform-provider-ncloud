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
