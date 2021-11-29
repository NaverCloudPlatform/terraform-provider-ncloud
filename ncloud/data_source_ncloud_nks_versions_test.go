package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudVersionBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVersionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_nks_version.versions"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudVersionConfig = `
data "ncloud_nks_version" "versions" {}
`
