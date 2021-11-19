package ncloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/terraform-providers/terraform-provider-ncloud/sdk/vnks"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudNKSNodePool_basic(t *testing.T) {
	var nodePool vnks.NodePoolRes
	clusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.20.11"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_node_pool.node_pool"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, k8sVersion, productCode, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_name", resourceName, "cluster_name"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_no", resourceName, "instance_no"),
					resource.TestCheckResourceAttrPair(resourceName, "k8s_version", resourceName, "k8s_version"),
					resource.TestCheckResourceAttrPair(resourceName, "node_pool_name", resourceName, "node_pool_name"),
					resource.TestCheckResourceAttrPair(resourceName, "node_count", resourceName, "node_count"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_no_list", resourceName, "subnet_no_list"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_name_list", resourceName, "subnet_name_list"),
					resource.TestCheckResourceAttrPair(resourceName, "product_code", resourceName, "product_code"),
					resource.TestCheckResourceAttrPair(resourceName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(resourceName, "autoscale.0.enabled", resourceName, "autoscale.0.enabled"),
					resource.TestCheckResourceAttrPair(resourceName, "autoscale.0.min", resourceName, "autoscale.0.min"),
					resource.TestCheckResourceAttrPair(resourceName, "autoscale.0.max", resourceName, "autoscale.0.max"),
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

func TestAccResourceNcloudNKSNodePool_disappears(t *testing.T) {
	var nodePool vnks.NodePoolRes
	clusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.20.11"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_node_pool.node_pool"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, k8sVersion, productCode, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
					testAccCheckNKSNodePoolDisappears(clusterName, &nodePool),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudSubnet_updateNodeCount(t *testing.T) {
	var nodePool vnks.NodePoolRes
	clusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.20.11"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_node_pool.node_pool"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, k8sVersion, productCode, 1),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
				),
			},
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName+"1", clusterType, k8sVersion, productCode, 1),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
				),
				ExpectError: regexp.MustCompile("Change 'name' is not support, Please set `name` as a old value"),
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
	k8sVersion := "1.20.11"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_node_pool.node_pool"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, k8sVersion, productCode, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
				),
			},
			{
				Config: testAccResourceNcloudNKSNodePoolUpdateAutoScaleConfig(clusterName, clusterType, k8sVersion, productCode, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "autoscale.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccResourceNcloudNodePool_InvalidNodeCount(t *testing.T) {
	clusterName := getTestClusterName()
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.20.11"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudNKSNodePoolConfig(clusterName, clusterType, k8sVersion, productCode, 0),
				ExpectError: regexp.MustCompile("The subnet must belong to the IPv4 CIDR of the specified VPC."),
			},
		},
	})
}

func testAccResourceNcloudNKSNodePoolConfig(name string, clusterType string, k8sVersion string, productCode string, nodeCount int) string {
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

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_name = ncloud_nks_cluster.cluster.name
  node_pool_name = "%[1]s"
  node_count     = "%[4]s"
  product_code   = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
  subnet_no_list              = [ ncloud_subnet.subnet.id ]
  autoscale {
    enabled = true
    min = 1
    max = 2
  }
}

`, name, clusterType, k8sVersion, productCode, nodeCount)
}

func testAccResourceNcloudNKSNodePoolUpdateAutoScaleConfig(name string, clusterType string, k8sVersion string, productCode string, nodeCount int) string {
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

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_name = ncloud_nks_cluster.cluster.name
  node_pool_name = "%[1]s"
  node_count     = "%[4]s"
  product_code   = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
  subnet_no_list              = [ ncloud_subnet.subnet.id ]
  autoscale {
    enabled = false
    min = 0
    max = 0
  }
}

`, name, clusterType, k8sVersion, productCode, nodeCount)
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

		clusterName, nodePoolName, err := NodePoolParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Id(%s) is not [ClusterName:NodePoolName] ", rs.Primary.ID)
		}
		config := testAccProvider.Meta().(*ProviderConfig)
		np, err := getNKSNodePool(context.Background(), config, &clusterName, &nodePoolName)
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

		clusterName, nodePoolName, err := NodePoolParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Id(%s) is not [ClusterName:NodePoolName] ", rs.Primary.ID)
		}
		np, err := getNKSNodePool(context.Background(), config, &clusterName, &nodePoolName)
		if err != nil {
			return err
		}

		if np != nil {
			return errors.New("Subnet still exists")
		}
	}

	return nil
}

func testAccCheckNKSNodePoolDisappears(clusterName string, np *vnks.NodePoolRes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		cluster, err := getNKSClusterWithName(context.Background(), config, clusterName)
		if err != nil {
			return err
		}

		intanceNo := fmt.Sprintf("%d", *np.InstanceNo)
		err = config.Client.vnks.V2Api.ClustersUuidNodePoolInstanceNoDelete(context.Background(), cluster.Uuid, &intanceNo)
		if err != nil {
			return err
		}

		d := resourceNcloudNKSNodePool().TestResourceData()
		d.SetId(NodePoolCreateResourceID(clusterName, *np.Name))
		if err := waitForNKSNodePoolDeletion(context.Background(), d, config); err != nil {
			return err
		}

		return err
	}
}
