package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudPublicIpInstance_basic(t *testing.T) {
	instance := map[string]interface{}{}

	description := fmt.Sprintf("test-public-ip-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	testCheckAttribute := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			config := testAccProvider.Meta().(*ProviderConfig)

			if instance["status"] != GetValueClassicOrVPC(config, "CREAT", "RUN") {
				return fmt.Errorf("invalid public ip status: %s ", instance["status"])
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicIpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpInstanceConfig(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicIpInstanceExists(resourceName, instance),
					testCheckAttribute(),
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

func TestAccResourceNcloudPublicIpInstance_classic_updateServerInstanceNo(t *testing.T) {
	instance := map[string]interface{}{}
	serverNameFoo := fmt.Sprintf("test-public-ip-foo-%s", acctest.RandString(5))
	serverNameBar := fmt.Sprintf("test-public-ip-bar-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicIpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, ""),
				SkipFunc: testOnlyClassic,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config:   testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, "${ncloud_server.foo.id}"),
				SkipFunc: testOnlyClassic,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config:   testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, "${ncloud_server.bar.id}"),
				SkipFunc: testOnlyClassic,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config:   testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, ""),
				SkipFunc: testOnlyClassic,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
		},
	})
}

func TestAccResourceNcloudPublicIpInstance_vpc_updateServerInstanceNo(t *testing.T) {
	instance := map[string]interface{}{}
	serverNameFoo := fmt.Sprintf("test-public-ip-foo-%s", acctest.RandString(5))
	serverNameBar := fmt.Sprintf("test-public-ip-bar-%s", acctest.RandString(5))
	resourceName := "ncloud_public_ip.public_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicIpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, ""),
				SkipFunc: testOnlyVpc,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config:   testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, "${ncloud_server.foo.id}"),
				SkipFunc: testOnlyVpc,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config:   testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, "${ncloud_server.bar.id}"),
				SkipFunc: testOnlyVpc,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
			{
				Config:   testAccPublicIpInstanceConfigVpcServer(serverNameFoo, serverNameBar, ""),
				SkipFunc: testOnlyVpc,
				Check:    resource.ComposeTestCheckFunc(testAccCheckPublicIpInstanceExists(resourceName, instance)),
			},
		},
	})
}

func testAccCheckPublicIpInstanceExists(n string, i map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := testAccProvider.Meta().(*ProviderConfig)

		instance, err := getPublicIp(config, rs.Primary.ID)

		if err != nil {
			return nil
		}

		if instance != nil {
			for k, v := range instance {
				i[k] = v
			}

			return nil
		}

		return fmt.Errorf("public ip instance not found")
	}
}

func testAccCheckPublicIpInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		config := testAccProvider.Meta().(*ProviderConfig)

		if rs.Type != "ncloud_public_ip" {
			continue
		}

		instance, err := getPublicIp(config, rs.Primary.ID)

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

func testAccPublicIpInstanceConfig(description string) string {
	return fmt.Sprintf(`
resource "ncloud_public_ip" "public_ip" {
	description = "%s"
}
`, description)
}

func testAccPublicIpInstanceConfigClassicServer(serverNameFoo, serverNameBar, serverInstanceNo string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "foo" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_server" "bar" {
	name = "%[2]s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_public_ip" "public_ip" {
	description = "%[1]s"
	server_instance_no = "%[3]s"
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
}
`, serverNameFoo, serverNameBar, serverInstanceNo)
}
