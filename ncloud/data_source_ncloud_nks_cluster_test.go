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
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSClusterConfig(testClusterName, clusterType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "uuid", resourceName, "uuid"),
					resource.TestCheckResourceAttrPair(dataName, "acg_name", resourceName, "acg_name"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "cluster_type", resourceName, "cluster_type"),
					resource.TestCheckResourceAttrPair(dataName, "created_at", resourceName, "created_at"),
					resource.TestCheckResourceAttrPair(dataName, "endpoint", resourceName, "endpoint"),
					resource.TestCheckResourceAttrPair(dataName, "region_code", resourceName, "region_code"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_name", resourceName, "subnet_name"),
					resource.TestCheckResourceAttrPair(dataName, "login_key_name", resourceName, "login_key_name"),
					resource.TestCheckResourceAttrPair(dataName, "k8s_version", resourceName, "k8s_version"),
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_name", resourceName, "vpc_name"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_lb_no", resourceName, "subnet_lb_no"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_lb_name", resourceName, "subnet_lb_name"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list.#", resourceName, "subnet_no_list.#"),
					resource.TestCheckResourceAttrPair(dataName, "subnet_no_list.0", resourceName, "subnet_no_list.0"),
				),
			},
		},
	})
}

func testAccDataSourceNKSClusterConfig(testClusterName string, clusterType string) string {
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

data "ncloud_nks_version" "version" {
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
  zone                     	  = "KR-1"
}

data "ncloud_nks_cluster" "cluster" {
	name = ncloud_nks_cluster.cluster.name
}


`, testClusterName, clusterType)
}
