package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccResourceNcloudLaunchConfiguration_classic_basic(t *testing.T) {
	var launchConfiguration LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SPSW0LINUX000046"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, testAccClassicProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, testAccClassicProvider),
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
	var launchConfiguration LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, testAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, testAccProvider),
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
	var launchConfiguration LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SPSW0LINUX000046"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, testAccClassicProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, testAccClassicProvider),
					testAccCheckResourceDisappears(testAccClassicProvider, resourceNcloudLaunchConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudLaunchConfiguration_vpc_disappears(t *testing.T) {
	var launchConfiguration LaunchConfiguration
	resourceName := "ncloud_launch_configuration.lc"
	serverImageProductCode := "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckLaunchConfigurationDestroy(state, testAccProvider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLaunchConfigurationConfig(serverImageProductCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLaunchConfigurationExists(resourceName, &launchConfiguration, testAccProvider),
					testAccCheckResourceDisappears(testAccProvider, resourceNcloudLaunchConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckLaunchConfigurationExists(n string, l *LaunchConfiguration, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LaunchConfiguration ID is set: %s", n)
		}

		config := provider.Meta().(*ProviderConfig)
		launchConfiguration, err := getLaunchConfiguration(config, rs.Primary.ID)
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
	config := provider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_launch_configuration" {
			continue
		}
		launchConfiguration, err := getClassicLaunchConfiguration(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		if launchConfiguration != nil {
			return fmt.Errorf("LaunchConfiguration(%s) still exists", ncloud.StringValue(launchConfiguration.LaunchConfigurationNo))
		}
	}
	return nil
}

func testAccLaunchConfigurationConfig(serverImageProductCode string) string {
	return fmt.Sprintf(`
resource "ncloud_launch_configuration" "lc" {
	server_image_product_code = "%[1]s"
}
`, serverImageProductCode)
}
