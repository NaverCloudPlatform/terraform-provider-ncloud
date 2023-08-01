package nks_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/nks"
)

// Create LoginKey Before NKS Test
const TF_TEST_NKS_LOGIN_KEY = "tf-test-nks-login-key"

func TestAccResourceNcloudNKSCluster_basic(t *testing.T) {
	var cluster vnks.Cluster
	name := GetTestClusterName()

	resourceName := "ncloud_nks_cluster.cluster"

	region, clusterType, _, k8sVersion := getRegionAndNKSType()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", clusterType),
					resource.TestMatchResourceAttr(resourceName, "k8s_version", regexp.MustCompile(k8sVersion)),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", TF_TEST_NKS_LOGIN_KEY),
					resource.TestCheckResourceAttr(resourceName, "zone", fmt.Sprintf("%s-1", region)),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "log.0.audit", "true"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.issuer_url", "https://keycloak.ncp.gimmetm.net/realms/nks"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.client_id", "nks-client"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.username_claim", "preferred_username"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.username_prefix", "oidc:"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_claim", "groups"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_prefix", "oidc:"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.required_claim", "iss=https://keycloak.ncp.gimmetm.net/realms/nks"),
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

func TestAccResourceNcloudNKSCluster_public_network(t *testing.T) {
	var cluster vnks.Cluster
	name := GetTestClusterName()
	resourceName := "ncloud_nks_cluster.cluster"

	region, clusterType, _, k8sVersion := getRegionAndNKSType()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterPublicNetworkConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", clusterType),
					resource.TestMatchResourceAttr(resourceName, "k8s_version", regexp.MustCompile(k8sVersion)),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", TF_TEST_NKS_LOGIN_KEY),
					resource.TestCheckResourceAttr(resourceName, "public_network", "true"),
					resource.TestCheckResourceAttr(resourceName, "zone", fmt.Sprintf("%s-1", region)),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_InvalidSubnet(t *testing.T) {
	name := GetTestClusterName()

	region, clusterType, _, k8sVersion := getRegionAndNKSType()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceNcloudNKSCluster_InvalidSubnetConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region),
				ExpectError: regexp.MustCompile(getInvalidSubnetExpectError()),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_Update(t *testing.T) {
	var cluster vnks.Cluster
	name := "m3-" + GetTestClusterName()

	region, clusterType, _, k8sVersion := getRegionAndNKSType()
	resourceName := "ncloud_nks_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "log.0.audit", "true"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "log.0.audit", "false"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSCluster_NoOIDCSpec(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "oidc.#", "0"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSCluster_AddSubnet(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "subnet_no_list.#", "3"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSCluster_AddSubnet(name, clusterType, "1.25.8-nks.1", TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", "1.25.8-nks.1"),
				),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_UpdateOnce(t *testing.T) {
	var cluster vnks.Cluster
	name := "m3-" + GetTestClusterName()

	region, clusterType, _, k8sVersion := getRegionAndNKSType()
	resourceName := "ncloud_nks_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSCluster_AddSubnet(name, clusterType, "1.25.8-nks.1", TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", "1.25.8-nks.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_no_list.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "log.0.audit", "false"),
				),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_VersionUpgrade(t *testing.T) {
	var cluster vnks.Cluster
	name := "m3-" + GetTestClusterName()

	region, clusterType, _, k8sVersion := getRegionAndNKSType()
	resourceName := "ncloud_nks_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, "1.25.8-nks.1", TF_TEST_NKS_LOGIN_KEY, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", "1.25.8-nks.1"),
				),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_OIDCSpec(t *testing.T) {
	var cluster vnks.Cluster
	name := GetTestClusterName()

	region, clusterType, _, k8sVersion := getRegionAndNKSType()
	resourceName := "ncloud_nks_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.issuer_url", "https://keycloak.ncp.gimmetm.net/realms/nks"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.client_id", "nks-client"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.username_claim", "preferred_username"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.username_prefix", "oidc:"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_claim", "groups"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_prefix", "oidc:"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.required_claim", "iss=https://keycloak.ncp.gimmetm.net/realms/nks"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSCluster_NoOIDCSpec(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "oidc.#", "0"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.issuer_url", "https://keycloak.ncp.gimmetm.net/realms/nks"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.client_id", "nks-client"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.username_claim", "preferred_username"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.username_prefix", "oidc:"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_claim", "groups"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_prefix", "oidc:"),
					resource.TestCheckResourceAttr(resourceName, "oidc.0.required_claim", "iss=https://keycloak.ncp.gimmetm.net/realms/nks"),
				),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_AuditLog(t *testing.T) {
	var cluster vnks.Cluster
	name := GetTestClusterName()

	region, clusterType, _, k8sVersion := getRegionAndNKSType()
	resourceName := "ncloud_nks_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "log.0.audit", "false"),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "log.0.audit", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_AddSubnet(t *testing.T) {
	var cluster vnks.Cluster
	name := GetTestClusterName()

	region, clusterType, _, k8sVersion := getRegionAndNKSType()
	resourceName := "ncloud_nks_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
				),
				Destroy: false,
			},
			{
				Config: testAccResourceNcloudNKSCluster_AddSubnet(name, clusterType, k8sVersion, TF_TEST_NKS_LOGIN_KEY, region, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "subnet_no_list.#", "3"),
				),
				Destroy: false,
			},
		},
	})
}

