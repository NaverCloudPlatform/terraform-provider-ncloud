package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudRedis_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_redis.by_id"
	resourceName := "ncloud_redis.test"
	testRedisName := fmt.Sprintf("tf-redis-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRedisConfig(testRedisName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "service_name", resourceName, "service_name"),
				),
			},
		},
	})
}

func testAccDataSourceRedisConfig(testRedisName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test_vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "test_subnet" {
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
}

resource "ncloud_redis_config_group" "example" {
    name               = "%[1]s"
    redis_version      = "7.0.13-simple"
    description        = "ACC TEST"
}

resource "ncloud_redis" "test" {
    service_name       = "%[1]s"
    server_name_prefix = "ex-svr"
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
    subnet_no          = ncloud_subnet.test_subnet.id
    config_group_no    = ncloud_redis_config_group.example.id
	image_product_code  = "SW.VRDS.OS.LNX64.ROCKY.0810.REDIS.B050"
	engine_version_code = "7.0.13"
    mode = "SIMPLE"
}

data "ncloud_redis" "by_id" {
    id = ncloud_redis.test.id
}
	`, testRedisName)
}
