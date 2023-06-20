package autoscaling

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudAutoScalingPolicy_classic_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	dataName := "data.ncloud_auto_scaling_policy.policy"
	resourceName := "ncloud_auto_scaling_policy.test-policy-CHANG"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingPolicyClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "adjustment_type_code", resourceName, "adjustment_type_code"),
					resource.TestCheckResourceAttrPair(dataName, "scaling_adjustment", resourceName, "scaling_adjustment"),
					resource.TestCheckResourceAttrPair(dataName, "cooldown", resourceName, "cooldown"),
					resource.TestCheckResourceAttrPair(dataName, "min_adjustment_step", resourceName, "min_adjustment_step"),
					resource.TestCheckResourceAttrPair(dataName, "auto_scaling_group_no", resourceName, "auto_scaling_group_no"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudAutoScalingPolicy_vpc_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	dataName := "data.ncloud_auto_scaling_policy.policy"
	resourceName := "ncloud_auto_scaling_policy.test-policy-CHANG"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAutoScalingPolicyVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "adjustment_type_code", resourceName, "adjustment_type_code"),
					resource.TestCheckResourceAttrPair(dataName, "scaling_adjustment", resourceName, "scaling_adjustment"),
					resource.TestCheckResourceAttrPair(dataName, "cooldown", resourceName, "cooldown"),
					resource.TestCheckResourceAttrPair(dataName, "min_adjustment_step", resourceName, "min_adjustment_step"),
					resource.TestCheckResourceAttrPair(dataName, "auto_scaling_group_no", resourceName, "auto_scaling_group_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAutoScalingPolicyClassicConfig(name string) string {
	return testAccNcloudAutoScalingPolicyClassicConfig(name) + fmt.Sprintf(`
data "ncloud_auto_scaling_policy" "policy" {
	id = ncloud_auto_scaling_policy.test-policy-CHANG.name
	auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	depends_on = [ncloud_auto_scaling_policy.test-policy-CHANG]
}
`)
}

func testAccDataSourceNcloudAutoScalingPolicyVpcConfig(name string) string {
	return testAccNcloudAutoScalingPolicyVpcConfig(name) + fmt.Sprintf(`
data "ncloud_auto_scaling_policy" "policy" {
	id = ncloud_auto_scaling_policy.test-policy-CHANG.name
	auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	depends_on = [ncloud_auto_scaling_policy.test-policy-CHANG]
}
`)
}
