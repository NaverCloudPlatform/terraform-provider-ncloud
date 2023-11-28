package hadoop_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoopAddOnbasic(t *testing.T) {
	dataName := "data.ncloud_hadoop_add_on.addon"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAddOnConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "add_on_list.0", "Presto (0.240)"),
					resource.TestCheckResourceAttr(dataName, "add_on_list.1", "HBASE (2.0.2)"),
				),
			},
		},
	})
}

func testAccDataSourceAddOnConfig() string {
	return fmt.Sprintf(`
data "ncloud_hadoop_add_on" "addon" {
	image_product_code= "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
	cluster_type_code= "CORE_HADOOP_WITH_SPARK"
}
`)
}
