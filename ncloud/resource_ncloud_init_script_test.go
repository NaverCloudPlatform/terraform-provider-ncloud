package ncloud

import (
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudInitScript_basic(t *testing.T) {
	var InitScript vserver.InitScript
	name := fmt.Sprintf("tf-init-script-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_init_script.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInitScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudInitScriptConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInitScriptExists(resourceName, &InitScript),
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

func TestAccResourceNcloudInitScript_disappears(t *testing.T) {
	var InitScript vserver.InitScript
	name := fmt.Sprintf("tf-init-script-disappear-%s", acctest.RandString(5))
	resourceName := "ncloud_init_script.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInitScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudInitScriptConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInitScriptExists(resourceName, &InitScript),
					testAccCheckInitScriptDisappears(&InitScript),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNcloudInitScriptConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_init_script" "foo" {
	name    = "%s"
	content = "#!/usr/bin/env\nls -al"
}
`, name)
}

func testAccCheckInitScriptExists(n string, InitScript *vserver.InitScript) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Init script id is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getInitScript(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*InitScript = *instance

		return nil
	}
}

func testAccCheckInitScriptDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_init_script" {
			continue
		}

		instance, err := getInitScript(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("init script still exists")
		}
	}

	return nil
}

func testAccCheckInitScriptDisappears(instance *vserver.InitScript) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)
		return deleteInitScript(config, *instance.InitScriptNo)
	}
}
