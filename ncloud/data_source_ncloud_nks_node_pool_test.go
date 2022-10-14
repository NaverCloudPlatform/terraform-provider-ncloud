package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNKSNodePool(t *testing.T) {
	dataName := "data.ncloud_nks_node_pool.node_pool"
	resourceName := "ncloud_nks_node_pool.node_pool"
	testClusterName := getTestClusterName()

	region, clusterType, productType, k8sVersion := getRegionAndNKSType()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSNodePoolConfig(testClusterName, clusterType, TF_TEST_NKS_LOGIN_KEY, k8sVersion, region, productType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "cluster_uuid", resourceName, "cluster_uuid"),
					resource.TestCheckResourceAttrPair(dataName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "k8s_version", resourceName, "k8s_version"),
					resource.TestCheckResourceAttrPair(dataName, "node_pool_name", resourceName, "node_pool_name"),
					resource.TestCheckResourceAttrPair(dataName, "node_count", resourceName, "node_count"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list", resourceName, "subnet_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "product_code", resourceName, "product_code"),
					resource.TestCheckResourceAttrPair(dataName, "autoscale.0.enabled", resourceName, "autoscale.0.enabled"),
					resource.TestCheckResourceAttrPair(dataName, "autoscale.0.min", resourceName, "autoscale.0.min"),
					resource.TestCheckResourceAttrPair(dataName, "autoscale.0.max", resourceName, "autoscale.0.max"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.name", resourceName, "nodes.0.name"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.instance_no", resourceName, "nodes.0.instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.spec", resourceName, "nodes.0.spec"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.private_ip", resourceName, "nodes.0.private_ip"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.public_ip", resourceName, "nodes.0.public_ip"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.node_status", resourceName, "nodes.0.node_status"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.container_version", resourceName, "nodes.0.container_version"),
					resource.TestCheckResourceAttrPair(dataName, "nodes.0.kernel_version", resourceName, "nodes.0.kernel_version"),
				),
			},
		},
	})
}

func testAccDataSourceNKSNodePoolConfig(testClusterName string, clusterType string, loginKey string, version string, region string, productType string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "subnet1" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-1"
	subnet             = "10.2.1.0/24"
	zone               = "%[5]s-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet_lb" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-lb"
	subnet             = "10.2.100.0/24"
	zone               = "%[5]s-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "LOADB"
}

resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = "%[4]s"
  login_key_name              = "%[3]s"
  lb_private_subnet_no        = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [
    ncloud_subnet.subnet1.id
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                     	  = "%[5]s-1"
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "%[1]s"
  node_count     = 1
  product_code   = "%[6]s"
  subnet_no      = ncloud_subnet.subnet1.id 
  autoscale {
    enabled = true
    min = 1
    max = 1
  }
}

data "ncloud_nks_node_pool" "node_pool"{
  cluster_uuid   = ncloud_nks_node_pool.node_pool.cluster_uuid
  node_pool_name = ncloud_nks_node_pool.node_pool.node_pool_name
}
`, testClusterName, clusterType, loginKey, version, region, productType)
}
