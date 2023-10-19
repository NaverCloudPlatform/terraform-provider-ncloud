package server_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccResourceNcloudPublicIpInstance_classic_basic(t *testing.T) {
	var instance *server.PublicIpInstance
	description := fmt.Sprintf("test-public-ip-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckPublicIpInstanceDestroy(s, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpInstanceClassicConfig(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(false)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "public_ip_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "public_ip", regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"server_instance_no"},
			},
		},
	})
}

func TestAccResourceNcloudPublicIpInstance_vpc_basic(t *testing.T) {
	var instance *server.PublicIpInstance

	name := fmt.Sprintf("test-public-ip-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckPublicIpInstanceDestroy(s, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpInstanceVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "public_ip_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "public_ip", regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)),
					resource.TestCheckResourceAttr(resourceName, "description", name),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"server_instance_no"},
			},
		},
	})
}

func TestAccResourceNcloudPublicIpInstance_classic_updateServerInstanceNo(t *testing.T) {
	var instance *server.PublicIpInstance
	serverNameFoo := fmt.Sprintf("test-public-ip-foo-%s", acctest.RandString(5))
	serverNameBar := fmt.Sprintf("test-public-ip-bar-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckPublicIpInstanceDestroy(s, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, ""),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(false))),
			},
			{
				Config: testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, "${ncloud_server.foo.id}"),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(false))),
			},
			{
				Config: testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, "${ncloud_server.bar.id}"),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(false))),
			},
			{
				Config: testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, ""),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(false))),
			},
		},
	})
}

func TestAccResourceNcloudPublicIpInstance_vpc_updateServerInstanceNo(t *testing.T) {
	var instance *server.PublicIpInstance
	serverNameFoo := fmt.Sprintf("test-public-ip-foo-%s", acctest.RandString(5))
	serverNameBar := fmt.Sprintf("test-public-ip-bar-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckPublicIpInstanceDestroy(s, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, ""),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(true))),
			},
			{
				Config: testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, "${ncloud_server.foo.id}"),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(true))),
			},
			{
				Config: testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, "${ncloud_server.bar.id}"),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(true))),
			},
			{
				Config: testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, ""),
				Check:  resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance, GetTestProvider(true))),
			},
		},
	})
}

func testAccCheckPublicIpInstanceExists(n string, i *server.PublicIpInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)

		instance, err := server.GetPublicIp(config, rs.Primary.ID)

		if err != nil {
			return nil
		}

		if instance != nil {
			i = instance

			return nil
		}

		return fmt.Errorf("public ip instance not found")
	}
}

func testAccCheckPublicIpInstanceDestroy(s *terraform.State, provider *schema.Provider) error {
	for _, rs := range s.RootModule().Resources {
		config := provider.Meta().(*conn.ProviderConfig)

		if rs.Type != "ncloud_public_ip" {
			continue
		}

		instance, err := server.GetPublicIp(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("Public IP still exists:\n%#v", instance)
		}

		break
	}

	return nil
}

func testAccPublicIpInstanceClassicConfig(description string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "server" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000045"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_public_ip" "public_ip" {
	description = "%[1]s"
	depends_on = [ncloud_server.server]
}
`, description)
}

func testAccPublicIpInstanceVpcConfig(name string) string {
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
	description = "%[1]s"
}`, name)
}

func testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, serverInstanceNo string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "foo" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000045"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_server" "bar" {
	name = "%[2]s"
	server_image_product_code = "SPSW0LINUX000045"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_public_ip" "public_ip" {
	description = "%[1]s"
	server_instance_no = "%[3]s"
	depends_on = [ncloud_server.foo]
}
`, serverNameFoo, serverNameBar, serverInstanceNo)
}

func testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, serverInstanceNo string) string {
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

resource "ncloud_server" "foo" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_server" "bar" {
	subnet_no = ncloud_subnet.test.id
	name = "%[2]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_public_ip" "public_ip" {
	description = "%[1]s"
	server_instance_no = "%[3]s"
	depends_on  = [ncloud_server.foo]
}
`, serverNameFoo, serverNameBar, serverInstanceNo)
}
