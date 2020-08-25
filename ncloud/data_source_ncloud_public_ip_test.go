package ncloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudPublicIpBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpConfig,
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				ExpectError: regexp.MustCompile("no results"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIpIsAssociated(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpAssociatedConfig,
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				ExpectError: regexp.MustCompile("no results"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIpSearch(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpSearchConfig,
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				ExpectError: regexp.MustCompile("no results"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_public_ip.test"),
					resource.TestCheckResourceAttrSet(
						"data.ncloud_public_ip.test",
						"server_instance.instance_no",
					),
				),
			},
		},
	})
}

var testAccDataSourceNcloudPublicIpConfig = `
data "ncloud_public_ip" "test" {}
`

var testAccDataSourceNcloudPublicIpAssociatedConfig = `
data "ncloud_public_ip" "test" {
	is_associated = "false"
}
`

var testAccDataSourceNcloudPublicIpSearchConfig = `
data "ncloud_public_ip" "test" {
	filter {
		name = "server_instance.server_name"
		values = ["clouddb"]
	}
}
`
