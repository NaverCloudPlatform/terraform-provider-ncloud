package autoscaling_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudAutoScalingAdjustmentTypes_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_auto_scaling_adjustment_types.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingAdjustmentTypesVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAutoScalingAdjustmentTypes_vpc_byFilterCode(t *testing.T) {
	dataName := "data.ncloud_auto_scaling_adjustment_types.by_filter"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingAdjustmentTypesByFilterCodeConfig("EXACT"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

var testAccDataSourceNcloudAutoScalingAdjustmentTypesVpcConfig = `
	data "ncloud_auto_scaling_adjustment_types" "test" { }
`

func testAccDataSourceNcloudAutoScalingAdjustmentTypesByFilterCodeConfig(code string) string {
	return fmt.Sprintf(`
	data "ncloud_auto_scaling_adjustment_types" "by_filter" {
		filter {
			name   = "code"
			values = ["%s"]
		}
	}
`, code)
}
