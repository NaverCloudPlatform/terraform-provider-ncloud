package nks_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/nks"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccResourceNcloudNKSNodePool_basic_XEN(t *testing.T) {
	validateAcctestEnvironment(t)

	clusterName := GetTestClusterName()
	resourceName := "ncloud_nks_node_pool.node_pool"

	nksInfo, err := getNKSTestInfo("XEN")
	if err != nil {
		t.Error(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 1),
				Check:  testAccResourceNcloudNKSNodePoolBasicCheck(resourceName, clusterName, nksInfo),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudNKSNodePool_basic_KVM(t *testing.T) {
	validateAcctestEnvironment(t)

	clusterName := GetTestClusterName()
	resourceName := "ncloud_nks_node_pool.node_pool"

	nksInfo, err := getNKSTestInfo("KVM")
	if err != nil {
		t.Error(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 1),
				Check:  testAccResourceNcloudNKSNodePoolBasicCheck(resourceName, clusterName, nksInfo),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudNKSNodePool_Update_XEN(t *testing.T) {
	validateAcctestEnvironment(t)

	var nodePool vnks.NodePool

	clusterName := fmt.Sprintf("m3-%s", GetTestClusterName())
	resourceName := "ncloud_nks_node_pool.node_pool"

	nksInfo, err := getNKSTestInfo("XEN")
	if err != nil {
		t.Error(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNKSNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSNodePoolConfigUpdateAll(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 2),
				Check:  testAccResourceNcloudNKSNodePoolUpdateAllCheck(resourceName, clusterName, nksInfo),
			},
		},
	})
}

func TestAccResourceNcloudNKSNodePool_Update_KVM(t *testing.T) {
	validateAcctestEnvironment(t)

	var nodePool vnks.NodePool
	clusterName := fmt.Sprintf("m3-%s", GetTestClusterName())
	resourceName := "ncloud_nks_node_pool.node_pool"

	nksInfo, err := getNKSTestInfo("KVM")
	if err != nil {
		t.Error(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNKSNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfig(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSNodePoolExists(resourceName, &nodePool),
				),
				Destroy: false,
			},
			{
				PreConfig: func() {
					time.Sleep(30 * time.Minute)
				},
				Config:  testAccResourceNcloudNKSNodePoolConfigUpdateAll(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 2),
				Check:   testAccResourceNcloudNKSNodePoolUpdateAllCheck(resourceName, clusterName, nksInfo),
				Destroy: false,
			},
		},
	})
}

func TestAccResourceNcloudNKSNodePool_publicNetwork_XEN(t *testing.T) {
	validateAcctestEnvironment(t)

	clusterName := GetTestClusterName()
	resourceName := "ncloud_nks_node_pool.node_pool"

	nksInfo, err := getNKSTestInfo("XEN")
	if err != nil {
		t.Error(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfigPublicNetwork(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 1),
				Check:  testAccResourceNcloudNKSNodePoolPublicNetworkCheck(resourceName, clusterName),
			},
		},
	})
}

func TestAccResourceNcloudNKSNodePool_publicNetwork_KVM(t *testing.T) {
	validateAcctestEnvironment(t)

	clusterName := GetTestClusterName()
	resourceName := "ncloud_nks_node_pool.node_pool"

	nksInfo, err := getNKSTestInfo("KVM")
	if err != nil {
		t.Error(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSNodePoolConfigPublicNetwork(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 1),
				Check:  testAccResourceNcloudNKSNodePoolPublicNetworkCheck(resourceName, clusterName),
			},
		},
	})
}

