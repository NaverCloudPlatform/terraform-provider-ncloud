package ses_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSESVersions(t *testing.T) {
	dataName := "data.ncloud_ses_versions.versions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSESVersionConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

var testAccDataSourceNcloudSESVersionConfig = `
data "ncloud_ses_versions" "versions" {}
`
