package hadoop_test

import (
	"fmt"

	randacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoop_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_hadoop.hadoop"
	resourceName := "ncloud_hadoop.hadoop"
	instanceName := fmt.Sprintf("tf-hadoop-%s", randacctest.RandString(3))
	bucketName := "bbb"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHadoopConfig(instanceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "cluster_name", resourceName, "cluster_name"),
					resource.TestCheckResourceAttrPair(dataName, "cluster_type_code", resourceName, "cluster_type_code"),
					resource.TestCheckResourceAttrPair(dataName, "image_product_code", resourceName, "image_product_code"),
					resource.TestCheckResourceAttrPair(dataName, "hadoop_server_list.1.server_name", resourceName, "hadoop_server_list.1.server_name"),
				),
			},
		},
	})
}

func testAccDataSourceHadoopConfig(name, bucketName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "edge_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-edge"
	subnet             = "10.5.0.0/18"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "master_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-master"
	subnet             = "10.5.64.0/19"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "worker_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-worker"
	subnet             = "10.5.96.0/20"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_hadoop" "hadoop" {
	vpc_no = ncloud_vpc.test.vpc_no
	cluster_name = "%[1]s"
	cluster_type_code = "CORE_HADOOP_WITH_SPARK"
	admin_user_name = "admin-test"
	admin_user_password = "Admin!2Admin"
	login_key_name = "naverCloud"
	master_node_subnet_no = ncloud_subnet.master_subnet.subnet_no
	edge_node_subnet_no = ncloud_subnet.edge_subnet.subnet_no
	worker_node_subnet_no = ncloud_subnet.worker_subnet.subnet_no
	bucket_name = "%[2]s"
	master_node_data_storage_type = "SSD"
	worker_node_data_storage_type = "SSD"
	master_node_data_storage_size = 100
	worker_node_data_storage_size = 100
}

data "ncloud_hadoop" "hadoop" {
	cluster_name = ncloud_hadoop.hadoop.cluster_name
}
`, name, bucketName)
}
