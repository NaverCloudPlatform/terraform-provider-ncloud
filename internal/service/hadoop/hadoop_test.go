package hadoop_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vhadoop"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	hadoopService "github.com/terraform-providers/terraform-provider-ncloud/internal/service/hadoop"
	"strconv"
	"strings"
	"testing"
)

func TestAccResourceNcloudHadoop_basic(t *testing.T) {
	var hadoopInstance vhadoop.CloudHadoopInstance
	testHadoopName := fmt.Sprintf("tf-hadoop-%s", acctest.RandString(4))
	resourceName := "ncloud_hadoop.hadoop"
	bucketName := "akj1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckHadoopDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHadoopConfig(testHadoopName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHadoopExistsWithProvider(resourceName, &hadoopInstance, GetTestProvider(true)),
				),
			},
		},
	})
}

func TestAccResourceNcloudHadoop_update(t *testing.T) {
	var hadoopInstance vhadoop.CloudHadoopInstance
	testHadoopName := fmt.Sprintf("tf-hadoop-%s", acctest.RandString(3))
	resourceName := "ncloud_hadoop.hadoop"
	imageProductCode := "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
	bucketName := "akj1"

	//masterProductBefore := "SVR.VCHDP.MSTDT.STAND.C004.M016.NET.HDD.B050.G002"
	//edgeProductBefore := "SVR.VCHDP.EDGND.STAND.C004.M016.NET.HDD.B050.G002"
	//workerProductBefore := "SVR.VCHDP.MSTDT.STAND.C004.M016.NET.HDD.B050.G002"
	//
	//masterProductAfter := "SVR.VCHDP.MSTDT.HICPU.C008.M016.NET.HDD.B050.G002"
	//edgeProductAfter := "SVR.VCHDP.EDGND.HICPU.C008.M016.NET.HDD.B050.G001"
	//workerProductAfter := "SVR.VCHDP.MSTDT.HICPU.C008.M016.NET.HDD.B050.G002"

	workerCountBefore := 2
	workerCountAfter := 3
	productCodeBefore := "SVR.VCHDP.MSTDT.STAND.C004.M016.NET.HDD.B050.G002"
	productCodeAfter := "SVR.VCHDP.MSTDT.HICPU.C008.M016.NET.HDD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckHadoopDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHadoopConfigUpdate(testHadoopName, imageProductCode, workerCountBefore, productCodeBefore, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHadoopExistsWithProvider(resourceName, &hadoopInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "worker_node_count", strconv.FormatInt(int64(workerCountBefore), 10)),
					resource.TestCheckResourceAttr(resourceName, "master_node_product_code", productCodeBefore),
					resource.TestCheckResourceAttr(resourceName, "worker_node_product_code", productCodeBefore),
				),
			},
			{
				Config: testAccHadoopConfigUpdate(testHadoopName, imageProductCode, workerCountAfter, productCodeAfter, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHadoopExistsWithProvider(resourceName, &hadoopInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "worker_node_count", strconv.FormatInt(int64(workerCountAfter), 10)),
					resource.TestCheckResourceAttr(resourceName, "master_node_product_code", productCodeAfter),
					resource.TestCheckResourceAttr(resourceName, "worker_node_product_code", productCodeAfter),
				),
			},
		},
	})
}

func testAccCheckHadoopExistsWithProvider(n string, hadoop *vhadoop.CloudHadoopInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		hadoopInstance, err := hadoopService.GetHadoopInstance(context.Background(), config, resource.Primary.ID)
		if err != nil {
			return err
		}

		if hadoopInstance != nil {
			*hadoop = *hadoopInstance
			return nil
		}

		return fmt.Errorf("hadoop instance not found")
	}
}

func testAccCheckHadoopDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_hadoop" {
			continue
		}
		instance, err := hadoopService.GetHadoopInstance(context.Background(), config, rs.Primary.ID)
		if err != nil && !checkNoInstanceResponse(err) {
			return err
		}

		if instance != nil {
			return errors.New("hadoop still exists")
		}
	}

	return nil
}

func checkNoInstanceResponse(err error) bool {
	return strings.Contains(err.Error(), "5001017")
}

func testAccHadoopConfig(testName, bucketName string) string {

	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "edge_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-edge"
	subnet             = "10.5.0.0/18"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "master_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-master"
	subnet             = "10.5.64.0/19"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "worker_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-worker"
	subnet             = "10.5.96.0/20"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_hadoop" "hadoop" {
	vpc_no = ncloud_vpc.test.vpc_no
	cluster_name = "%[1]s"
	cluster_type_code = "CORE_HADOOP_WITH_SPARK"
	admin_user_name = "admin-test"
	admin_user_password = "Admin!2Admin"
	login_key_name = "naverCloud"
	master_node_subnet_no = ncloud_subnet.master_subnet.subnet_no
	edge_node_subnet_no = ncloud_subnet.edge_subnet.subnet_no
	worker_node_subnet_no = ncloud_subnet.worker_subnet.subnet_no
	bucket_name = "%[2]s"
	master_node_data_storage_type = "SSD"
	worker_node_data_storage_type = "SSD"
	master_node_data_storage_size = 100
	worker_node_data_storage_size = 100
}
`, testName, bucketName)
}

func testAccHadoopConfigUpdate(testName, imageProduct string, workerCount int, productCode, bucketName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "edge_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-edge"
	subnet             = "10.5.0.0/18"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "master_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-master"
	subnet             = "10.5.64.0/19"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "worker_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-worker"
	subnet             = "10.5.96.0/20"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_hadoop" "hadoop" {
	vpc_no = ncloud_vpc.test.vpc_no
	cluster_name = "%[1]s"
	cluster_type_code = "CORE_HADOOP_WITH_SPARK"
	admin_user_name = "admin-test"
	admin_user_password = "Admin!2Admin"
	login_key_name = "naverCloud"
	master_node_subnet_no =  ncloud_subnet.master_subnet.subnet_no
	edge_node_subnet_no = ncloud_subnet.edge_subnet.subnet_no
	worker_node_subnet_no =  ncloud_subnet.worker_subnet.subnet_no
	bucket_name = "%[5]s"
	master_node_data_storage_type = "SSD"
	worker_node_data_storage_type = "SSD"
	master_node_data_storage_size = 100
	worker_node_data_storage_size = 100
	
	image_product_code = "%[2]s"
	worker_node_count = %[3]d
	master_node_product_code = "%[4]s"
	worker_node_product_code = "%[4]s"
}
`, testName, imageProduct, workerCount, productCode, bucketName)
}
