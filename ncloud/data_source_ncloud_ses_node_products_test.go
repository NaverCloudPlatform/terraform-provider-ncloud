package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceNcloudSESNodeProductCodes(t *testing.T) {
	dataName := "data.ncloud_ses_node_products.product_codes"
	region := os.Getenv("NCLOUD_REGION")
	testName := "ses-product-test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: getTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSESNodeProductConfig(testName, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSESNodeProductConfig(testName string, region string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name               = "%[1]s"
	ipv4_cidr_block    = "192.168.0.0/16"
}

resource "ncloud_subnet" "node_subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "%[1]s"
	subnet             = "192.168.1.0/24"
	zone               = "%[2]s-1"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

data "ncloud_ses_node_os_images" "os_images" {
}

data "ncloud_ses_node_products" "product_codes" {
  os_image_code = data.ncloud_ses_node_os_images.os_images.images.0.id
  subnet_no = ncloud_subnet.node_subnet.id
}
`, testName, region)
}
