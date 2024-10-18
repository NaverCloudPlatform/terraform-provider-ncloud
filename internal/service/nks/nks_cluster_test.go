package nks_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/nks"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

// Create LoginKey Before NKS Test
const TF_TEST_NKS_LOGIN_KEY = "tf-test-nks-login-key"

type NKSTestInfo struct {
	Vpc                 *vpc.Vpc
	DefaultAcl          *vpc.NetworkAcl
	PrivateSubnetList   []*vpc.Subnet
	PublicSubnetList    []*vpc.Subnet
	PrivateLbSubnetList []*vpc.Subnet
	PublicLbSubnetList  []*vpc.Subnet
	Region              string
	ClusterType         string
	ProductType         string
	K8sVersion          string
	UpgradeK8sVersion   string
	HypervisorCode      string
	IsFin               bool
	IsCaaS              bool
	needPublicLb        bool
}

func validateAcctestEnvironment(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Accetance Test skipped unless env 'TF_ACC' is set")
	}
}

func TestAccResourceNcloudNKSCluster_basic_XEN(t *testing.T) {
	validateAcctestEnvironment(t)

	name := GetTestClusterName()

	resourceName := "ncloud_nks_cluster.cluster"

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
				Config: testAccResourceNcloudNKSClusterDefaultConfig(name, TF_TEST_NKS_LOGIN_KEY, true, nksInfo),
				Check:  testAccResourceNcloudNKSClusterDefaultConfigCheck(resourceName, name, nksInfo),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_basic_KVM(t *testing.T) {
	validateAcctestEnvironment(t)

	name := GetTestClusterName()

	resourceName := "ncloud_nks_cluster.cluster"

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
				Config: testAccResourceNcloudNKSClusterDefaultConfig(name, TF_TEST_NKS_LOGIN_KEY, true, nksInfo),
				Check:  testAccResourceNcloudNKSClusterDefaultConfigCheck(resourceName, name, nksInfo),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_public_network_XEN(t *testing.T) {
	validateAcctestEnvironment(t)

	name := GetTestClusterName()
	resourceName := "ncloud_nks_cluster.cluster"

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
				Config: testAccResourceNcloudNKSClusterPublicNetworkConfig(name, TF_TEST_NKS_LOGIN_KEY, nksInfo),
				Check:  testAccResourceNcloudNKSClusterPublicNetworkConfigCheck(name, resourceName, nksInfo),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_public_network_KVM(t *testing.T) {
	validateAcctestEnvironment(t)

	name := GetTestClusterName()
	resourceName := "ncloud_nks_cluster.cluster"

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
				Config: testAccResourceNcloudNKSClusterPublicNetworkConfig(name, TF_TEST_NKS_LOGIN_KEY, nksInfo),
				Check:  testAccResourceNcloudNKSClusterPublicNetworkConfigCheck(name, resourceName, nksInfo),
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_Update_XEN(t *testing.T) {
	validateAcctestEnvironment(t)

	name := fmt.Sprintf("m3-%s", GetTestClusterName())
	resourceName := "ncloud_nks_cluster.cluster"

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
				Config:  testAccResourceNcloudNKSClusterDefaultConfig(name, TF_TEST_NKS_LOGIN_KEY, true, nksInfo),
				Check:   testAccResourceNcloudNKSClusterDefaultConfigCheck(resourceName, name, nksInfo),
				Destroy: false,
			},
			{
				Config:  testAccResourceNcloudNKSClusterUpdateConfig(name, TF_TEST_NKS_LOGIN_KEY, false, nksInfo),
				Check:   testAccResourceNcloudNKSClusterUpdateConfigCheck(resourceName, nksInfo),
				Destroy: false,
			},
		},
	})
}

