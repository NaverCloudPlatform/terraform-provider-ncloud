package ncloud

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudNKSNodePool(t *testing.T) {
	validateAcctestEnvironment(t)

	dataName := "data.ncloud_nks_node_pool.node_pool"
	resourceName := "ncloud_nks_node_pool.node_pool"
	clusterName := getTestClusterName()
	nksInfo, err := getNKSTestInfo("XEN")
	if err != nil {
		t.Error(err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSNodePoolConfig(clusterName, TF_TEST_NKS_LOGIN_KEY, nksInfo, 1),
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

func testAccDataSourceNKSNodePoolConfig(name string, loginKeyName string, nksInfo *NKSTestInfo, nodeCount int32) string {
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
    %[7]s
  ]
  vpc_no                      = %[8]s
  zone                        = "%[9]s-1"
`, name, nksInfo.ClusterType, nksInfo.K8sVersion, loginKeyName, *nksInfo.PrivateLbSubnetList[0].SubnetNo, nksInfo.HypervisorCode, *nksInfo.PrivateSubnetList[0].SubnetNo, *nksInfo.Vpc.VpcNo, nksInfo.Region))

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
    values = ["ubuntu-20.04"]
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

data "ncloud_nks_node_pool" "node_pool"{
  cluster_uuid   = ncloud_nks_node_pool.node_pool.cluster_uuid
  node_pool_name = ncloud_nks_node_pool.node_pool.node_pool_name
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid   = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "%[2]s"
  node_count     = %[3]d
  k8s_version    = "%[4]s"
  subnet_no_list = [ %[5]s ]
  autoscale {
    enabled = false
	max = 0
	min = 0
  }

  label {
    key = "foo"
    value = "bar"
  }

  taint {
    key = "foo"
    effect = "NoSchedule"
    value = "bar"
  }

  software_code = data.ncloud_nks_server_images.image.images.0.value
`, nksInfo.Region, name, nodeCount, nksInfo.K8sVersion, *nksInfo.PrivateSubnetList[0].SubnetNo))
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
