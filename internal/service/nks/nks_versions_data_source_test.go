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

func TestAccDataSourceNcloudVersions_XEN(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVersionConfig_XEN,
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

var testAccDataSourceNcloudVersionConfig = `
data "ncloud_nks_versions" "versions" {}
`

var testAccDataSourceNcloudVersionConfig_XEN = `
data "ncloud_nks_versions" "versions" {
hypervisor_code = "XEN"
}
`

var testAccDataSourceNcloudVersionConfig_KVM = `
data "ncloud_nks_versions" "versions" {
hypervisor_code = "KVM"
}
`
