package devtools

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourcePipelineTriggerTimeZone_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineTriggerTimeZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcepipeline_trigger_timezone.time_zone"),
				),
			},
		},
	})
}
func TestAccDataSourceNcloudSourcePipelineTriggerTimeZone_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineTriggerTimeZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_sourcepipeline_trigger_timezone.time_zone"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourcePipelineTriggerTimeZoneConfig() string {
	return fmt.Sprintf(`
data "ncloud_sourcepipeline_trigger_timezone" "time_zone" {
}
`)
}
