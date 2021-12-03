package ncloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudNKSNodePool_basic(t *testing.T) {
	var nodePool vnks.NodePoolRes
	clusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_node_pool.node_pool"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, productCode, 2, TF_TEST_NKS_LOGIN_KEY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "node_pool_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "product_code", productCode),
					resource.TestCheckResourceAttr(resourceName, "autoscale.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "autoscale.0.min", "2"),
					resource.TestCheckResourceAttr(resourceName, "autoscale.0.max", "2"),
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

func TestAccResourceNcloudNKSNodePool_updateNodeCountAndAutoScale(t *testing.T) {
	var nodePool vnks.NodePoolRes
	clusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_node_pool.node_pool"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, productCode, 2, TF_TEST_NKS_LOGIN_KEY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "node_count", "2"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSNodePoolUpdateAutoScaleConfig(clusterName, clusterType, productCode, 1, TF_TEST_NKS_LOGIN_KEY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "autoscale.0.enabled", "false"),
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

func TestAccResourceNcloudNKSNodePool_invalidNodeCount(t *testing.T) {
	clusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, productCode, 0, TF_TEST_NKS_LOGIN_KEY),
				ExpectError: regexp.MustCompile("nodeCount must not be less than 1"),
			},
		},
	})
}

func testAccResourceNcloudNKSNodePoolConfig(name string, clusterType string, productCode string, nodeCount int, loginKey string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.2.1.0/24"
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
}

resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = data.ncloud_nks_versions.version.versions.0.value
  login_key_name              = "%[5]s"
  lb_subnet_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [
    ncloud_subnet.subnet.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "KR-1"
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid   = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "%[1]s"
  node_count     = %[4]d
  product_code   = "%[3]s"
  subnet_no      = ncloud_subnet.subnet.id
  autoscale {
    enabled = true
    min = 2
    max = 2
  }
}

`, name, clusterType, productCode, nodeCount, loginKey)
}

func testAccResourceNcloudNKSNodePoolUpdateAutoScaleConfig(name string, clusterType string, productCode string, nodeCount int, loginkey string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.2.1.0/24"
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
}

resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = data.ncloud_nks_versions.version.versions.0.value
  login_key_name              = "%[5]s"
  lb_subnet_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [
    ncloud_subnet.subnet.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "KR-1"
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid   = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "%[1]s"
  node_count     = %[4]d
  product_code   = "%[3]s"
  subnet_no      = ncloud_subnet.subnet.id
  autoscale {
    enabled = false
    min = 0
    max = 0
  }
}
`, name, clusterType, productCode, nodeCount, loginkey)
}

func testAccCheckNKSNodePoolExists(n string, nodePool *vnks.NodePoolRes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No nodepool no is set")
		}

		clusterUuid, nodePoolName, err := NodePoolParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Id(%s) is not [ClusterName:NodePoolName] ", rs.Primary.ID)
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		if err != nil {
			return err
		}

		np, err := getNKSNodePool(context.Background(), config, clusterUuid, nodePoolName)
		if err != nil {
			return err
		}

		*nodePool = *np

		return nil
	}
}

func testAccCheckNKSNodePoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nks_node_pool" {
			continue
		}

		clusterUuid, nodePoolName, err := NodePoolParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Id(%s) is not [ClusterName:NodePoolName] ", rs.Primary.ID)
		}

		clusters, err := getNKSClusters(context.Background(), config)
		if err != nil {
			return err
		}

		for _, cluster := range clusters {
			if ncloud.StringValue(cluster.Uuid) == clusterUuid {
				np, err := getNKSNodePool(context.Background(), config, clusterUuid, nodePoolName)
				if err != nil {
					return err
				}

				if np != nil {
					return errors.New("NodePool still exists")
				}
			}
		}

	}

	return nil
}