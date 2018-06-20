package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDataSourceNcloudServerImageBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudServerImageDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImageFilterByName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageFilterByNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudServerImageDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerImageFilterByType(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerImageFilterByTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudServerImageDataSourceID("data.ncloud_server_image.test"),
				),
			},
		},
	})
}

func testAccCheckNcloudServerImageDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Server Image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server Image data source ID not set")
		}
		return nil
	}
}

var testAccDataSourceNcloudServerImageConfig = `
data "ncloud_server_image" "test" {
}
`

var testAccDataSourceNcloudServerImageFilterByNameConfig = `
data "ncloud_server_image" "test" {
  "product_name_regex" = "Server.*2016"
}
`

var testAccDataSourceNcloudServerImageFilterByTypeConfig = `
data "ncloud_server_image" "test" {
  "product_type_code" = "WINNT"
}
`
