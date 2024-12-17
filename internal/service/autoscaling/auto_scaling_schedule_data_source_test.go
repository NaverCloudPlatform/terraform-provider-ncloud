package autoscaling_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudAutoScalingSchedule_classic_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	dataName := "data.ncloud_auto_scaling_schedule.schedule"
	resourceName := "ncloud_auto_scaling_schedule.test-schedule"
	start := testAccNcloudAutoscalingScheduleValidStart(t)
	end := testAccNcloudAutoscalingScheduleValidEnd(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingScheduleClassicConfig(name, start, end),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "desired_capacity", resourceName, "desired_capacity"),
					resource.TestCheckResourceAttrPair(dataName, "min_size", resourceName, "min_size"),
					resource.TestCheckResourceAttrPair(dataName, "max_size", resourceName, "max_size"),
					resource.TestCheckResourceAttrPair(dataName, "start_time", resourceName, "start_time"),
					resource.TestCheckResourceAttrPair(dataName, "end_time", resourceName, "end_time"),
					resource.TestCheckResourceAttrPair(dataName, "recurrence", resourceName, "recurrence"),
					resource.TestCheckResourceAttrPair(dataName, "auto_scaling_group_no", resourceName, "auto_scaling_group_no"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAutoScalingSchedule_vpc_basic(t *testing.T) {
	t.Skip()
	{
		// Skip: deprecated server_image_product_code
	}

	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	dataName := "data.ncloud_auto_scaling_schedule.schedule"
	resourceName := "ncloud_auto_scaling_schedule.test-schedule"
	start := testAccNcloudAutoscalingScheduleValidStart(t)
	end := testAccNcloudAutoscalingScheduleValidEnd(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingScheduleVpcConfig(name, start, end),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "desired_capacity", resourceName, "desired_capacity"),
					resource.TestCheckResourceAttrPair(dataName, "min_size", resourceName, "min_size"),
					resource.TestCheckResourceAttrPair(dataName, "max_size", resourceName, "max_size"),
					resource.TestCheckResourceAttrPair(dataName, "start_time", resourceName, "start_time"),
					resource.TestCheckResourceAttrPair(dataName, "end_time", resourceName, "end_time"),
					resource.TestCheckResourceAttrPair(dataName, "recurrence", resourceName, "recurrence"),
					resource.TestCheckResourceAttrPair(dataName, "auto_scaling_group_no", resourceName, "auto_scaling_group_no"),
					resource.TestCheckResourceAttrPair(dataName, "time_zone", resourceName, "time_zone"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAutoScalingScheduleVpcConfig(name string, start string, end string) string {
	return testAccNcloudAutoScalingScheduleVpcConfig(name, start, end) + `
data "ncloud_auto_scaling_schedule" "schedule" {
	id = ncloud_auto_scaling_schedule.test-schedule.name
	auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	depends_on = [ncloud_auto_scaling_schedule.test-schedule]
}`
}

func testAccDataSourceNcloudAutoScalingScheduleClassicConfig(name string, start string, end string) string {
	return testAccNcloudAutoScalingScheduleClassicConfig(name, start, end) + `
data "ncloud_auto_scaling_schedule" "schedule" {
	id = ncloud_auto_scaling_schedule.test-schedule.name
	auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	depends_on = [ncloud_auto_scaling_schedule.test-schedule]
}`
}
