package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceNcloudLbTargetGroup_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-testacc-tg-%s", acctest.RandString(5))
	dataName := "data.ncloud_lb_target_group.tg_test"
	resourceName := "ncloud_lb_target_group.my-tg"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudLbTargetGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "protocol", resourceName, "protocol"),
					resource.TestCheckResourceAttrPair(dataName, "target_type", resourceName, "target_type"),
					resource.TestCheckResourceAttrPair(dataName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "health_check", resourceName, "health_check"),
					resource.TestCheckResourceAttrPair(dataName, "algorithm_type", resourceName, "algorithm_type"),
					resource.TestCheckResourceAttrPair(dataName, "use_sticky_session", resourceName, "use_sticky_session"),
				),
				SkipFunc: nil,
			},
		},
	})
}

func testAccDataSourceNcloudLbTargetGroupConfig(name string) string {
	return testAccResourceNcloudLbTargetGroupConfig(name) + fmt.Sprintf(`
data "ncloud_lb_target_group" "tg_test" {
	id = ncloud_lb_target_group.my-tg.target_group_no
}
`)
}
