package autoscaling_test

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/autoscaling"
)

func TestAccResourceNcloudAutoScalingPolicy_classic_basic(t *testing.T) {
	// Images are all deprecated in Classic
	t.Skip()

	var policy autoscaling.AutoScalingPolicy
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceCHANG := "ncloud_auto_scaling_policy.test-policy-CHANG"
	resourceEXACT := "ncloud_auto_scaling_policy.test-policy-EXACT"
	resourcePRCNT := "ncloud_auto_scaling_policy.test-policy-PRCNT"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingPolicyDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNcloudAutoScalingPolicyClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingPolicyExists(resourceCHANG, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceCHANG, "adjustment_type_code", "CHANG"),
					resource.TestCheckResourceAttr(resourceCHANG, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceCHANG, "name", name+"-chang"),
					resource.TestCheckResourceAttr(resourceCHANG, "cooldown", "300"),

					testAccCheckNcloudAutoScalingPolicyExists(resourceEXACT, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceEXACT, "adjustment_type_code", "EXACT"),
					resource.TestCheckResourceAttr(resourceEXACT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceEXACT, "name", name+"-exact"),
					resource.TestCheckResourceAttr(resourceEXACT, "cooldown", "300"),

					testAccCheckNcloudAutoScalingPolicyExists(resourcePRCNT, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourcePRCNT, "adjustment_type_code", "PRCNT"),
					resource.TestCheckResourceAttr(resourcePRCNT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourcePRCNT, "name", name+"-prcnt"),
					resource.TestCheckResourceAttr(resourcePRCNT, "cooldown", "300"),
				),
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingPolicy_classic_zero_value(t *testing.T) {
	// Images are all deprecated in Classic
	t.Skip()

	var policy autoscaling.AutoScalingPolicy
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceCHANG := "ncloud_auto_scaling_policy.test-policy-CHANG"
	resourceEXACT := "ncloud_auto_scaling_policy.test-policy-EXACT"
	resourcePRCNT := "ncloud_auto_scaling_policy.test-policy-PRCNT"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingPolicyDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				//default
				Config: testAccNcloudAutoScalingPolicyClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingPolicyExists(resourceCHANG, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceCHANG, "adjustment_type_code", "CHANG"),
					resource.TestCheckResourceAttr(resourceCHANG, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceCHANG, "name", name+"-chang"),
					resource.TestCheckResourceAttr(resourceCHANG, "cooldown", "300"),

					testAccCheckNcloudAutoScalingPolicyExists(resourceEXACT, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceEXACT, "adjustment_type_code", "EXACT"),
					resource.TestCheckResourceAttr(resourceEXACT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceEXACT, "name", name+"-exact"),
					resource.TestCheckResourceAttr(resourceEXACT, "cooldown", "300"),

					testAccCheckNcloudAutoScalingPolicyExists(resourcePRCNT, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourcePRCNT, "adjustment_type_code", "PRCNT"),
					resource.TestCheckResourceAttr(resourcePRCNT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourcePRCNT, "name", name+"-prcnt"),
					resource.TestCheckResourceAttr(resourcePRCNT, "cooldown", "300"),
				),
			},
			{
				//zero-value
				Config: testAccNcloudAutoScalingPolicyClassicConfigWhenSetZero(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingPolicyExists(resourceCHANG, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceCHANG, "adjustment_type_code", "CHANG"),
					resource.TestCheckResourceAttr(resourceCHANG, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceCHANG, "cooldown", "0"),

					testAccCheckNcloudAutoScalingPolicyExists(resourceEXACT, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceEXACT, "adjustment_type_code", "EXACT"),
					resource.TestCheckResourceAttr(resourceEXACT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceEXACT, "name", name+"-exact"),
					resource.TestCheckResourceAttr(resourceEXACT, "cooldown", "0"),

					testAccCheckNcloudAutoScalingPolicyExists(resourcePRCNT, &policy, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourcePRCNT, "adjustment_type_code", "PRCNT"),
					resource.TestCheckResourceAttr(resourcePRCNT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourcePRCNT, "name", name+"-prcnt"),
					resource.TestCheckResourceAttr(resourcePRCNT, "cooldown", "0"),
				),
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingPolicy_vpc_zero_value(t *testing.T) {
	t.Skip()
	{
		// Skip: deprecated server_image_product_code
	}

	var policy autoscaling.AutoScalingPolicy
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceCHANG := "ncloud_auto_scaling_policy.test-policy-CHANG"
	resourceEXACT := "ncloud_auto_scaling_policy.test-policy-EXACT"
	resourcePRCNT := "ncloud_auto_scaling_policy.test-policy-PRCNT"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingPolicyDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				//default
				Config: testAccNcloudAutoScalingPolicyVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingPolicyExists(resourceCHANG, &policy, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceCHANG, "adjustment_type_code", "CHANG"),
					resource.TestCheckResourceAttr(resourceCHANG, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceCHANG, "name", name+"-chang"),
					resource.TestCheckResourceAttr(resourceCHANG, "cooldown", "300"),

					testAccCheckNcloudAutoScalingPolicyExists(resourceEXACT, &policy, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceEXACT, "adjustment_type_code", "EXACT"),
					resource.TestCheckResourceAttr(resourceEXACT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceEXACT, "name", name+"-exact"),
					resource.TestCheckResourceAttr(resourceEXACT, "cooldown", "300"),

					testAccCheckNcloudAutoScalingPolicyExists(resourcePRCNT, &policy, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourcePRCNT, "adjustment_type_code", "PRCNT"),
					resource.TestCheckResourceAttr(resourcePRCNT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourcePRCNT, "name", name+"-prcnt"),
					resource.TestCheckResourceAttr(resourcePRCNT, "cooldown", "300"),
				),
			},
			{
				//zero-value
				Config: testAccNcloudAutoScalingPolicyVpcConfigWhenSetZero(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingPolicyExists(resourceCHANG, &policy, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceCHANG, "adjustment_type_code", "CHANG"),
					resource.TestCheckResourceAttr(resourceCHANG, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceCHANG, "name", name+"-chang"),
					resource.TestCheckResourceAttr(resourceCHANG, "cooldown", "0"),

					testAccCheckNcloudAutoScalingPolicyExists(resourceEXACT, &policy, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceEXACT, "adjustment_type_code", "EXACT"),
					resource.TestCheckResourceAttr(resourceEXACT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceEXACT, "name", name+"-exact"),
					resource.TestCheckResourceAttr(resourceEXACT, "cooldown", "0"),

					testAccCheckNcloudAutoScalingPolicyExists(resourcePRCNT, &policy, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourcePRCNT, "adjustment_type_code", "PRCNT"),
					resource.TestCheckResourceAttr(resourcePRCNT, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourcePRCNT, "name", name+"-prcnt"),
					resource.TestCheckResourceAttr(resourcePRCNT, "cooldown", "0"),
				),
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingPolicy_classic_disappears(t *testing.T) {
	// Images are all deprecated in Classic
	t.Skip()

	var policy autoscaling.AutoScalingPolicy
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceCHANG := "ncloud_auto_scaling_policy.test-policy-CHANG"
	resourceEXACT := "ncloud_auto_scaling_policy.test-policy-EXACT"
	resourcePRCNT := "ncloud_auto_scaling_policy.test-policy-PRCNT"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingPolicyDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNcloudAutoScalingPolicyClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingPolicyExists(resourceCHANG, &policy, GetTestProvider(false)),
					testAccCheckNcloudAutoScalingPolicyExists(resourceEXACT, &policy, GetTestProvider(false)),
					testAccCheckNcloudAutoScalingPolicyExists(resourcePRCNT, &policy, GetTestProvider(false)),

					TestAccCheckResourceDisappears(GetTestProvider(false), autoscaling.ResourceNcloudAutoScalingPolicy(), resourceCHANG),
					TestAccCheckResourceDisappears(GetTestProvider(false), autoscaling.ResourceNcloudAutoScalingPolicy(), resourceEXACT),
					TestAccCheckResourceDisappears(GetTestProvider(false), autoscaling.ResourceNcloudAutoScalingPolicy(), resourcePRCNT),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudAutoScalingPolicy_vpc_disappears(t *testing.T) {
	t.Skip()
	{
		// Skip: deprecated server_image_product_code
	}

	var policy autoscaling.AutoScalingPolicy
	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))
	resourceCHANG := "ncloud_auto_scaling_policy.test-policy-CHANG"
	resourceEXACT := "ncloud_auto_scaling_policy.test-policy-EXACT"
	resourcePRCNT := "ncloud_auto_scaling_policy.test-policy-PRCNT"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNcloudAutoScalingPolicyDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNcloudAutoScalingPolicyVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNcloudAutoScalingPolicyExists(resourceCHANG, &policy, GetTestProvider(true)),
					testAccCheckNcloudAutoScalingPolicyExists(resourceEXACT, &policy, GetTestProvider(true)),
					testAccCheckNcloudAutoScalingPolicyExists(resourcePRCNT, &policy, GetTestProvider(true)),

					TestAccCheckResourceDisappears(GetTestProvider(true), autoscaling.ResourceNcloudAutoScalingPolicy(), resourceCHANG),
					TestAccCheckResourceDisappears(GetTestProvider(true), autoscaling.ResourceNcloudAutoScalingPolicy(), resourceEXACT),
					TestAccCheckResourceDisappears(GetTestProvider(true), autoscaling.ResourceNcloudAutoScalingPolicy(), resourcePRCNT),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckNcloudAutoScalingPolicyExists(n string, p *autoscaling.AutoScalingPolicy, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AutoScalingPolicy ID is set: %s", n)
		}

		config := provider.Meta().(*conn.ProviderConfig)
		autoScalingPolicy, err := autoscaling.GetAutoScalingPolicy(config, rs.Primary.ID, rs.Primary.Attributes["auto_scaling_group_no"])
		if err != nil {
			return err
		}
		if autoScalingPolicy == nil {
			return fmt.Errorf("Not found AutoScalingPolicy : %s", rs.Primary.ID)
		}
		*p = *autoScalingPolicy
		return nil
	}
}

func testAccCheckNcloudAutoScalingPolicyDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_auto_scaling_policy" {
			continue
		}
		autoScalingPolicy, err := autoscaling.GetAutoScalingPolicy(config, rs.Primary.ID, rs.Primary.Attributes["auto_scaling_group_no"])
		if err != nil {
			return err
		}

		if autoScalingPolicy != nil {
			return fmt.Errorf("AutoScalingPolicy(%s) still exists", ncloud.StringValue(autoScalingPolicy.AutoScalingPolicyName))
		}
	}
	return nil
}

func testAccNcloudAutoScalingPolicyVpcConfigBase(name string) string {
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

func testAccNcloudAutoScalingPolicyVpcConfig(name string) string {
	return testAccNcloudAutoScalingPolicyVpcConfigBase(name) + fmt.Sprintf(`
resource "ncloud_auto_scaling_policy" "test-policy-CHANG" {
    name = "%[1]s-chang"
    adjustment_type_code = "CHANG"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}

resource "ncloud_auto_scaling_policy" "test-policy-EXACT" {
    name = "%[1]s-exact"
    adjustment_type_code = "EXACT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}

resource "ncloud_auto_scaling_policy" "test-policy-PRCNT" {
    name = "%[1]s-prcnt"
    adjustment_type_code = "PRCNT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}
`, name)
}

func testAccNcloudAutoScalingPolicyVpcConfigWhenSetZero(name string) string {
	return testAccNcloudAutoScalingPolicyVpcConfigBase(name) + fmt.Sprintf(`
resource "ncloud_auto_scaling_policy" "test-policy-CHANG" {
    name = "%[1]s-chang"
    adjustment_type_code = "CHANG"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	cooldown = 0
}

resource "ncloud_auto_scaling_policy" "test-policy-EXACT" {
    name = "%[1]s-exact"
    adjustment_type_code = "EXACT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	cooldown = 0
}

resource "ncloud_auto_scaling_policy" "test-policy-PRCNT" {
    name = "%[1]s-prcnt"
    adjustment_type_code = "PRCNT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	cooldown = 0
}
`, name)
}

func testAccNcloudAutoScalingPolicyClassicConfigBase(name string) string {
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

func testAccNcloudAutoScalingPolicyClassicConfig(name string) string {
	return testAccNcloudAutoScalingPolicyClassicConfigBase(name) + fmt.Sprintf(`
resource "ncloud_auto_scaling_policy" "test-policy-CHANG" {
    name = "%[1]s-chang"
    adjustment_type_code = "CHANG"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}

resource "ncloud_auto_scaling_policy" "test-policy-EXACT" {
    name = "%[1]s-exact"
    adjustment_type_code = "EXACT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}

resource "ncloud_auto_scaling_policy" "test-policy-PRCNT" {
    name = "%[1]s-prcnt"
    adjustment_type_code = "PRCNT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
}
`, name)
}

func testAccNcloudAutoScalingPolicyClassicConfigWhenSetZero(name string) string {
	return testAccNcloudAutoScalingPolicyClassicConfigBase(name) + fmt.Sprintf(`
resource "ncloud_auto_scaling_policy" "test-policy-CHANG" {
    name = "%[1]s-chang"
    adjustment_type_code = "CHANG"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	cooldown = 0
}

resource "ncloud_auto_scaling_policy" "test-policy-EXACT" {
    name = "%[1]s-exact"
    adjustment_type_code = "EXACT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	cooldown = 0
}

resource "ncloud_auto_scaling_policy" "test-policy-PRCNT" {
    name = "%[1]s-prcnt"
    adjustment_type_code = "PRCNT"
    scaling_adjustment = 2
    auto_scaling_group_no = ncloud_auto_scaling_group.test.auto_scaling_group_no
	cooldown = 0
}
`, name)
}