func testAccResourceNcloudNKSNodePoolConfig(name string, loginKeyName string, nksInfo *NKSTestInfo, nodeCount int32) string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf(`
resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = %[5]s
  hypervisor_code             = "%[6]s"
  kube_network_plugin         = "cilium"
  subnet_no_list              = [
    %[7]s,  %[10]s
  ]
  vpc_no                      = %[8]s
  zone                        = "%[9]s-1"
`, name, nksInfo.ClusterType, nksInfo.K8sVersion, loginKeyName, *nksInfo.PrivateLbSubnetList[0].SubnetNo, nksInfo.HypervisorCode, *nksInfo.PrivateSubnetList[0].SubnetNo, *nksInfo.Vpc.VpcNo, nksInfo.Region, *nksInfo.PrivateSubnetList[1].SubnetNo))

	if nksInfo.needPublicLb {
		b.WriteString(fmt.Sprintf(`
  lb_public_subnet_no = %[1]s
`, *nksInfo.PublicLbSubnetList[0].SubnetNo))
	}

	b.WriteString(`
}
`)

	b.WriteString(fmt.Sprintf(`
data "ncloud_nks_server_images" "image"{
  hypervisor_code = ncloud_nks_cluster.cluster.hypervisor_code
    filter {
    name = "label"
    values = ["%[6]s"]
    regex = true
  }

}

data "ncloud_nks_server_products" "product"{
  software_code = data.ncloud_nks_server_images.image.images[0].value
  zone = "%[1]s-1"
  filter {
    name = "product_type"
    values = [ "STAND"]
  }
  
  filter {
    name = "cpu_count"
    values = [ "2"]
  }
  
  filter {
    name = "memory_size"
    values = [ "8GB" ]
  }
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid   = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "%[2]s"
  node_count     = %[3]d
  k8s_version    = "%[4]s"
  subnet_no_list = [ %[5]s ]
  autoscale {
    enabled = false
	min = 0
    max = 0
  }

  label {
    key = "foo"
    value = "bar"
  }

  taint {
    key = "foo"
    effect = "PreferNoSchedule"
    value = "bar"
  }

  software_code = data.ncloud_nks_server_images.image.images.0.value
`, nksInfo.Region, name, nodeCount, nksInfo.K8sVersion, *nksInfo.PrivateSubnetList[0].SubnetNo, nksInfo.UbuntuImageVersion))
	if nksInfo.HypervisorCode == "KVM" {
		b.WriteString(`
  server_spec_code = data.ncloud_nks_server_products.product.products.0.value
  storage_size = 100
}
		`)

	} else {
		b.WriteString(`
  product_code = data.ncloud_nks_server_products.product.products.0.value
}
		`)
	}
	return b.String()

}

func testAccResourceNcloudNKSNodePoolConfigPublicNetwork(name string, loginKeyName string, nksInfo *NKSTestInfo, nodeCount int32) string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf(`
resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = %[5]s
  hypervisor_code             = "%[6]s"
  kube_network_plugin         = "cilium"
  auth_type                   = "CONFIG_MAP"
  subnet_no_list              = [
    %[7]s
  ]
  vpc_no                      = %[8]s
  zone                        = "%[9]s-1"
  public_network              = true
`, name, nksInfo.ClusterType, nksInfo.K8sVersion, loginKeyName, *nksInfo.PrivateLbSubnetList[0].SubnetNo, nksInfo.HypervisorCode, *nksInfo.PublicSubnetList[0].SubnetNo, *nksInfo.Vpc.VpcNo, nksInfo.Region))

	if nksInfo.needPublicLb {
		b.WriteString(fmt.Sprintf(`
  lb_public_subnet_no = %[1]s
`, *nksInfo.PublicLbSubnetList[0].SubnetNo))
	}

	b.WriteString(`
}
`)

	b.WriteString(fmt.Sprintf(`
data "ncloud_nks_server_images" "image"{
  hypervisor_code = ncloud_nks_cluster.cluster.hypervisor_code
    filter {
    name = "label"
    values = ["%[6]s"]
    regex = true
  }

}

data "ncloud_nks_server_products" "product"{
  software_code = data.ncloud_nks_server_images.image.images[0].value
  zone = "%[1]s-1"
  filter {
    name = "product_type"
    values = [ "STAND"]
  }
  
  filter {
    name = "cpu_count"
    values = [ "2"]
  }
  
  filter {
    name = "memory_size"
    values = [ "8GB" ]
  }
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid   = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "%[2]s"
  node_count     = %[3]d
  k8s_version    = "%[4]s"
  subnet_no_list = [ %[5]s ]
  autoscale {
    enabled = false
    min = 0
    max = 0
  }

  label {
    key = "bar"
    value = "foo"
  }

  taint {
    key = "bar"
    effect = "PreferNoSchedule"
    value = "foo"
  }

  software_code = data.ncloud_nks_server_images.image.images.0.value
`, nksInfo.Region, name, nodeCount, nksInfo.K8sVersion, *nksInfo.PublicSubnetList[0].SubnetNo, nksInfo.UbuntuImageVersion))
	if nksInfo.HypervisorCode == "KVM" {
		b.WriteString(`
  server_spec_code = data.ncloud_nks_server_products.product.products.0.value
  storage_size = 100
}
		`)

	} else {
		b.WriteString(`
  product_code = data.ncloud_nks_server_products.product.products.0.value
}
		`)
	}
	return b.String()

}

