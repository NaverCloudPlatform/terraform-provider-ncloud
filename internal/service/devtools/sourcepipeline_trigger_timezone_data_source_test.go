package devtools_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudSourcePipelineTriggerTimeZone_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ClassicProtoV5ProviderFactories,
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
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
	return `
data "ncloud_sourcepipeline_trigger_timezone" "time_zone" {
}
`
}
