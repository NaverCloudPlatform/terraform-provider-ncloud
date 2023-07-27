package autoscaling_test

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/autoscaling"
)

func TestAccResourceNcloudLaunchConfiguration_classic_basic(t *testing.T) {
	var launchConfiguration autoscaling.LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SPSW0LINUX000046"
	serverProductCode := "SPSVRSSD00000003"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode, serverProductCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, GetTestProvider(false)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudLaunchConfiguration_vpc_basic(t *testing.T) {
	var launchConfiguration autoscaling.LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, GetTestProvider(true)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudLaunchConfiguration_classic_disappears(t *testing.T) {
	var launchConfiguration autoscaling.LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SPSW0LINUX000046"
	serverProductCode := "SPSVRSSD00000003"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(false),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode, serverProductCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, GetTestProvider(false)),
					TestAccCheckResourceDisappears(GetTestProvider(false), autoscaling.ResourceNcloudLaunchConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudLaunchConfiguration_vpc_disappears(t *testing.T) {
	var launchConfiguration autoscaling.LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, GetTestProvider(true)),
					TestAccCheckResourceDisappears(GetTestProvider(true), autoscaling.ResourceNcloudLaunchConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckLaunchConfigurationExists(n string, l *autoscaling.LaunchConfiguration, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LaunchConfiguration ID is set: %s", n)
		}

		config := provider.Meta().(*conn.ProviderConfig)
		launchConfiguration, err := autoscaling.GetLaunchConfiguration(config, rs.Primary.ID)
		if err != nil {
			return err
		}
		if launchConfiguration == nil {
			return fmt.Errorf("Not found LaunchConfiguration : %s", rs.Primary.ID)
		}
		*l = *launchConfiguration
		return nil
	}
}

func testAccCheckLaunchConfigurationDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_launch_configuration" {
			continue
		}
		launchConfiguration, err := autoscaling.GetClassicLaunchConfiguration(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		if launchConfiguration != nil {
			return fmt.Errorf("LaunchConfiguration(%s) still exists", ncloud.StringValue(launchConfiguration.LaunchConfigurationNo))
		}
	}
	return nil
}

func testAccLaunchConfigurationConfig(serverImageProductCode string, serverProductCode string) string {
	return fmt.Sprintf(`
resource "ncloud_launch_configuration" "lc" {
	server_image_product_code = "%[1]s"
	server_product_code = "%[2]s"
}
`, serverImageProductCode, serverProductCode)
}
