package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourcePipelineTriggerTimeZone_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineTriggerTimeZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcepipeline_trigger_timezone.time_zone"),
				),
			},
		},
	})
}
func TestAccDataSourceNcloudSourcePipelineTriggerTimeZone_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineTriggerTimeZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcepipeline_trigger_timezone.time_zone"),
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
