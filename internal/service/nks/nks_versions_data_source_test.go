package nks_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudVersions(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVersionConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_nks_versions.versions"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVersions_KVM(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVersionConfig_KVM,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_nks_versions.versions"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudVersions_regionalSupport(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVersionConfig_regionalSupportTrue,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_nks_versions.versions"),
					resource.TestCheckResourceAttr("data.ncloud_nks_versions.versions", "regional_support", "true"),
				),
			},
		},
	})
}

// regional_support = false exercises the explicit-false path: GetOkExists must
// distinguish it from an unset value and still send isRegionalSupport=false.
func TestAccDataSourceNcloudVersions_regionalSupportFalse(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVersionConfig_regionalSupportFalse,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_nks_versions.versions"),
					resource.TestCheckResourceAttr("data.ncloud_nks_versions.versions", "regional_support", "false"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudVersionConfig = `
data "ncloud_nks_versions" "versions" {}
`

var testAccDataSourceNcloudVersionConfig_KVM = `
data "ncloud_nks_versions" "versions" {
hypervisor_code = "KVM"
}
`

var testAccDataSourceNcloudVersionConfig_regionalSupportTrue = `
data "ncloud_nks_versions" "versions" {
  regional_support = true
}
`

var testAccDataSourceNcloudVersionConfig_regionalSupportFalse = `
data "ncloud_nks_versions" "versions" {
  regional_support = false
}
`
