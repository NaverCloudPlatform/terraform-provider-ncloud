package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudServerProducts_classic_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		IsUnitTest:               false,
		ProtoV5ProviderFactories: ClassicProtoV5ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductsConfig("SPSW0LINUX000045"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_products.all"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudServerProducts_vpc_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		IsUnitTest:               false,
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudServerProductsConfig("SW.VSVR.OS.LNX64.CNTOS.0703.B050"),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_server_products.all"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudServerProductsConfig(productCode string) string {
	return fmt.Sprintf(`
data "ncloud_server_products" "all" {
	server_image_product_code = "%s"

	filter {
		name = "product_code"
		values = ["SSD"]
		regex = true
	}

	filter {
		name = "cpu_count"
		values = ["2"]
	}

	filter {
		name = "memory_size"
		values = ["8GB"]
	}

	filter {
		name = "base_block_storage_size"
		values = ["50GB"]
	}

	filter {
		name = "product_type"
		values = ["STAND"]
	}
}
`, productCode)
}
