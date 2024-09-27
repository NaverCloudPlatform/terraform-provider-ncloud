package server_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServerImageNumbers_basic(t *testing.T) {
	dataName := "data.ncloud_server_image_numbers.images"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerImageNumbersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataName, "image_number_list.0.server_image_number", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataName, "image_number_list.1.server_image_number", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataName, "image_number_list.2.server_image_number", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

var testAccDataSourceServerImageNumbersConfig = `
data "ncloud_server_image_numbers" "images" { }
`
