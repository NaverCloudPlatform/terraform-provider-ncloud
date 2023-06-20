package nasvolume

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudNasVolumes_classic_basic(t *testing.T) {
	testAccDataSourceNcloudNasVolumesBasic(t, false)
}

func TestAccDataSourceNcloudNasVolumes_vpc_basic(t *testing.T) {
	testAccDataSourceNcloudNasVolumesBasic(t, true)
}

func testAccDataSourceNcloudNasVolumesBasic(t *testing.T, isVpc bool) {
	postfix := GetTestPrefix()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(isVpc),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNasVolumesConfig(postfix),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_nas_volumes.by_id"),
					TestAccCheckDataSourceID("data.ncloud_nas_volumes.by_filter"),
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
