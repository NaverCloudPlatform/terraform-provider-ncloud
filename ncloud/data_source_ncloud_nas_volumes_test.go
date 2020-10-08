package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudNasVolumes_classic_basic(t *testing.T) {
	testAccDataSourceNcloudNasVolumesBasic(t, false)
}

func TestAccDataSourceNcloudNasVolumes_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudNasVolumesBasic(t, true)
}

func testAccDataSourceNcloudNasVolumesBasic(t *testing.T, isVpc bool) {
	postfix := getTestPrefix()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNasVolumesConfig(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_nas_volumes.by_id"),
					testAccCheckDataSourceID("data.ncloud_nas_volumes.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNasVolumesConfig(volumeNamePostfix string) string {
	return fmt.Sprintf(`
resource "ncloud_nas_volume" "test" {
	volume_name_postfix = "%s"
	volume_size = "500"
	volume_allotment_protocol_type = "NFS"
}

data "ncloud_nas_volumes" "by_id" {
	no_list = [ncloud_nas_volume.test.id]
}

data "ncloud_nas_volumes" "by_filter" {
	filter {
		name = "nas_volume_no"
		values = [ncloud_nas_volume.test.id]
	}
}
`, volumeNamePostfix)
}
