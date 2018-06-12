package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					testAccCheckNcloudNasVolumesDataSourceID("data.ncloud_nas_volumes.volumes"),
				),
			},
		},
	})
}

func testAccCheckNcloudNasVolumesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudNasVolumesConfig = `
data "ncloud_nas_volumes" "volumes" {}
`
