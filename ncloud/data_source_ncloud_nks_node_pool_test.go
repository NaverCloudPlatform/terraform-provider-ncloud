package ncloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNKSNodePool_basic(t *testing.T) {
	dataName := "data.ncloud_nks_node_pool.node_pool"
	resourceName := "ncloud_nks_node_pool.node_pool"
	testClusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSNodePoolConfig(testClusterName, clusterType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					//resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrPair(dataName, "cluster_name", resourceName, "cluster_name"),
					resource.TestCheckResourceAttrPair(dataName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "k8s_version", resourceName, "k8s_version"),
					resource.TestCheckResourceAttrPair(dataName, "node_pool_name", resourceName, "node_pool_name"),
					resource.TestCheckResourceAttrPair(dataName, "node_count", resourceName, "node_count"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list", resourceName, "subnet_no_list"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_name_list", resourceName, "subnet_name_list"),
					resource.TestCheckResourceAttrPair(dataName, "product_code", resourceName, "product_code"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "autoscale.0.enabled", resourceName, "autoscale.0.enabled"),
					resource.TestCheckResourceAttrPair(dataName, "autoscale.0.min", resourceName, "autoscale.0.min"),
					resource.TestCheckResourceAttrPair(dataName, "autoscale.0.max", resourceName, "autoscale.0.max"),
				),
			},
			{
				Config:  testAccDataSourceNKSNodePoolConfig(testClusterName, clusterType),
				Destroy: true,
			},
		},
	})
}

func testAccDataSourceNKSNodePoolConfig(testClusterName string, clusterType string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
  key_name = "%[1]s"
}

resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "subnet1" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-1"
	subnet             = "10.2.1.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet2" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-2"
	subnet             = "10.2.2.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet_lb" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-lb"
	subnet             = "10.2.100.0/24"
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "LOADB"
}

data "ncloud_nks_version" "version"{
}


resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = data.ncloud_nks_version.version.versions.0.value
  login_key_name              = ncloud_login_key.loginkey.key_name
  subnet_lb_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone_no                     = "2"
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_name = ncloud_nks_cluster.cluster.name
  node_pool_name = "%[1]s"
  node_count     = 1
  product_code   = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
  subnet_no_list              = [ ncloud_subnet.subnet1.id ]
  autoscale {
    enabled = true
    min = 1
    max = 2
  }
}

data "ncloud_nks_node_pool" "node_pool"{
  cluster_name    = ncloud_nks_cluster.cluster.name
  node_pool_name = "%[1]s"
}
`, testClusterName, clusterType)
}
