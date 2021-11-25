package ncloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/terraform-providers/terraform-provider-ncloud/sdk/vnks"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudNKSCluster_basic(t *testing.T) {
	var cluster vnks.Cluster
	name := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.19.14-nks.1"
	resourceName := "ncloud_nks_cluster.cluster"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", clusterType),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", k8sVersion),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", name),
					resource.TestCheckResourceAttr(resourceName, "zone", "KR-1"),
					resource.TestCheckResourceAttr(resourceName, "region_code", "KR"),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_InvalidSubnet(t *testing.T) {
	name := fmt.Sprintf("test-nksdisappears-%s", acctest.RandString(5))
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.19.14"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudNKSCluster_InvalidSubnetConfig(name, clusterType, k8sVersion),
				ExpectError: regexp.MustCompile("중 하나여야 합니다."),
			},
		},
	})
}

func testAccResourceNcloudNKSClusterConfig(name string, clusterType string, k8sVersion string) string {
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
  filter {
    name = "value"
    values = ["%[3]s"]
    regex = true
  }
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
  zone                        = "KR-1"
}
`, name, clusterType, k8sVersion)
}

func testAccResourceNcloudNKSCluster_InvalidSubnetConfig(name string, clusterType string, k8sVersion string) string {
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
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet2" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-2"
	subnet             = "10.2.2.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet_lb" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-lb"
	subnet             = "10.2.100.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "LOADB"
}

data "ncloud_nks_version" "version" {
  filter {
    name = "value"
    values = ["%[3]s"]
    regex = true
  }
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
  zone                        = "KR-1"
}
`, name, clusterType, k8sVersion)
}

func testAccCheckNKSClusterExists(n string, cluster *vnks.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster name is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		resp, err := getNKSClusterWithName(context.Background(), config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*cluster = *resp

		return nil
	}
}

func testAccCheckNKSClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nks_cluster" {
			continue
		}

		cluster, err := getNKSClusterWithName(context.Background(), config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if cluster != nil {
			return errors.New("Cluster still exists")
		}
	}

	return nil
}

func getTestClusterName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testClusterName := fmt.Sprintf("tf-%d-cluster", rInt)
	return testClusterName
}
