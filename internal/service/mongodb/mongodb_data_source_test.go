package mongodb_test

import (
	"fmt"
	"testing"

	randacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudMongoDb_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_mongodb.by_id"
	resourceName := "ncloud_mongodb.mongodb"
	testMongoDbName := fmt.Sprintf("tf-mongodb-%s", randacctest.RandString(4))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongodbConfig(testMongoDbName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "service_name", resourceName, "service_name"),
					resource.TestCheckResourceAttrPair(dataName, "image_product_code", resourceName, "image_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "backup_time", resourceName, "backup_time"),
					resource.TestCheckResourceAttrPair(dataName, "backup_file_retention_period", resourceName, "backup_file_retention_period"),
				),
			},
		},
	})
}

func testAccDataSourceMongodbConfig(testMongoDbName string) string {
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
	vpc_no = 	ncloud_vpc.vpc.vpc_no
	subnet_no = ncloud_subnet.subnet.id
	service_name = "%[1]s"
    server_name_prefix = "ex-svr"
	user_name = "testuser"
	user_password = "t123456789!"
	cluster_type_code = "STAND_ALONE"
}

data "ncloud_mongodb" "by_id" {
	id = "${ncloud_mongodb.mongodb.id}"
}
`, testMongoDbName)
}
