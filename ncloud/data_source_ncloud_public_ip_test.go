package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudPublicIPBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIPConfig,
				// ignore check: may be empty created data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIPMostRecent(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIPMostRecentConfig,
				// ignore check: may be empty created data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIPIsAssociated(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIPAssociatedConfig,
				// ignore check: may be empty created data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIPSearch(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIPSearchConfig,
				// ignore check: may be empty created data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//	resource.TestCheckResourceAttrSet(
				//		"data.ncloud_public_ip.test",
				//		"server_instance.server_instance_no",
				//	),
				//),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIPSorting(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIPSortingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudPublicIPConfig = `
data "ncloud_public_ip" "test" {}
`

var testAccDataSourceNcloudPublicIPMostRecentConfig = `
data "ncloud_public_ip" "test" {
  "most_recent" = "true"
}
`

var testAccDataSourceNcloudPublicIPAssociatedConfig = `
data "ncloud_public_ip" "test" {
  "is_associated" = "false"
  "most_recent" = "true"
}
`

var testAccDataSourceNcloudPublicIPSearchConfig = `
data "ncloud_public_ip" "test" {
  "search_filter_name" = "associatedServerName" // Public IP (publicIp) | Associated server name (associatedServerName)
  "search_filter_value" = "clouddb"
  "most_recent" = "true"
}
`

var testAccDataSourceNcloudPublicIPSortingConfig = `
data "ncloud_public_ip" "test" {
  "sorted_by" = "publicIp" // Public IP (publicIp) | Public IP instance number (publicIpInstanceNo) [case insensitive]
  "sorting_order" = "ascending" // Ascending (ascending) | Descending (descending) [case insensitive]
  "most_recent" = "true"
}
`
