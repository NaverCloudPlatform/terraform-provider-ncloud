package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudSESSoftwareProductCodes(t *testing.T) {
	dataName := "data.ncloud_ses_software_product.codes"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSESSoftwareProductConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSESSoftwareProductCodesFilter(t *testing.T) {
	dataName := "data.ncloud_ses_software_product.filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSESSoftwareProductWithFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "codes.0.value", "SW.VELST.OS.LNX64.CNTOS.0708.B050"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudSESSoftwareProductConfig = `
data "ncloud_ses_software_product" "codes" {}
`

func testAccDataSourceNcloudSESSoftwareProductWithFilterConfig() string {
	return fmt.Sprintf(`
data "ncloud_ses_software_product" "filter" {
	filter {
		name = "value"
		values = ["SW.VELST.OS.LNX64.CNTOS.0708.B050"]
	}
}
`)
}
