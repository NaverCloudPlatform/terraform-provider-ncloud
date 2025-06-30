package mongodb_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vmongodb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	mongodbservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/mongodb"
)

func TestAccResourceNcloudMongoDb_vpc_basic(t *testing.T) {
	var mongodbInstance vmongodb.CloudMongoDbInstance
	name := fmt.Sprintf("tf-mongodb-%s", sdkacctest.RandString(4))
	resourceName := "ncloud_mongodb.mongodb"
	clusterTypeCode := "STAND_ALONE"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDbVpcConfig(name, clusterTypeCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDbExists(resourceName, &mongodbInstance, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "service_name", name),
					resource.TestCheckResourceAttr(resourceName, "user_name", "testuser"),
					resource.TestCheckResourceAttr(resourceName, "user_password", "t123456789!"),
					resource.TestCheckResourceAttr(resourceName, "backup_time", "02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_file_retention_period", "1"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type_code", "STAND_ALONE"),
				),
			},
		},
	})
}

func TestAccResourceNcloudMongoDb_vpc_sharding(t *testing.T) {
	var mongodbInstance vmongodb.CloudMongoDbInstance
	name := fmt.Sprintf("tf-mongodb-%s", sdkacctest.RandString(4))
	resourceName := "ncloud_mongodb.mongodb"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDbVpcConfigShard(name, 2, 3, 0, 2, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDbExists(resourceName, &mongodbInstance, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "shard_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "member_server_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "arbiter_server_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "mongos_server_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "config_server_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type_code", "SHARDED_CLUSTER"),
				),
			},
		},
	})
}

func testAccCheckMongoDbExists(n string, mongodb *vmongodb.CloudMongoDbInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		mongodbInstance, err := mongodbservice.GetCloudMongoDbInstance(context.Background(), config, resource.Primary.ID)
		if err != nil {
			return err
		}

		if mongodbInstance != nil {
			*mongodb = *mongodbInstance
			return nil
		}

		return fmt.Errorf("mongodb instance not found")
	}
}

func testAccCheckMongoDbDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_mongodb" {
			continue
		}
		instance, err := mongodbservice.GetCloudMongoDbInstance(context.Background(), config, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("mongodb still exists")
		}
	}

	return nil
}

func testAccMongoDbVpcConfig(testMongoDbName string, clusterTypeCode string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}

resource "ncloud_mongodb" "mongodb" {
	vpc_no = ncloud_vpc.vpc.vpc_no
	subnet_no = ncloud_subnet.subnet.id
	service_name = "%[1]s"
    server_name_prefix = "ex-svr"
	cluster_type_code = "%[2]s"
	user_name = "testuser"
	user_password = "t123456789!"
}
`, testMongoDbName, clusterTypeCode)
}

func testAccMongoDbVpcConfigShard(testMongoDbName string, shardCount int, memberServerCount int, arbiterServerCount int, mongosServerCount int, configServerCount int) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.0.0.0/16"
}
	
resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}

resource "ncloud_mongodb" "mongodb" {
	vpc_no = ncloud_vpc.vpc.vpc_no
	subnet_no = ncloud_subnet.subnet.id
	service_name = "%[1]s"
    server_name_prefix = "ex-svr"
	cluster_type_code = "SHARDED_CLUSTER"
	user_name = "testuser"
	user_password = "t123456789!"

	shard_count = "%[2]v"
	member_server_count = "%[3]v"
	arbiter_server_count = "%[4]v"
	mongos_server_count = "%[5]v"
	config_server_count = "%[6]v"
}
`, testMongoDbName, shardCount, memberServerCount, arbiterServerCount, mongosServerCount, configServerCount)
}
