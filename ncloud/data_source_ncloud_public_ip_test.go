package ncloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudPublicIp_classic_basic(t *testing.T) {
	testAccDataSourceNcloudPublicIpBasic(t, false)
}

func TestAccDataSourceNcloudPublicIp_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudPublicIpBasic(t, true)
}

func testAccDataSourceNcloudPublicIpBasic(t *testing.T, isVpc bool) {
	resourceName := "ncloud_public_ip.public_ip"
	dataName := "data.ncloud_public_ip.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "server_name", resourceName, "server_name"),

					// Classic only
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "internet_line_type", resourceName, "internet_line_type"),
					resource.TestCheckResourceAttrPair(dataName, "kind_type", resourceName, "kind_type"),
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
						"server_instance.server_instance_no",
					),
				),
			},
		},
	})
}

var testAccDataSourceNcloudPublicIpConfig = `
resource "ncloud_public_ip" "public_ip" {}
data "ncloud_public_ip" "test" {
	public_ip_no = ncloud_public_ip.public_ip.id
}
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
		values = ["tf-2807-1"]
	}
}
`