func testAccResourceNcloudNKSNodePoolConfigUpdateAll(name string, loginKeyName string, nksInfo *NKSTestInfo, nodeCount int32) string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf(`
resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = %[5]s
  hypervisor_code             = "%[6]s"
  kube_network_plugin         = "cilium"
  subnet_no_list              = [
    %[7]s,  %[10]s
  ]
  vpc_no                      = %[8]s
  zone                        = "%[9]s-1"
`, name, nksInfo.ClusterType, nksInfo.UpgradeK8sVersion, loginKeyName, *nksInfo.PrivateLbSubnetList[0].SubnetNo, nksInfo.HypervisorCode, *nksInfo.PrivateSubnetList[0].SubnetNo, *nksInfo.Vpc.VpcNo, nksInfo.Region, *nksInfo.PrivateSubnetList[1].SubnetNo))

	if nksInfo.needPublicLb {
		b.WriteString(fmt.Sprintf(`
  lb_public_subnet_no = %[1]s
`, *nksInfo.PublicLbSubnetList[0].SubnetNo))
	}

	b.WriteString(`
}
`)

	b.WriteString(fmt.Sprintf(`
data "ncloud_nks_server_images" "image"{
  hypervisor_code = ncloud_nks_cluster.cluster.hypervisor_code
    filter {
    name = "label"
    values = ["%[7]s"]
    regex = true
  }

}

data "ncloud_nks_server_products" "product"{
  software_code = data.ncloud_nks_server_images.image.images[0].value
  zone = "%[1]s-1"
  filter {
    name = "product_type"
    values = [ "STAND"]
  }
  
  filter {
    name = "cpu_count"
    values = [ "2"]
  }
  
  filter {
    name = "memory_size"
    values = [ "8GB" ]
  }
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid   = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "%[2]s"

  k8s_version    = "%[4]s"
  subnet_no_list = [ %[5]s, %[6]s ]
  autoscale {
    enabled = true
    min = 1
    max = 2
  }

  label {
    key = "bar"
    value = "foo"
  }

  taint {
    key = "bar"
    effect = "PreferNoSchedule"
    value = ""
  }

  software_code = data.ncloud_nks_server_images.image.images.0.value
`, nksInfo.Region, name, nodeCount, nksInfo.UpgradeK8sVersion, *nksInfo.PrivateSubnetList[0].SubnetNo, *nksInfo.PrivateSubnetList[1].SubnetNo, nksInfo.UbuntuImageVersion))
	if nksInfo.HypervisorCode == "KVM" {
		b.WriteString(`
  server_spec_code = data.ncloud_nks_server_products.product.products.0.value
  storage_size = 100

  lifecycle {
    ignore_changes = [software_code,server_spec_code ]
  }
}
		`)

	} else {
		b.WriteString(`
  product_code = data.ncloud_nks_server_products.product.products.0.value
}
		`)
	}
	return b.String()

}

