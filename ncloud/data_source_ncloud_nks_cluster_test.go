package ncloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNKSCluster(t *testing.T) {
	dataName := "data.ncloud_nks_cluster.cluster"
	resourceName := "ncloud_nks_cluster.cluster"
	testClusterName := getTestClusterName()
	region, clusterType, _, k8sVersion := getRegionAndNKSType()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSClusterConfig(testClusterName, clusterType, TF_TEST_NKS_LOGIN_KEY, k8sVersion, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "uuid", resourceName, "uuid"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "cluster_type", resourceName, "cluster_type"),
					resource.TestCheckResourceAttrPair(dataName, "endpoint", resourceName, "endpoint"),
					resource.TestCheckResourceAttrPair(dataName, "login_key_name", resourceName, "login_key_name"),
					resource.TestCheckResourceAttrPair(dataName, "k8s_version", resourceName, "k8s_version"),
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "lb_private_subnet_no", resourceName, "lb_private_subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "lb_public_subnet_no", resourceName, "lb_public_subnet_no"),
					resource.TestCheckResourceAttrPair(dataName, "kube_network_plugin", resourceName, "kube_network_plugin"),
					resource.TestCheckResourceAttrPair(dataName, "public_network", resourceName, "public_network"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list.#", resourceName, "subnet_no_list.#"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list.0", resourceName, "subnet_no_list.0"),
					resource.TestCheckResourceAttrPair(dataName, "acg_no", resourceName, "acg_no"),
				),
			},
		},
	})
}

func testAccDataSourceNKSClusterConfig(testClusterName string, clusterType string, loginKey string, version string, region string) string {
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

resource "ncloud_subnet" "subnet2" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-2"
	subnet             = "10.2.2.0/24"
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
  kube_network_plugin         = "cilium"
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                     	  = "%[5]s-1"
}

data "ncloud_nks_cluster" "cluster" {
	uuid = ncloud_nks_cluster.cluster.uuid
}


`, testClusterName, clusterType, loginKey, version, region)
}
