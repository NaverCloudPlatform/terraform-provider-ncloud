package nks_test

import (
	"bytes"
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudNKSCluster(t *testing.T) {
	dataName := "data.ncloud_nks_cluster.cluster"
	resourceName := "ncloud_nks_cluster.cluster"
	name := GetTestClusterName()
	nksInfo, err := getNKSTestInfo("XEN")
	if err != nil {
		t.Error(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSClusterConfig(name, TF_TEST_NKS_LOGIN_KEY, true, nksInfo),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
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

func testAccDataSourceNKSClusterConfig(name string, loginKeyName string, auditLog bool, nksInfo *NKSTestInfo) string {
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
  log {
    audit                     = %[10]t
  }
  oidc {
    issuer_url                = "https://keycloak.ncp.gimmetm.net/realms/nks"
    client_id                 = "nks-client"
    username_claim            = "preferred_username"
    username_prefix           = "oidc:"
    groups_claim              = "groups"
    groups_prefix             = "oidc:"
    required_claim           = "iss=https://keycloak.ncp.gimmetm.net/realms/nks"
  }
`, name, nksInfo.ClusterType, nksInfo.K8sVersion, loginKeyName, *nksInfo.PrivateLbSubnetList[0].SubnetNo, nksInfo.HypervisorCode, *nksInfo.PrivateSubnetList[0].SubnetNo, *nksInfo.Vpc.VpcNo, nksInfo.Region, auditLog))

	if nksInfo.needPublicLb {
		b.WriteString(fmt.Sprintf(`
  lb_public_subnet_no = %[1]s
`, *nksInfo.PublicLbSubnetList[0].SubnetNo))
	}

	b.WriteString(`
}
`)

	b.WriteString(`
data "ncloud_nks_cluster" "cluster" {
	uuid = ncloud_nks_cluster.cluster.uuid
}
`)
	return b.String()
}
