package redis_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	redisservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/redis"
)

func TestAccResourceNcloudRedisConfigGroup_vpc_basic(t *testing.T) {
	testConfigGroupName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resourceName := "ncloud_redis_config_group.test"
	version := "7.0.13-simple"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRedisConfigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRedisConfigGroupConfig(testConfigGroupName, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", testConfigGroupName),
					resource.TestCheckResourceAttr(resourceName, "redis_version", version),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     testConfigGroupName,
			},
		},
	})
}

func testAccCheckRedisConfigGroupDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_redis_config_group" {
			continue
		}
		instance, err := redisservice.GetRedisConfigGroup(context.Background(), config, rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("still exists")
		}
	}

	return nil
}

func testAccResourceRedisConfigGroupConfig(testConfigGroupName string, version string) string {
	return fmt.Sprintf(`
resource "ncloud_redis_config_group" "test" {
    name               = "%[1]s"
    redis_version      = "%[2]s"
    description        = "ACC TEST"
}

`, testConfigGroupName, version)
}
