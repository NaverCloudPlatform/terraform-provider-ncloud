package devtools_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourceCommitRepositories(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceCommitRepositoriesConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcecommit_repositories.all"),
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
