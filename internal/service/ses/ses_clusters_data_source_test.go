package ses_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSESClusters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSESClustersConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_ses_clusters.all"),
				),
			},
		},
	})
}

func testAccDataSourceSESClustersConfig() string {
	return `
data "ncloud_ses_clusters" "all" {}
`
}
