package ncloud

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vses2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

const TF_TEST_SES_LOGIN_KEY = "tf-ses-login-key"

func TestAccResourceNcloudSESCluster_basic(t *testing.T) {
	var cluster vses2.OpenApiGetClusterInfoResponseVo
	resourceName := "ncloud_ses_cluster.cluster"
	testClusterName := getTestClusterName()
	searchEngineVersionCode := "133"
	region := os.Getenv("NCLOUD_REGION")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSESClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSESClusterConfig(testClusterName, TF_TEST_SES_LOGIN_KEY, searchEngineVersionCode, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSESClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", testClusterName),
				),
			},
		},
	})
}

func testAccResourceSESClusterConfig(testClusterName string, loginKey string, version string, region string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "172.16.0.0/16"
}

resource "ncloud_subnet" "node_subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s"
	subnet             = "172.16.1.0/24"
	zone               = "%[4]s-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}
data "ncloud_ses_versions" "version" {
}

data "ncloud_ses_node_os_images" "os_images" {
}

data "ncloud_ses_node_products" "product_codes" {
  os_image_code = data.ncloud_ses_node_os_images.os_images.images.0.id
  subnet_no = ncloud_subnet.node_subnet.id
}

resource "ncloud_login_key" "loginkey" {
  key_name = "%[2]s"
}

resource "ncloud_ses_cluster" "cluster" {
  cluster_name                  = "%[1]s"
  os_image_code         		= data.ncloud_ses_node_os_images.os_images.images.0.id
  vpc_no                        = ncloud_vpc.vpc.id
  search_engine {
	  version_code    			= "%[3]s"
	  user_name       			= "admin"
	  user_password   			= "qwe123!@#"
      dashboard_port            = "5601"
  }
  manager_node {  
	  is_dual_manager           = false
	  product_code     			= data.ncloud_ses_node_products.product_codes.codes.0.id
	  subnet_no        			= ncloud_subnet.node_subnet.id
  }
  data_node {
	  product_code       		= data.ncloud_ses_node_products.product_codes.codes.0.id
	  subnet_no           		= ncloud_subnet.node_subnet.id
	  count            		    = 3
	  storage_size        		= 100
  }
  login_key_name                = ncloud_login_key.loginkey.key_name
}
`, testClusterName, loginKey, version, region)
}

func testAccCheckSESClusterExists(n string, cluster *vses2.OpenApiGetClusterInfoResponseVo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster service_group_instance_no is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		resp, err := getSESCluster(context.Background(), config, rs.Primary.ID)
		if err != nil {
			return err
		}

		cluster = resp

		return nil
	}
}

func testAccCheckSESClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_ses_cluster" {
			continue
		}

		cluster, err := getSESCluster(context.Background(), config, rs.Primary.ID)
		if err != nil {
			return err
		}

		if *cluster.ClusterStatus != "return" {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}
