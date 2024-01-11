package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudRedisConfigGroup_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_redis_config_group.by_name"
	resourceName := "ncloud_redis_config_group.test"
	testConfigGroupName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	version := "7.0.13-simple"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRedisConfigGroupConfig(testConfigGroupName, version),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceRedisConfigGroupConfig(testConfigGroupName string, version string) string {
	return fmt.Sprintf(`
resource "ncloud_redis_config_group" "test" {
    name               = "%[1]s"
    redis_version      = "%[2]s"
    description        = "ACC TEST"
}

data "ncloud_redis_config_group" "by_name" {
    name = ncloud_redis_config_group.test.name
}
	`, testConfigGroupName, version)
}
