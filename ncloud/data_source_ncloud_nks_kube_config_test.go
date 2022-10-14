package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNKSKubeConfig(t *testing.T) {
	dataName := "data.ncloud_nks_kube_config.kube_config"
	resourceName := "ncloud_nks_cluster.cluster"
	testClusterName := getTestClusterName()
	region, clusterType, _, k8sVersion := getRegionAndNKSType()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSKubeConfigConfig(testClusterName, clusterType, TF_TEST_NKS_LOGIN_KEY, k8sVersion, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "cluster_uuid", resourceName, "uuid"),
					resource.TestCheckResourceAttrPair(dataName, "host", resourceName, "endpoint"),
				),
			},
		},
	})
}

func testAccDataSourceNKSKubeConfigConfig(testClusterName string, clusterType string, loginKey string, version string, region string) string {
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
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                     	  = "%[5]s-1"
}

data "ncloud_nks_kube_config" "kube_config" {
	cluster_uuid = ncloud_nks_cluster.cluster.uuid
}


`, testClusterName, clusterType, loginKey, version, region)
}
