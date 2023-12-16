package nks_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudNKSKubeConfig(t *testing.T) {
	validateAcctestEnvironment(t)

	dataName := "data.ncloud_nks_kube_config.kube_config"
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
				Config: testAccDataSourceNKSKubeConfigConfig(name, TF_TEST_NKS_LOGIN_KEY, true, nksInfo),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "cluster_uuid", resourceName, "uuid"),
					resource.TestCheckResourceAttrPair(dataName, "host", resourceName, "endpoint"),
				),
			},
		},
	})
}

func testAccDataSourceNKSKubeConfigConfig(name string, loginKeyName string, auditLog bool, nksInfo *NKSTestInfo) string {
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
    issuer_url                = "https://keycloak.url/realms/nks"
    client_id                 = "nks-client"
    username_claim            = "preferred_username"
    username_prefix           = "oidc:"
    groups_claim              = "groups"
    groups_prefix             = "oidc:"
    required_claim           = "iss=https://keycloak.url/realms/nks"
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
	data "ncloud_nks_kube_config" "kube_config" {
		cluster_uuid = ncloud_nks_cluster.cluster.uuid
	}
`)
	return b.String()
}