func TestAccResourceNcloudNKSCluster_Update_KVM(t *testing.T) {
	validateAcctestEnvironment(t)

	name := fmt.Sprintf("m3-%s", GetTestClusterName())
	resourceName := "ncloud_nks_cluster.cluster"

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
				Config:  testAccResourceNcloudNKSClusterDefaultConfig(name, TF_TEST_NKS_LOGIN_KEY, true, nksInfo),
				Check:   testAccResourceNcloudNKSClusterDefaultConfigCheck(resourceName, name, nksInfo),
				Destroy: false,
			},
			{
				Config:  testAccResourceNcloudNKSClusterUpdateConfig(name, TF_TEST_NKS_LOGIN_KEY, false, nksInfo),
				Check:   testAccResourceNcloudNKSClusterUpdateConfigCheck(resourceName, nksInfo),
				Destroy: false,
			},
		},
	})
}

func testAccResourceNcloudNKSClusterDefaultConfig(name string, loginKeyName string, auditLog bool, nksInfo *NKSTestInfo) string {
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

	if !nksInfo.IsFin {
		b.WriteString(`
  ip_acl_default_action = "deny"
  ip_acl {
    action = "allow"
    address = "223.130.195.0/24"
    comment = "allow ip"
  }
  
`)
	}

	if nksInfo.needPublicLb {
		b.WriteString(fmt.Sprintf(`
  lb_public_subnet_no = %[1]s
`, *nksInfo.PublicLbSubnetList[0].SubnetNo))
	}

	b.WriteString(`
}
`)
	return b.String()
}

func testAccResourceNcloudNKSClusterUpdateConfig(name string, loginKeyName string, auditLog bool, nksInfo *NKSTestInfo) string {
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
    %[7]s, 
 	%[8]s
  ]
  vpc_no                      = %[9]s
  zone                        = "%[10]s-1"
  log {
    audit                     = %[11]t
  }
  oidc {
    issuer_url                = "https://keycloak.url/realms/update"
    client_id                 = "update-client"
  }
`, name, nksInfo.ClusterType, nksInfo.UpgradeK8sVersion, loginKeyName, *nksInfo.PrivateLbSubnetList[0].SubnetNo, nksInfo.HypervisorCode, *nksInfo.PrivateSubnetList[0].SubnetNo, *nksInfo.PrivateSubnetList[1].SubnetNo, *nksInfo.Vpc.VpcNo, nksInfo.Region, auditLog))
	if !nksInfo.IsFin {
		b.WriteString(`
  ip_acl_default_action = "allow"
  
`)
	}

	if nksInfo.needPublicLb {
		b.WriteString(fmt.Sprintf(`
  lb_public_subnet_no = %[1]s
`, *nksInfo.PublicLbSubnetList[0].SubnetNo))
	}

	b.WriteString(`
}
`)
	return b.String()
}

func testAccResourceNcloudNKSClusterPublicNetworkConfig(name string, loginKeyName string, nksInfo *NKSTestInfo) string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf(`
resource "ncloud_nks_cluster" "cluster" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = "%[3]s"
  login_key_name              = "%[4]s"
  hypervisor_code             = "%[5]s"
  lb_private_subnet_no        = %[6]s
  kube_network_plugin         = "cilium"
  public_network              = "true"
  subnet_no_list              = [
  %[7]s
  ]
  vpc_no                      = %[8]s
  zone                        = "%[9]s-1"
`, name, nksInfo.ClusterType, nksInfo.K8sVersion, loginKeyName, nksInfo.HypervisorCode, *nksInfo.PrivateLbSubnetList[0].SubnetNo, *nksInfo.PublicSubnetList[0].SubnetNo, *nksInfo.Vpc.VpcNo, nksInfo.Region))

	if nksInfo.needPublicLb {
		b.WriteString(fmt.Sprintf(`
  lb_public_subnet_no = %[1]s
