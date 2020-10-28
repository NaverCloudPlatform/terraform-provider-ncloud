package ncloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudBlockStorageSnapshot_basic(t *testing.T) {
	dataName := "data.ncloud_block_storage_snapshot.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// ignore check: may be empty created data
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				Config: testAccDataSourceNcloudBlockStorageSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataName, "snapshot_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataName, "block_storage_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataName, "volume_size", regexp.MustCompile(`^\d+$`)),
					testAccCheckDataSourceID("data.ncloud_block_storage_snapshot.by_filter"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudBlockStorageSnapshotConfig = `
data "ncloud_block_storage_snapshot" "by_id" {
	id = "5192089"
}

data "ncloud_block_storage_snapshot" "by_filter" {
	filter {
		name = "snapshot_no"
		values = ["5192089"]
	}
}
`
