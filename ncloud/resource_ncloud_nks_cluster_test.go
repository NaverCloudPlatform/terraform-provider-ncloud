package ncloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ncloud/sdk/vnks"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudNKSCluster_basic(t *testing.T) {
	var cluster vnks.Cluster
	name := fmt.Sprintf("test-nksbasic-%s", acctest.RandString(5))
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.20.11"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, productCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", clusterType),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", k8sVersion),
					resource.TestCheckResourceAttr(resourceName, "product_code", productCode),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
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

func TestAccResourceNcloudNKSCluster_disappears(t *testing.T) {
	var cluster vnks.Cluster
	name := fmt.Sprintf("test-nksdisappears-%s", acctest.RandString(5))
	clusterType := "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
	k8sVersion := "1.20.11"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
	resourceName := "ncloud_nks_cluster.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNKSClusterConfig(name, clusterType, k8sVersion, productCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSClusterExists(resourceName, &cluster),
					testAccCheckNKSClusterDisappears(&cluster),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

//func TestAccResourceNcloudNKSCluster_updateName(t *testing.T) {
//	var subnet vpc.Subnet
//	name := fmt.Sprintf("test-nksname-%s", acctest.RandString(5))
//	cidr := "10.2.2.0/24"
//	resourceName := "ncloud_nks_cluster.bar"
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckNKSClusterDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: TestAccResourceNcloudNKSClusterConfig(name),
//
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNKSClusterExists(resourceName, &subnet),
//				),
//			},
//			{
//				Config: TestAccResourceNcloudNKSClusterConfig("testacc-subnet-update"),
//
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNKSClusterExists(resourceName, &subnet),
//				),
//				ExpectError: regexp.MustCompile("Change 'name' is not support, Please set `name` as a old value"),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}
//
//func TestAccResourceNcloudNKSCluster_updateNetworkACL(t *testing.T) {
//	var subnet vpc.Subnet
//	name := fmt.Sprintf("test-nksupdate-nacl-%s", acctest.RandString(5))
//	cidr := "10.2.2.0/24"
//	resourceName := "ncloud_nks_cluster.bar"
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckNKSClusterDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: TestAccResourceNcloudNKSClusterConfig(name),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNKSClusterExists(resourceName, &subnet),
//				),
//			},
//			{
//				Config: TestAccResourceNcloudNKSClusterConfigUpdateNetworkACL(name, cidr),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNKSClusterExists(resourceName, &subnet),
//				),
//			},
//		},
//	})
//}
//
//func TestAccResourceNcloudNKSCluster_InvalidCIDR(t *testing.T) {
//	name := fmt.Sprintf("test-nksupdate-nacl-%s", acctest.RandString(5))
//	cidr := "10.3.2.0/24"
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckNKSClusterDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config:      TestAccResourceNcloudNKSClusterConfigInvalidCIDR(name, cidr),
//				ExpectError: regexp.MustCompile("The subnet must belong to the IPv4 CIDR of the specified VPC."),
//			},
//		},
//	})
//}

func testAccResourceNcloudNKSClusterConfig(name string, clusterType string, k8sVersion string, productCode string) string {
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

resource "ncloud_nks_cluster" "test" {
  name                        = "%[1]s"
  cluster_type                = "%[2]s"
  k8s_version                 = "%[3]s"
  login_key_name              = ncloud_login_key.loginkey.key_name
  subnet_lb_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
  ]
  vpc_no                      = ncloud_vpc.vpc.vpc_no
  zone_no                     = "2"

  node_pool {
    is_default     = true
    name           = "%[1]s"
    node_count     = 1
    product_code   = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
    subnet_no_list              = [
    ncloud_subnet.subnet1.id,
    ncloud_subnet.subnet2.id,
    ]
  }
}
`, name, clusterType, k8sVersion, productCode)
}

//
//func TestAccResourceNcloudNKSClusterConfigUpdateNetworkACL(name, cidr string) string {
//	return fmt.Sprintf(`
//resource "ncloud_vpc" "foo" {
//	name               = "%[1]s"
//	ipv4_cidr_block    = "10.2.0.0/16"
//}
//
//resource "ncloud_network_acl" "nacl" {
//	vpc_no      = ncloud_vpc.foo.vpc_no
//	name        = "%[1]s"
//	description = "for test acc"
//}
//
//resource "ncloud_nks_cluster" "bar" {
//	vpc_no             = ncloud_vpc.foo.vpc_no
//	name               = "%[1]s"
//	subnet             = "%[2]s"
//	zone               = "KR-1"
//	network_acl_no     = ncloud_network_acl.nacl.network_acl_no
//	subnet_type        = "PUBLIC"
//	usage_type         = "GEN"
//}
//`, name, cidr)
//}
//
//func TestAccResourceNcloudNKSClusterConfigInvalidCIDR(name, cidr string) string {
//	return fmt.Sprintf(`
//resource "ncloud_vpc" "foo" {
//	name               = "%[1]s"
//	ipv4_cidr_block    = "10.2.0.0/16"
//}
//
//resource "ncloud_nks_cluster" "bar" {
//	vpc_no             = ncloud_vpc.foo.vpc_no
//	name               = "%[1]s"
//	subnet             = "%s"
//	zone               = "KR-1"
//	network_acl_no     = ncloud_vpc.foo.default_network_acl_no
//	subnet_type        = "PUBLIC"
//	usage_type         = "GEN"
//}
//`, name, cidr)
//}

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
		resp, err := getNKSClusterCluster(context.Background(), config, rs.Primary.ID)
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

		cluster, err := getNKSClusterCluster(context.Background(), config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if cluster != nil {
			return errors.New("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckNKSClusterDisappears(cluster *vnks.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		err := config.Client.vnks.V2Api.ClustersUuidDelete(context.Background(), cluster.Uuid)

		d := &schema.ResourceData{}
		d.SetId(ncloud.StringValue(cluster.Uuid))
		if err := waitForNKSClusterDeletion(context.Background(), d, config); err != nil {
			return err
		}

		return err
	}
}
