package ncloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudPublicIpBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIpMostRecent(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpMostRecentConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIpIsAssociated(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpAssociatedConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIpSearch(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpSearchConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
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

func TestAccDataSourceNcloudPublicIpSorting(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpSortingConfig,
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				//ExpectError: regexp.MustCompile("no results"), // may be no data
				//Check: resource.ComposeTestCheckFunc(
				//	testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				//),
			},
		},
	})
}

var testAccDataSourceNcloudPublicIpConfig = `
data "ncloud_public_ip" "test" {}
`

var testAccDataSourceNcloudPublicIpMostRecentConfig = `
data "ncloud_public_ip" "test" {
  "most_recent" = "true"
}
`

var testAccDataSourceNcloudPublicIpAssociatedConfig = `
data "ncloud_public_ip" "test" {
  "is_associated" = "false"
  "most_recent" = "true"
}
`

var testAccDataSourceNcloudPublicIpSearchConfig = `
data "ncloud_public_ip" "test" {
  "search_filter_name" = "associatedServerName" // Public IP (publicIp) | Associated server name (associatedServerName)
  "search_filter_value" = "clouddb"
  "most_recent" = "true"
}
`

var testAccDataSourceNcloudPublicIpSortingConfig = `
data "ncloud_public_ip" "test" {
  "sorted_by" = "publicIp" // Public IP (publicIp) | Public IP instance number (publicIpInstanceNo) [case insensitive]
  "sorting_order" = "ascending" // Ascending (ascending) | Descending (descending) [case insensitive]
  "most_recent" = "true"
}
`
