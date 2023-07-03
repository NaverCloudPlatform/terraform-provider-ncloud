package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudLbTargetGroup_basic(t *testing.T) {
	name := fmt.Sprintf("terraform-testacc-tg-%s", acctest.RandString(5))
	dataName := "data.ncloud_lb_target_group.test"
	resourceName := "ncloud_lb_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudLbTargetGroupConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "protocol", resourceName, "protocol"),
					resource.TestCheckResourceAttrPair(dataName, "target_type", resourceName, "target_type"),
					resource.TestCheckResourceAttrPair(dataName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "health_check", resourceName, "health_check"),
					resource.TestCheckResourceAttrPair(dataName, "algorithm_type", resourceName, "algorithm_type"),
					resource.TestCheckResourceAttrPair(dataName, "use_sticky_session", resourceName, "use_sticky_session"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudLbTargetGroupConfig(name string) string {
	return testAccResourceNcloudLbTargetGroupConfig(name) + fmt.Sprintf(`
data "ncloud_lb_target_group" "test" {
	id = ncloud_lb_target_group.test.target_group_no
}
`)
}
