package ncloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNKSClustersBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNKSClustersConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_nks_clusters.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNKSClustersConfig() string {
	return `
data "ncloud_nks_clusters" "all" {}
`
}
