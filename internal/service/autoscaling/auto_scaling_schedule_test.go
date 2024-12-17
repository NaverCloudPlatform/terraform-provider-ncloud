package autoscaling_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/autoscaling"
)

func TestAccResourceNcloudAutoScalingSchedule_classic_basic(t *testing.T) {
	var schedule autoscaling.AutoScalingSchedule
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceName := "ncloud_auto_scaling_schedule.test-schedule"
	start := testAccNcloudAutoscalingScheduleValidStart(t)
	end := testAccNcloudAutoscalingScheduleValidEnd(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingScheduleDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNcloudAutoScalingScheduleClassicConfig(name, start, end),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingScheduleExists(resourceName, &schedule, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceName, "desired_capacity", "1"),
				),
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingSchedule_vpc_basic(t *testing.T) {
	t.Skip()
	{
		// Skip: deprecated server_image_product_code
	}

	var schedule autoscaling.AutoScalingSchedule
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceName := "ncloud_auto_scaling_schedule.test-schedule"
	start := testAccNcloudAutoscalingScheduleValidStart(t)
	end := testAccNcloudAutoscalingScheduleValidEnd(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingScheduleDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNcloudAutoScalingScheduleVpcConfig(name, start, end),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingScheduleExists(resourceName, &schedule, GetTestProvider(true)),
				),
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingSchedule_classic_disappears(t *testing.T) {
	var schedule autoscaling.AutoScalingSchedule
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceName := "ncloud_auto_scaling_schedule.test-schedule"
	start := testAccNcloudAutoscalingScheduleValidStart(t)
	end := testAccNcloudAutoscalingScheduleValidEnd(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingScheduleDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNcloudAutoScalingScheduleClassicConfig(name, start, end),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingScheduleExists(resourceName, &schedule, GetTestProvider(false)),
					TestAccCheckResourceDisappears(GetTestProvider(false), autoscaling.ResourceNcloudAutoScalingSchedule(), resourceName),
					resource.TestCheckResourceAttr(resourceName, "desired_capacity", "1"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingSchedule_vpc_disappears(t *testing.T) {
	t.Skip()
	{
		// Skip: deprecated server_image_product_code
	}

	var schedule autoscaling.AutoScalingSchedule
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceName := "ncloud_auto_scaling_schedule.test-schedule"
	start := testAccNcloudAutoscalingScheduleValidStart(t)
	end := testAccNcloudAutoscalingScheduleValidEnd(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingScheduleDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNcloudAutoScalingScheduleVpcConfig(name, start, end),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingScheduleExists(resourceName, &schedule, GetTestProvider(true)),
					TestAccCheckResourceDisappears(GetTestProvider(true), autoscaling.ResourceNcloudAutoScalingSchedule(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccNcloudAutoscalingScheduleValidEnd(t *testing.T) string {
	return testAccNcloudAutoscalingScheduleTime(t, "2h")
}

func testAccNcloudAutoscalingScheduleValidStart(t *testing.T) string {
	return testAccNcloudAutoscalingScheduleTime(t, "1h")
}

func testAccNcloudAutoscalingScheduleTime(t *testing.T, duration string) string {
	// TODO: Multi region
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		t.Fatalf("err parsing location: %s", err)
	}
	n := time.Now().In(loc)
	d, err := time.ParseDuration(duration)
	if err != nil {
		t.Fatalf("err parsing time duration: %s", err)
	}
	return n.Add(d).Format(autoscaling.SCHEDULE_TIME_FORMAT)
}

func testAccCheckNcloudAutoScalingScheduleDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_auto_scaling_schedule" {
			continue
		}
		autoScalingSchedule, err := autoscaling.GetAutoScalingSchedule(config, rs.Primary.ID, rs.Primary.Attributes["auto_scaling_group_no"])
		if err != nil {
			return err
		}

		if autoScalingSchedule != nil {
			return fmt.Errorf("AutoScalingSchedule(%s) still exists", ncloud.StringValue(autoScalingSchedule.ScheduledActionName))
		}
	}
	return nil
}

func testAccCheckNcloudAutoScalingScheduleExists(n string, schedule *autoscaling.AutoScalingSchedule, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AutoScalingSchdule ID is set: %s", n)
		}

		config := provider.Meta().(*conn.ProviderConfig)
		autoScalingSchedule, err := autoscaling.GetAutoScalingSchedule(config, rs.Primary.ID, rs.Primary.Attributes["auto_scaling_group_no"])
		if err != nil {
			return err
		}
		if autoScalingSchedule == nil {
			return fmt.Errorf("Not found AutoScalingSchedule : %s", rs.Primary.ID)
		}
		*schedule = *autoScalingSchedule
		return nil
	}
}

func testAccNcloudAutoScalingScheduleClassicConfigBase(name string) string {
	return fmt.Sprintf(`
resource "ncloud_launch_configuration" "test" {
    name = "%[1]s"
    server_image_product_code = "SPSW0LINUX000046"
    server_product_code = "SPSVRSSD00000003"
}

resource "ncloud_auto_scaling_group" "test" {
	name = "%[1]s"
	launch_configuration_no = ncloud_launch_configuration.test.launch_configuration_no
	min_size = 1
	max_size = 1
	zone_no_list = ["2"]
}
`, name)
}

func testAccNcloudAutoScalingScheduleClassicConfig(name string, start string, end string) string {
	return testAccNcloudAutoScalingScheduleClassicConfigBase(name) + fmt.Sprintf(`
resource "ncloud_auto_scaling_schedule" "test-schedule" {
	name = "%[1]s"
	min_size = 1
	max_size = 1
	desired_capacity = 1
	start_time = "%[2]s"
	end_time = "%[3]s"
	auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}`, name, start, end)
}

func testAccNcloudAutoScalingScheduleVpcConfigBase(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_launch_configuration" "test" {
    name = "%[1]s"
    server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
    server_product_code = "SVR.VSVR.HICPU.C002.M004.NET.SSD.B050.G002"
}

resource "ncloud_auto_scaling_group" "test" {
	name = "%[1]s"
	access_control_group_no_list = [ncloud_vpc.test.default_access_control_group_no]
	subnet_no = ncloud_subnet.test.subnet_no
	launch_configuration_no = ncloud_launch_configuration.test.launch_configuration_no
	min_size = 1
	max_size = 1
}
`, name)
}

func testAccNcloudAutoScalingScheduleVpcConfig(name string, start string, end string) string {
	return testAccNcloudAutoScalingScheduleVpcConfigBase(name) + fmt.Sprintf(`
resource "ncloud_auto_scaling_schedule" "test-schedule" {
	name = "%[1]s"
	min_size = 1
	max_size = 1
	desired_capacity = 1
	start_time = "%[2]s"
	end_time = "%[3]s"
	auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}`, name, start, end)
}
