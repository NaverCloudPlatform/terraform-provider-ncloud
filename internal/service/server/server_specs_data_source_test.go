package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServerSpecs_basic(t *testing.T) {
	dataName := "data.ncloud_server_specs.specs"
	generationCode := "G3"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerSpecsConfig(generationCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "server_spec_list.0.generation_code", "G3"),
					resource.TestCheckResourceAttr(dataName, "server_spec_list.0.hypervisor_type", "KVM"),
				),
			},
		},
	})
}

func testAccDataSourceServerSpecsConfig(generationCode string) string {
	return fmt.Sprintf(`
data "ncloud_server_specs" "specs" {
	filter {
			name = "generation_code"
			values = ["%s"]
	}
}
`, generationCode)
}