`, *nksInfo.PublicLbSubnetList[0].SubnetNo))
	}

	b.WriteString(`
}
`)
	return b.String()
}

func testAccResourceNcloudNKSClusterDefaultConfigCheck(resourceName string, name string, nksInfo *NKSTestInfo) (check resource.TestCheckFunc) {
	var cluster vnks.Cluster
	check = resource.ComposeTestCheckFunc(
		testAccCheckNKSClusterExists(resourceName, &cluster),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "cluster_type", nksInfo.ClusterType),
		resource.TestMatchResourceAttr(resourceName, "k8s_version", regexp.MustCompile(nksInfo.K8sVersion)),
		resource.TestCheckResourceAttr(resourceName, "login_key_name", TF_TEST_NKS_LOGIN_KEY),
		resource.TestCheckResourceAttr(resourceName, "zone", fmt.Sprintf("%s-1", nksInfo.Region)),
		resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
		resource.TestCheckResourceAttr(resourceName, "log.0.audit", "true"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.issuer_url", "https://keycloak.url/realms/nks"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.client_id", "nks-client"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.username_claim", "preferred_username"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.username_prefix", "oidc:"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_claim", "groups"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.groups_prefix", "oidc:"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.required_claim", "iss=https://keycloak.url/realms/nks"),
	)
	if !nksInfo.IsFin {

		check = resource.ComposeTestCheckFunc(
			check,
			resource.TestCheckResourceAttr(resourceName, "ip_acl_default_action", "deny"),
			resource.TestCheckResourceAttr(resourceName, "ip_acl.0.action", "allow"),
			resource.TestCheckResourceAttr(resourceName, "ip_acl.0.address", "223.130.195.0/24"),
			resource.TestCheckResourceAttr(resourceName, "ip_acl.0.comment", "allow ip"),
		)
	}
	return
}

func testAccResourceNcloudNKSClusterUpdateConfigCheck(resourceName string, nksInfo *NKSTestInfo) (check resource.TestCheckFunc) {
	var cluster vnks.Cluster
	return resource.ComposeTestCheckFunc(
		testAccCheckNKSClusterExists(resourceName, &cluster),
		resource.TestMatchResourceAttr(resourceName, "k8s_version", regexp.MustCompile(nksInfo.UpgradeK8sVersion)),
		resource.TestCheckResourceAttr(resourceName, "log.0.audit", "false"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.issuer_url", "https://keycloak.url/realms/update"),
		resource.TestCheckResourceAttr(resourceName, "oidc.0.client_id", "update-client"),
		resource.TestCheckResourceAttr(resourceName, "ip_acl_default_action", "allow"),
		resource.TestCheckResourceAttr(resourceName, "ip_acl.#", "0"),
	)
}

func testAccResourceNcloudNKSClusterPublicNetworkConfigCheck(name string, resourceName string, nksInfo *NKSTestInfo) (check resource.TestCheckFunc) {
	var cluster vnks.Cluster
	return resource.ComposeTestCheckFunc(
		testAccCheckNKSClusterExists(resourceName, &cluster),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "cluster_type", nksInfo.ClusterType),
		resource.TestMatchResourceAttr(resourceName, "k8s_version", regexp.MustCompile(nksInfo.K8sVersion)),
		resource.TestCheckResourceAttr(resourceName, "login_key_name", TF_TEST_NKS_LOGIN_KEY),
		resource.TestCheckResourceAttr(resourceName, "public_network", "true"),
		resource.TestCheckResourceAttr(resourceName, "zone", fmt.Sprintf("%s-1", nksInfo.Region)),
		resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
		resource.TestCheckResourceAttr(resourceName, "ip_acl_default_action", "allow"),
		resource.TestCheckResourceAttr(resourceName, "ip_acl.#", "0"),
	)
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

func getNKSTestInfo(hypervisor string) (*NKSTestInfo, error) {

	nksInfo := &NKSTestInfo{
		Region:            os.Getenv("NCLOUD_REGION"),
		K8sVersion:        "1.27.9",
		UpgradeK8sVersion: "1.27.9",
		HypervisorCode:    hypervisor,
	}
	zoneCode := ncloud.String(fmt.Sprintf("%s-1", nksInfo.Region))
	nksInfo.IsFin = strings.HasPrefix(nksInfo.Region, "F")
	nksInfo.IsCaaS = nksInfo.Region[1:2] == "CS"
	nksInfo.needPublicLb = true

	if nksInfo.IsFin {
		nksInfo.needPublicLb = false
	}
	if hypervisor == "KVM" {
		switch nksInfo.Region {
		case "FKR":
			nksInfo.ClusterType = "SVR.VNKS.STAND.C002.M008.G003"
		case "PCS02":
			nksInfo.ClusterType = "SVR.VNKS.STAND.C008.M032.G003"
		default:
			nksInfo.ClusterType = "SVR.VNKS.STAND.C002.M008.G003"
		}
		nksInfo.K8sVersion = fmt.Sprintf("%s-nks.2", nksInfo.K8sVersion)
		nksInfo.UpgradeK8sVersion = fmt.Sprintf("%s-nks.2", nksInfo.UpgradeK8sVersion)
	} else {
		switch nksInfo.Region {
		case "FKR":
			nksInfo.ClusterType = "SVR.VNKS.STAND.C002.M008.NET.HDD.B050.G001"
		default:
			nksInfo.ClusterType = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
		}
		nksInfo.K8sVersion = fmt.Sprintf("%s-nks.1", nksInfo.K8sVersion)
		nksInfo.UpgradeK8sVersion = fmt.Sprintf("%s-nks.1", nksInfo.UpgradeK8sVersion)
	}

	vpcName := ncloud.String("tf-test-vpc")
	apiKeys := &ncloud.APIKey{
		AccessKey: os.Getenv("NCLOUD_ACCESS_KEY"),
		SecretKey: os.Getenv("NCLOUD_SECRET_KEY"),
	}

	vpcClient := vpc.NewAPIClient(vpc.NewConfiguration(apiKeys))

	reqParams := &vpc.GetVpcListRequest{
		RegionCode: &nksInfo.Region,
		VpcName:    vpcName,
	}
	vpcResp, err := vpcClient.V2Api.GetVpcList(reqParams)
	if err != nil {
		return nil, err
	}

	if len(vpcResp.VpcList) == 0 {
		createVpcReq := &vpc.CreateVpcRequest{
			RegionCode:    &nksInfo.Region,
			VpcName:       vpcName,
			Ipv4CidrBlock: ncloud.String("10.0.0.0/16"),
		}
		createVpcResp, err := vpcClient.V2Api.CreateVpc(createVpcReq)
		if err != nil {
			return nil, err
		}
		nksInfo.Vpc = createVpcResp.VpcList[0]

	} else {
		nksInfo.Vpc = vpcResp.VpcList[0]
	}

	aclReq := &vpc.GetNetworkAclListRequest{
		RegionCode: &nksInfo.Region,
		VpcNo:      nksInfo.Vpc.VpcNo,
	}
	aclResp, err := vpcClient.V2Api.GetNetworkAclList(aclReq)
	if err != nil {
		return nil, err
	}
	for _, acl := range aclResp.NetworkAclList {
		if *acl.IsDefault {
			nksInfo.DefaultAcl = acl
		}
	}

	subnetReqParams := &vpc.GetSubnetListRequest{
		VpcNo:      nksInfo.Vpc.VpcNo,
		RegionCode: &nksInfo.Region,
	}

	subnetResp, err := vpcClient.V2Api.GetSubnetList(subnetReqParams)
	if err != nil {
		return nil, err
	}

	for _, subnet := range subnetResp.SubnetList {
		if *subnet.UsageType.Code == "GEN" && *subnet.SubnetType.Code == "PRIVATE" {
			nksInfo.PrivateSubnetList = append(nksInfo.PrivateSubnetList, subnet)
		} else if *subnet.UsageType.Code == "GEN" && *subnet.SubnetType.Code == "PUBLIC" {
			nksInfo.PublicSubnetList = append(nksInfo.PublicSubnetList, subnet)
		} else if *subnet.UsageType.Code == "LOADB" && *subnet.SubnetType.Code == "PRIVATE" {
			nksInfo.PrivateLbSubnetList = append(nksInfo.PrivateLbSubnetList, subnet)
		} else if *subnet.UsageType.Code == "LOADB" && *subnet.SubnetType.Code == "PUBLIC" {
			nksInfo.PublicLbSubnetList = append(nksInfo.PublicLbSubnetList, subnet)
		}
	}

	if len(nksInfo.PrivateSubnetList) == 0 {
		for i := 1; i <= 2; i++ {

			createSubnetReq := &vpc.CreateSubnetRequest{
				VpcNo:          nksInfo.Vpc.VpcNo,
				RegionCode:     &nksInfo.Region,
				ZoneCode:       zoneCode,
				SubnetTypeCode: ncloud.String("PRIVATE"),
				UsageTypeCode:  ncloud.String("GEN"),
				NetworkAclNo:   nksInfo.DefaultAcl.NetworkAclNo,
				Subnet:         ncloud.String(fmt.Sprintf("10.0.%d.0/24", i)),
				SubnetName:     ncloud.String(fmt.Sprintf("tf-subnet-priv-%d", i)),
			}

			subnetResp, err := vpcClient.V2Api.CreateSubnet(createSubnetReq)
			if err != nil {
				return nil, err
			}

			nksInfo.PrivateSubnetList = append(nksInfo.PrivateSubnetList, subnetResp.SubnetList[0])
		}
	}

	if len(nksInfo.PublicSubnetList) == 0 && !nksInfo.IsCaaS {
		createSubnetReq := &vpc.CreateSubnetRequest{
			VpcNo:          nksInfo.Vpc.VpcNo,
			RegionCode:     &nksInfo.Region,
			ZoneCode:       zoneCode,
			SubnetTypeCode: ncloud.String("PUBLIC"),
			UsageTypeCode:  ncloud.String("GEN"),
			NetworkAclNo:   nksInfo.DefaultAcl.NetworkAclNo,
			Subnet:         ncloud.String("10.0.10.0/24"),
			SubnetName:     ncloud.String("tf-subnet-pub"),
		}

		subnetResp, err := vpcClient.V2Api.CreateSubnet(createSubnetReq)
		if err != nil {
			return nil, err
		}

		nksInfo.PublicSubnetList = append(nksInfo.PublicSubnetList, subnetResp.SubnetList[0])
	}

	if len(nksInfo.PrivateLbSubnetList) == 0 {
		createSubnetReq := &vpc.CreateSubnetRequest{
			VpcNo:          nksInfo.Vpc.VpcNo,
			RegionCode:     &nksInfo.Region,
			ZoneCode:       zoneCode,
			SubnetTypeCode: ncloud.String("PRIVATE"),
			UsageTypeCode:  ncloud.String("LOADB"),
			NetworkAclNo:   nksInfo.DefaultAcl.NetworkAclNo,
			Subnet:         ncloud.String("10.0.100.0/24"),
			SubnetName:     ncloud.String("tf-subnet-lb-priv"),
		}

		subnetResp, err := vpcClient.V2Api.CreateSubnet(createSubnetReq)
		if err != nil {
			return nil, err
		}

		nksInfo.PrivateLbSubnetList = append(nksInfo.PrivateLbSubnetList, subnetResp.SubnetList[0])
	}

	if len(nksInfo.PublicLbSubnetList) == 0 && nksInfo.needPublicLb {
		createSubnetReq := &vpc.CreateSubnetRequest{
			VpcNo:          nksInfo.Vpc.VpcNo,
			RegionCode:     &nksInfo.Region,
			ZoneCode:       zoneCode,
			SubnetTypeCode: ncloud.String("PUBLIC"),
			UsageTypeCode:  ncloud.String("LOADB"),
			NetworkAclNo:   nksInfo.DefaultAcl.NetworkAclNo,
			Subnet:         ncloud.String("10.0.101.0/24"),
			SubnetName:     ncloud.String("tf-subnet-lb-pub"),
		}

		subnetResp, err := vpcClient.V2Api.CreateSubnet(createSubnetReq)
		if err != nil {
			return nil, err
		}

		nksInfo.PublicLbSubnetList = append(nksInfo.PublicLbSubnetList, subnetResp.SubnetList[0])
	}

	return nksInfo, nil

}
