package ncloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudServer_classic_basic(t *testing.T) {
	var serverInstance NcloudServerInstance
	testServerName := getTestServerName()
	resourceName := "ncloud_server.server"
	productCode := "SPSVRSTAND000004" // vCPU 2EA, Memory 4GB, Disk 50GB

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if *serverInstance.ServerName != testServerName {
				return fmt.Errorf("not found: %s", testServerName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerClassicConfig(testServerName, productCode),
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return config.SupportVPC, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(resourceName, &serverInstance),
					testCheck(),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "server_image_product_code", "SPSW0LINUX000032"),
					resource.TestCheckResourceAttr(resourceName, "server_product_code", productCode),
					resource.TestCheckResourceAttr(resourceName, "name", testServerName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "internet_line_type", "PUBLC"),
					resource.TestCheckResourceAttr(resourceName, "zone", "KR-1"),
					resource.TestCheckResourceAttr(resourceName, "base_block_storage_disk_type", "NET"),
					resource.TestCheckResourceAttr(resourceName, "base_block_storage_size", "53687091200"),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "4294967296"),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "operation", "NULL"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUN"),
					resource.TestCheckResourceAttr(resourceName, "platform_type", "LNX32"),
					resource.TestCheckResourceAttr(resourceName, "is_fee_charging_monitoring", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_protect_server_termination", "false"),
					resource.TestCheckResourceAttr(resourceName, "server_image_name", "centos-6.3-32"),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", fmt.Sprintf("%s-key", testServerName)),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "port_forwarding_public_ip", regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)),
					resource.TestMatchResourceAttr(resourceName, "private_ip", regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)),
					resource.TestCheckResourceAttr(resourceName, "public_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_operation", "NULL"),
					resource.TestCheckResourceAttr(resourceName, "instance_status", "RUN"),
					resource.TestCheckResourceAttr(resourceName, "instance_status_name", "running"),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return config.SupportVPC, nil
				},
				ResourceName:            "ncloud_server.server",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"login_key_name", "server_product_code"},
			},
		},
	})
}

func TestAccResourceNcloudServer_vpc_basic(t *testing.T) {
	var serverInstance NcloudServerInstance
	testServerName := getTestServerName()
	resourceName := "ncloud_server.server"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerVpcConfig(testServerName, productCode),
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return !config.SupportVPC, nil
				},
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExists("ncloud_server.server", &serverInstance),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "server_image_product_code", "SW.VSVR.OS.LNX64.CNTOS.0708.B050"),
					resource.TestCheckResourceAttr(resourceName, "server_product_code", productCode),
					resource.TestCheckResourceAttr(resourceName, "name", testServerName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "zone", "KR-1"),
					resource.TestCheckResourceAttr(resourceName, "base_block_storage_disk_type", "NET"),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "8589934592"),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "instance_operation", "NULL"),
					resource.TestCheckResourceAttr(resourceName, "instance_status", "RUN"),
					resource.TestCheckResourceAttr(resourceName, "instance_status_name", "running"),
					resource.TestCheckResourceAttr(resourceName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(resourceName, "is_protect_server_termination", "false"),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", fmt.Sprintf("%s-key", testServerName)),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "public_ip", ""),
					// VPC only
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "network_interfaces.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "network_interfaces.0", regexp.MustCompile(`^\d+$`)),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return !config.SupportVPC, nil
				},
				ResourceName:            "ncloud_server.server",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"login_key_name", "server_product_code"},
			},
		},
	})
}

func TestAccResourceNcloudServer_classic_changeSpec(t *testing.T) {
	var before NcloudServerInstance
	var after NcloudServerInstance
	testServerName := getTestServerName()
	resourceName := "ncloud_server.server"
	productCode := "SPSVRSTAND000004"       // vCPU 2EA, Memory 4GB, Disk 50GB
	targetProductCode := "SPSVRSTAND000005" // vCPU 4EA, Memory 8GB, Disk 50GB

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerClassicConfig(testServerName, productCode),
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return config.SupportVPC, nil
				},
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "4294967296"),
				),
			},
			{
				Config: testAccServerClassicConfig(testServerName, targetProductCode),
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return config.SupportVPC, nil
				},
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "8589934592"),
					testAccCheckInstanceNotRecreated(t, &before, &after),
				),
			},
			{
				ResourceName:      "ncloud_server.server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudServer_vpc_changeSpec(t *testing.T) {
	var before NcloudServerInstance
	var after NcloudServerInstance
	testServerName := getTestServerName()
	resourceName := "ncloud_server.server"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"       // vCPU 2EA, Memory 8GB, Disk 50GB
	targetProductCode := "SVR.VSVR.STAND.C004.M016.NET.HDD.B050.G002" // vCPU 4EA, Memory 16GB, Disk 50GB

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerVpcConfig(testServerName, productCode),
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return !config.SupportVPC, nil
				},
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "8589934592"),
				),
			},
			{
				Config: testAccServerVpcConfig(testServerName, targetProductCode),
				SkipFunc: func() (bool, error) {
					config := testAccProvider.Meta().(*ProviderConfig)
					return !config.SupportVPC, nil
				},
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "17179869184"),
					testAccCheckInstanceNotRecreated(t, &before, &after),
				),
			},
			{
				ResourceName:      "ncloud_server.server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckServerExists(n string, i *NcloudServerInstance) resource.TestCheckFunc {
	return testAccCheckInstanceExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckInstanceExistsWithProvider(n string, i *NcloudServerInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		config := provider.Meta().(*ProviderConfig)
		instance, err := getServerInstance(config, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if instance != nil {
			*i = *instance
			return nil
		}

		return fmt.Errorf("server instance not found")
	}
}

func testAccCheckInstanceNotRecreated(t *testing.T, before, after *NcloudServerInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *before.ServerInstanceNo != *after.ServerInstanceNo {
			t.Fatalf("Ncloud Instance IDs have changed. Before %s. After %s", *before.ServerInstanceNo, *after.ServerInstanceNo)
		}
		return nil
	}
}

func testAccCheckServerDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_server" {
			continue
		}
		instance, err := getServerInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance == nil {
			continue
		}

		if *instance.ServerInstanceStatusName != "terminating" {
			return fmt.Errorf("found unterminated instance: %s", *instance.ServerInstanceNo)
		}
	}

	return nil
}

func getTestServerName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testServerName := fmt.Sprintf("tf-%d-vm", rInt)
	return testServerName
}

func testAccServerVpcConfig(testServerName, productCode string) string {
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
	zone               = "KR-1"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_network_interface" "test" {
	name                  = "%[1]s"
	description           = "for acc test"
	subnet_no             = ncloud_subnet.test.id
	access_control_groups = [ncloud_vpc.test.default_access_control_group_no]
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0708.B050"
	server_product_code = "%[2]s"
	login_key_name = ncloud_login_key.loginkey.key_name
	network_interfaces = [ncloud_network_interface.test.id]
}
`, testServerName, productCode)
}

func testAccServerClassicConfig(testServerName, productCode string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "server" {
	name = "%[1]s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "%[2]s"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}
`, testServerName, productCode)
}
