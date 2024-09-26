package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServerImageNumbers_basic(t *testing.T) {
	dataName := "data.ncloud_server_image_numbers.images"
	imageName := "rocky-8.10-gpu"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerImageNumbersConfig(imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataName, "image_number_list.0.name", "rocky-8.10-gpu"),
					resource.TestCheckResourceAttr(dataName, "image_number_list.0.number", "25623982"),
				),
			},
		},
	})
}

func testAccDataSourceServerImageNumbersConfig(imageName string) string {
	return fmt.Sprintf(`
data "ncloud_server_image_numbers" "images" {
	filter {
			name = "name"
			values = ["%s"]
	}
}
`, imageName)
}
