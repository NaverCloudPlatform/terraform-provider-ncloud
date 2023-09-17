package server_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

// ignore test : should use real access_control_group_configuration_no
func testAccDataSourceNcloudAccessControlRulesBasic(t *testing.T) {
	testId := os.Getenv("TEST_ID")
	if testId == "" {
		log.Println("[ERROR] ENV 'TEST_ID' is required")
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudAccessControlRulesConfig(testId),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_access_control_rules.test"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudAccessControlRulesConfig(testConfigNo string) string {
	return fmt.Sprintf(`
data "ncloud_access_control_rules" "test" {
	access_control_group_configuration_no = "%s"
}
`, testConfigNo)

}
