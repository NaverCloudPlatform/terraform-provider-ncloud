package hadoop_test

import (
	"fmt"
	randacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoops_basic(t *testing.T) {
	dataName := "data.ncloud_hadoops.hadoops"
	filteredDataName := "data.ncloud_hadoops.hadoops_by_filter"
	resourceName := "ncloud_hadoop.hadoop"
	instanceName := fmt.Sprintf("tf-hadoop-%s", randacctest.RandString(3))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHadoopsConfig(instanceName),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "hadoops.0.cluster_name", resourceName, "cluster_name"),
					resource.TestCheckResourceAttrPair(dataName, "hadoops.0.cluster_type_code", resourceName, "cluster_type_code"),
					resource.TestCheckResourceAttrPair(filteredDataName, "hadoops_by_filter.0.cluster_name", resourceName, "cluster_name"),
					resource.TestCheckResourceAttrPair(filteredDataName, "hadoops_by_filter.0.cluster_type_code", resourceName, "cluster_type_code"),
				),
			},
		},
	})
}

func testAccDataSourceHadoopsConfig(name string) string {
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

data "ncloud_hadoops" "hadoops" {
	cluster_name = "%[1]s"
}

data "ncloud_hadoops" "hadoops_by_filter" {
	cluster_name = "%[1]s"
}
`, name)
}
