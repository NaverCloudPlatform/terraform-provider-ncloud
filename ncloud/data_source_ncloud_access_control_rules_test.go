package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
	"testing"
)

// ignore test : should use real access_control_group_configuration_no
func testAccDataSourceNcloudAccessControlRulesBasic(t *testing.T) {
	t.Parallel()

	testId := os.Getenv("TEST_ID")
	if testId == "" {
		log.Println("[ERROR] ENV 'TEST_ID' is required")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlRulesConfig(testId),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceNcloudAccessControlRulesDataSourceID("data.ncloud_access_control_rules.test"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAccessControlRulesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("source ID not set")
		}
		return nil
	}
}

func testAccDataSourceNcloudAccessControlRulesConfig(testConfigNo string) string {
	return fmt.Sprintf(`
data "ncloud_access_control_rules" "test" {
	"access_control_group_configuration_no" = "%s"
}
`, testConfigNo)

}
