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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckHadoopDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHadoopConfig(testHadoopName),
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
				Config: testAccHadoopConfigUpdate(testHadoopName, imageProductCode, workerCountBefore, productCodeBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHadoopExistsWithProvider(resourceName, &hadoopInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "worker_node_count", strconv.FormatInt(int64(workerCountBefore), 10)),
					resource.TestCheckResourceAttr(resourceName, "master_node_product_code", productCodeBefore),
					resource.TestCheckResourceAttr(resourceName, "worker_node_product_code", productCodeBefore),
				),
			},
			{
				Config: testAccHadoopConfigUpdate(testHadoopName, imageProductCode, workerCountAfter, productCodeAfter),
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

func testAccHadoopConfig(testName string) string {
	return fmt.Sprintf(`
resource "ncloud_hadoop" "hadoop" {
	vpc_no = "49956"
	cluster_name = "%[1]s"
	cluster_type_code = "CORE_HADOOP_WITH_SPARK"
	admin_user_name = "admin-test"
	admin_user_password = "Admin!2Admin"
	login_key_name = "naverCloud"
	master_node_subnet_no = "111983"
	edge_node_subnet_no = "111985"
	worker_node_subnet_no = "111984"
	bucket_name = "ddd1"
	master_node_data_storage_type = "SSD"
	worker_node_data_storage_type = "SSD"
	master_node_data_storage_size = 100
	worker_node_data_storage_size = 100
}
`, testName)
}

func testAccHadoopConfigUpdate(testName, imageProduct string, workerCount int, productCode string) string {
	return fmt.Sprintf(`
resource "ncloud_hadoop" "hadoop" {
	vpc_no = "49956"
	cluster_name = "%[1]s"
	cluster_type_code = "CORE_HADOOP_WITH_SPARK"
	admin_user_name = "admin-test"
	admin_user_password = "Admin!2Admin"
	login_key_name = "naverCloud"
	master_node_subnet_no = "111983"
	edge_node_subnet_no = "111985"
	worker_node_subnet_no = "111984"
	bucket_name = "ddd1"
	master_node_data_storage_type = "SSD"
	worker_node_data_storage_type = "SSD"
	master_node_data_storage_size = 100
	worker_node_data_storage_size = 100
	
	image_product_code = "%[2]s"
	worker_node_count = %[3]d
	master_node_product_code = "%[4]s"
	worker_node_product_code = "%[4]s"
}
`, testName, imageProduct, workerCount, productCode)
}
