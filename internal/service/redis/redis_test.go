package redis_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vredis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	redisservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/redis"
)

func TestAccResourceNcloudRedis_vpc_basic(t *testing.T) {
	var redisInstance vredis.CloudRedisInstance
	testRedisName := fmt.Sprintf("tf-redis-%s", acctest.RandString(5))
	resourceName := "ncloud_redis.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRedisDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRedisConfig(testRedisName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedisExistsWithProvider(resourceName, &redisInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "service_name", testRedisName),
					resource.TestCheckResourceAttr(resourceName, "server_name_prefix", "ex-svr"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine_version_code"},
			},
		},
	})
}

func testAccCheckRedisExistsWithProvider(n string, redis *vredis.CloudRedisInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		redisInstance, err := redisservice.GetRedisDetail(context.Background(), config, resource.Primary.ID)
		if err != nil {
			return err
		}

		if redisInstance != nil {
			*redis = *redisInstance
			return nil
		}

		return fmt.Errorf("redis instance not found")
	}
}

func testAccCheckRedisDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_redis" {
			continue
		}
		instance, err := redisservice.GetRedisDetail(context.Background(), config, rs.Primary.ID)
		if err != nil && !checkNoInstanceResponse(err) {
			return err
		}

		if instance != nil {
			return errors.New("redis still exists")
		}
	}

	return nil
}

func checkNoInstanceResponse(err error) bool {
	return strings.Contains(err.Error(), "5001017")
}

func testAccResourceRedisConfig(testRedisName string) string {
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
    description        = "example"
}

resource "ncloud_redis" "test" {
    service_name        = "%[1]s"
    server_name_prefix  = "ex-svr"
	vpc_no              = ncloud_vpc.test_vpc.vpc_no
    subnet_no           = ncloud_subnet.test_subnet.id
    config_group_no     = ncloud_redis_config_group.example.id
	image_product_code  = "SW.VRDS.OS.LNX64.ROCKY.0810.REDIS.B050"
	engine_version_code = "7.0.13"
    mode = "SIMPLE"
}
`, testRedisName)
}