func testAccResourceNcloudNKSNodePoolBasicCheck(resourceName string, name string, nksInfo *NKSTestInfo) (check resource.TestCheckFunc) {
	var nodePool vnks.NodePool
	check = resource.ComposeTestCheckFunc(
		testAccCheckNKSNodePoolExists(resourceName, &nodePool),
		resource.TestCheckResourceAttr(resourceName, "node_pool_name", name),
		resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
		resource.TestCheckResourceAttr(resourceName, "autoscale.0.enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "autoscale.0.min", "0"),
		resource.TestCheckResourceAttr(resourceName, "autoscale.0.max", "0"),
		resource.TestCheckResourceAttr(resourceName, "k8s_version", nksInfo.K8sVersion),
		resource.TestCheckResourceAttr(resourceName, "subnet_no_list.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "label.0.key", "foo"),
		resource.TestCheckResourceAttr(resourceName, "label.0.value", "bar"),
		resource.TestCheckResourceAttr(resourceName, "taint.0.key", "foo"),
		resource.TestCheckResourceAttr(resourceName, "taint.0.value", "bar"),
		resource.TestCheckResourceAttr(resourceName, "taint.0.effect", "PreferNoSchedule"),
	)

	if nksInfo.HypervisorCode == "KVM" {
		check = resource.ComposeTestCheckFunc(
			check,
			resource.TestCheckResourceAttr(resourceName, "storage_size", "100"),
		)
	}
	return check
}

func testAccResourceNcloudNKSNodePoolPublicNetworkCheck(resourceName string, name string) (check resource.TestCheckFunc) {
	var nodePool vnks.NodePool
	return resource.ComposeTestCheckFunc(
		testAccCheckNKSNodePoolExists(resourceName, &nodePool),
		resource.TestCheckResourceAttr(resourceName, "node_pool_name", name),
		resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
	)
}

func testAccResourceNcloudNKSNodePoolUpdateAllCheck(resourceName string, name string, nksInfo *NKSTestInfo) (check resource.TestCheckFunc) {
	var nodePool vnks.NodePool
	return resource.ComposeTestCheckFunc(
		testAccCheckNKSNodePoolExists(resourceName, &nodePool),
		resource.TestCheckResourceAttr(resourceName, "node_pool_name", name),
		resource.TestCheckResourceAttr(resourceName, "autoscale.0.enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "autoscale.0.min", "1"),
		resource.TestCheckResourceAttr(resourceName, "autoscale.0.max", "2"),
		resource.TestCheckResourceAttr(resourceName, "k8s_version", nksInfo.UpgradeK8sVersion),
		resource.TestCheckResourceAttr(resourceName, "subnet_no_list.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "label.0.key", "bar"),
		resource.TestCheckResourceAttr(resourceName, "label.0.value", "foo"),
		resource.TestCheckResourceAttr(resourceName, "taint.0.key", "bar"),
		resource.TestCheckResourceAttr(resourceName, "taint.0.value", ""),
		resource.TestCheckResourceAttr(resourceName, "taint.0.effect", "PreferNoSchedule"),
	)
}

func testAccCheckNKSNodePoolExists(n string, nodePool *vnks.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No nodepool no is set")
		}

		clusterUuid, nodePoolName, err := nks.NodePoolParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Id(%s) is not [ClusterName:NodePoolName] ", rs.Primary.ID)
		}

		config := TestAccProvider.Meta().(*conn.ProviderConfig)
		if err != nil {
			return err
		}

		np, err := nks.GetNKSNodePool(context.Background(), config, clusterUuid, nodePoolName)
		if err != nil {
			return err
		}

		*nodePool = *np

		return nil
	}
}

func testAccCheckNKSNodePoolDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nks_node_pool" {
			continue
		}

		clusterUuid, nodePoolName, err := nks.NodePoolParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Id(%s) is not [ClusterName:NodePoolName] ", rs.Primary.ID)
		}

		clusters, err := nks.GetNKSClusters(context.Background(), config)
		if err != nil {
			return err
		}

		for _, cluster := range clusters {
			if ncloud.StringValue(cluster.Uuid) == clusterUuid {
				np, err := nks.GetNKSNodePool(context.Background(), config, clusterUuid, nodePoolName)
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

func testAccCheckNKSClusterDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nks_cluster" {
			continue
		}

		clusters, err := nks.GetNKSClusters(context.Background(), config)
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
