package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceCommitRepositories(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceCommitRepositoriesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcecommit_repositories.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceCommitRepositoriesConfig() string {
	return `
data "ncloud_sourcecommit_repositories" "all" {}
`
}
