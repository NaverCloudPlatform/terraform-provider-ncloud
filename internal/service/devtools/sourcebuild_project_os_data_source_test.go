package devtools_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceBuildOs(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildOsConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcebuild_project_os.os"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceBuildOsConfig() string {
	return `
data "ncloud_sourcebuild_project_os" "os" {}
`
}
