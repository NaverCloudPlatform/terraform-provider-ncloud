package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudPublicIp_classic_basic(t *testing.T) {
	resourceName := "ncloud_public_ip.public_ip"
	dataName := "data.ncloud_public_ip.test"
	name := fmt.Sprintf("tf-public-ip-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "server_name", resourceName, "server_name"),

					// Classic only
					resource.TestCheckResourceAttrPair(dataName, "zone", resourceName, "zone"),
					resource.TestCheckResourceAttrPair(dataName, "kind_type", resourceName, "kind_type"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIp_vpc_basic(t *testing.T) {
	resourceName := "ncloud_public_ip.public_ip"
	dataName := "data.ncloud_public_ip.test"
	name := fmt.Sprintf("tf-public-ip-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "public_ip", resourceName, "public_ip"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "server_name", resourceName, "server_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIpIsAssociated(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpAssociatedConfig,
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				ExpectError: regexp.MustCompile("no results"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_public_ip.test"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudPublicIpSearch(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudPublicIpSearchConfig,
				SkipFunc: func() (bool, error) {
					return skipNoResultsTest, nil
				},
				ExpectError: regexp.MustCompile("no results"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_public_ip.test"),
					resource.TestCheckResourceAttrSet(
						"data.ncloud_public_ip.test",
						"server_instance.server_instance_no",
					),
				),
			},
		},
	})
}

func testAccDataSourceNcloudPublicIpClassicConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "server" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_public_ip" "public_ip" {
	description = "%[1]s"
	depends_on = [ncloud_server.server]
}

data "ncloud_public_ip" "test" {
	id = ncloud_public_ip.public_ip.id
}
`, name)
}

func testAccDataSourceNcloudPublicIpVpcConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_public_ip" "public_ip" {
	depends_on = [ncloud_server.server]
}

data "ncloud_public_ip" "test" {
	id = ncloud_public_ip.public_ip.id
}
`, name)
}

var testAccDataSourceNcloudPublicIpAssociatedConfig = `
data "ncloud_public_ip" "test" {
	is_associated = "false"
}
`

var testAccDataSourceNcloudPublicIpSearchConfig = `
data "ncloud_public_ip" "test" {
	filter {
		name = "server_instance.server_name"
		values = ["tf-2807-1"]
	}
}
`
