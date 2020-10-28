package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNcloudPlacementGroup_basic(t *testing.T) {
	name := fmt.Sprintf("tf-pl-group-data-%s", acctest.RandString(5))
	resourceName := "ncloud_placement_group.foo"
	dataName := "data.ncloud_placement_group.by_id"
	dataNameFilter := "data.ncloud_placement_group.by_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPlacementGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					testAccCheckDataSourceID(dataNameFilter),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "placement_group_no", resourceName, "placement_group_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "placement_group_type", resourceName, "placement_group_type"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudPlacementGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_placement_group" "foo" {
	name = "%[1]s"
}

data "ncloud_placement_group" "by_id" {
	id = ncloud_placement_group.foo.id
}

data "ncloud_placement_group" "by_filter" {
	filter {
		name   = "id"
		values = [ncloud_placement_group.foo.id]
	}
}
`, name)
}
