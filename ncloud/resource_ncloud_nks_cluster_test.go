package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Create LoginKey Before NKS Test
const TF_TEST_NKS_LOGIN_KEY = "tf-test-nks-login-key"

func TestAccResourceNcloudNKSCluster_basic(t *testing.T) {
	var cluster vnks.Cluster
	name := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.19.14-nks.1"
	resourceName := "ncloud_nks_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", clusterType),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", k8sVersion),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", TF_TEST_NKS_LOGIN_KEY),
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
	name := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.19.14"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudNKSCluster_InvalidSubnetConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY),
				ExpectError: regexp.MustCompile("중 하나여야 합니다."),
			},
		},
	})
}

func testAccResourceNcloudNKSClusterConfig(name string, clusterType string, k8sVersion string, loginKeyName string) string {
	return fmt.Sprintf(`
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

data "ncloud_nks_versions" "version" {
  filter {
    name = "value"
    values = ["%[3]s"]
    regex = true
  }
}

resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = data.ncloud_nks_versions.version.versions.0.value
  login_key_name              = "%[4]s"
  subnet_lb_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "KR-1"
}
`, name, clusterType, k8sVersion, loginKeyName)
}

func testAccResourceNcloudNKSCluster_InvalidSubnetConfig(name string, clusterType string, k8sVersion string, loginKeyName string) string {
	return fmt.Sprintf(`
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

data "ncloud_nks_versions" "version" {
  filter {
    name = "value"
    values = ["%[3]s"]
    regex = true
  }
}

resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = data.ncloud_nks_versions.version.versions.0.value
  login_key_name              = "%[4]s"
  subnet_lb_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "KR-1"
}
`, name, clusterType, k8sVersion, loginKeyName)
}

func testAccCheckNKSClusterExists(n string, cluster *vnks.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster uuid is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		resp, err := getNKSCluster(context.Background(), config, rs.Primary.ID)
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

		clusters, err := getNKSClusters(context.Background(), config)
		if err != nil {
			return err
		}

		for _, cluster := range clusters {
			if ncloud.StringValue(cluster.Uuid) == rs.Primary.ID {
				return fmt.Errorf("Cluster still exists")
			}
		}
	}

	return nil
}

func getTestClusterName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testClusterName := fmt.Sprintf("tf-%d-cluster", rInt)
	return testClusterName
}