func testAccResourceNcloudNKSClusterConfig(name string, clusterType string, k8sVersion string, loginKeyName string, region string, auditLog bool) string {
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
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = ncloud_subnet.subnet_lb.id
  kube_network_plugin         = "cilium"
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "%[5]s-1"
  log {
    audit                     = %[6]t
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
}
`, name, clusterType, k8sVersion, loginKeyName, region, auditLog)
}

func testAccResourceNcloudNKSClusterPublicNetworkConfig(name string, clusterType string, k8sVersion string, loginKeyName string, region string) string {
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
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet2" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-2"
	subnet             = "10.2.2.0/24"
	zone               = "%[5]s-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
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
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = ncloud_subnet.subnet_lb.id
  kube_network_plugin         = "cilium"
  public_network              = "true"
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "%[5]s-1"
}
`, name, clusterType, k8sVersion, loginKeyName, region)
}

func testAccResourceNcloudNKSCluster_InvalidSubnetConfig(name string, clusterType string, k8sVersion string, loginKeyName string, region string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.2.0.0/16"
}

resource "ncloud_subnet" "subnet1" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-1"
	subnet             = "10.2.1.0/24"
	zone               = "%[5]s-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet2" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-2"
	subnet             = "10.2.2.0/24"
	zone               = "%[5]s-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "subnet_lb" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-lb"
	subnet             = "10.2.100.0/24"
	zone               = "%[5]s-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "LOADB"
}

resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = ncloud_subnet.subnet_lb.id
  kube_network_plugin		  = "cilium"
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "%[5]s-1"
  log {
	audit = true
  }
}
`, name, clusterType, k8sVersion, loginKeyName, region)
}

func testAccResourceNcloudNKSCluster_NoOIDCSpec(name string, clusterType string, k8sVersion string, loginKeyName string, region string, auditLog bool) string {
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
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = ncloud_subnet.subnet_lb.id
  kube_network_plugin         = "cilium"
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "%[5]s-1"
  log {
    audit                     = "%[6]t"
  }
}
`, name, clusterType, k8sVersion, loginKeyName, region, auditLog)
}

func testAccResourceNcloudNKSCluster_AddSubnet(name string, clusterType string, k8sVersion string, loginKeyName string, region string, auditLog bool) string {
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

resource "ncloud_subnet" "subnet3" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s-3"
	subnet             = "10.2.4.0/24"
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
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  lb_private_subnet_no        = ncloud_subnet.subnet_lb.id
  kube_network_plugin         = "cilium"
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
    ncloud_subnet.subnet3.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone                        = "%[5]s-1"
  log {
    audit                     = "%[6]t"
  }
}
`, name, clusterType, k8sVersion, loginKeyName, region, auditLog)
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

		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		resp, err := nks.GetNKSCluster(context.Background(), config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*cluster = *resp

		return nil
	}
}

func getRegionAndNKSType() (region string, clusterType string, productType string, k8sVersion string) {
	region = os.Getenv("NCLOUD_REGION")
	if region == "FKR" {
		clusterType = "SVR.VNKS.STAND.C002.M008.NET.HDD.B050.G001"
		productType = "SVR.VSVR.STAND.C002.M004.NET.SSD.B050.G001"
	} else {
		clusterType = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
		productType = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	}
	k8sVersion = "1.24.10-nks.1"
	return
}

func getInvalidSubnetExpectError() string {
	apigw := os.Getenv("NCLOUD_API_GW")
	if strings.Contains(apigw, "gov-ntruss.com") {
		return "Not found zone"
	}
	return "Subnet is undefined"
}
