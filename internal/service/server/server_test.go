package server_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/common"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	serverservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccResourceNcloudServer_vpc_basic(t *testing.T) {
	var serverInstance serverservice.ServerInstance
	testServerName := GetTestServerName()
	resourceName := "ncloud_server.server"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerVpcConfig(testServerName, productCode),
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExistsWithProvider("ncloud_server.server", &serverInstance, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "server_image_product_code", "SW.VSVR.OS.LNX64.ROCKY.0810.B050"),
					resource.TestCheckResourceAttr(resourceName, "server_product_code", productCode),
					resource.TestCheckResourceAttr(resourceName, "name", testServerName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "zone", "KR-2"),
					resource.TestCheckResourceAttr(resourceName, "base_block_storage_disk_type", "NET"),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "8589934592"),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(resourceName, "is_protect_server_termination", "false"),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", fmt.Sprintf("%s-key", testServerName)),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "public_ip", ""),
					// VPC only
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "network_interface.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.order"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.network_interface_no"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.subnet_no"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.private_ip"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudServer_vpc_kvm(t *testing.T) {
	var serverInstance serverservice.ServerInstance
	testServerName := GetTestServerName()
	resourceName := "ncloud_server.server"
	specCode := "s2-g3"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerImageNumberVpcConfig(testServerName, specCode),
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExistsWithProvider("ncloud_server.server", &serverInstance, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "server_spec_code", specCode),
					resource.TestCheckResourceAttr(resourceName, "name", testServerName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "is_protect_server_termination", "false"),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", fmt.Sprintf("%s-key", testServerName)),
					// VPC only
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "network_interface.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.order"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.network_interface_no"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.subnet_no"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.private_ip"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudServer_vpc_networkInterface(t *testing.T) {
	var serverInstance serverservice.ServerInstance
	testServerName := GetTestServerName()
	resourceName := "ncloud_server.server"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerVpcConfigNetworkInterface(testServerName, productCode),
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExistsWithProvider("ncloud_server.server", &serverInstance, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "server_image_product_code", "SW.VSVR.OS.LNX64.ROCKY.0810.B050"),
					resource.TestCheckResourceAttr(resourceName, "server_product_code", productCode),
					resource.TestCheckResourceAttr(resourceName, "name", testServerName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "zone", "KR-2"),
					resource.TestCheckResourceAttr(resourceName, "base_block_storage_disk_type", "NET"),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "8589934592"),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "platform_type", "LNX64"),
					resource.TestCheckResourceAttr(resourceName, "is_protect_server_termination", "false"),
					resource.TestCheckResourceAttr(resourceName, "login_key_name", fmt.Sprintf("%s-key", testServerName)),
					resource.TestMatchResourceAttr(resourceName, "instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "public_ip", ""),
					// VPC only
					resource.TestMatchResourceAttr(resourceName, "subnet_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "network_interface.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.order"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.network_interface_no"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.subnet_no"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.private_ip"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudServer_vpc_changeSpec(t *testing.T) {
	var before serverservice.ServerInstance
	var after serverservice.ServerInstance
	testServerName := GetTestServerName()
	resourceName := "ncloud_server.server"
	productCode := "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"       // vCPU 2EA, Memory 8GB, Disk 50GB
	targetProductCode := "SVR.VSVR.STAND.C004.M016.NET.HDD.B050.G002" // vCPU 4EA, Memory 16GB, Disk 50GB

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerVpcConfig(testServerName, productCode),
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExistsWithProvider(resourceName, &before, TestAccProvider),
					resource.TestCheckResourceAttr(resourceName, "cpu_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "8589934592"),
				),
			},
			{
				Config: testAccServerVpcConfig(testServerName, targetProductCode),
				Check: resource.ComposeTestCheckFunc(testAccCheckServerExistsWithProvider(resourceName, &after, TestAccProvider),
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

func TestConvertToMap(t *testing.T) {
	i := &serverservice.ServerInstance{
		ZoneNo:                     ncloud.String("KR-1"),
		ServerName:                 ncloud.String("tf-server"),
		IsProtectServerTermination: ncloud.Bool(true),
		CpuCount:                   ncloud.Int32(2),
	}

	m := common.ConvertToMap(i)

	if m["cpu_count"].(float64) != 2 {
		t.Fatalf("'cpu_count' expected '2' but %s", m["cpu_count"])
	}

	if m["is_protect_server_termination"].(bool) != true {
		t.Fatalf("'is_protect_server_termination' expected 'true' but %s", m["is_protect_server_termination"])
	}

	if m["name"].(string) != "tf-server" {
		t.Fatalf("'cpu_count' expected '2' but %s", m["name"])
	}

	if _, ok := m["network_interface"]; !ok {
		t.Fatalf("'network_interface' expected 'nil' but %s", m["network_interface"])
	}
}

func TestConvertToArrayMap(t *testing.T) {
	i := &serverservice.ServerInstance{
		ZoneNo:                     ncloud.String("KR-1"),
		ServerName:                 ncloud.String("tf-server"),
		IsProtectServerTermination: ncloud.Bool(true),
		CpuCount:                   ncloud.Int32(2),
	}
	var list []*serverservice.ServerInstance
	list = append(list, i)

	m := common.ConvertToArrayMap(list)

	if m[0]["cpu_count"].(float64) != 2 {
		t.Fatalf("'cpu_count' expected '2' but %s", m[0]["cpu_count"])
	}

	if m[0]["is_protect_server_termination"].(bool) != true {
		t.Fatalf("'is_protect_server_termination' expected 'true' but %s", m[0]["is_protect_server_termination"])
	}

	if m[0]["name"].(string) != "tf-server" {
		t.Fatalf("'cpu_count' expected '2' but %s", m[0]["name"])
	}

	if _, ok := m[0]["network_interface"]; !ok {
		t.Fatalf("'network_interface' expected 'nil' but %s", m[0]["network_interface"])
	}
}

func testAccCheckServerExistsWithProvider(n string, i *serverservice.ServerInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		instance, err := serverservice.GetServerInstance(config, rs.Primary.ID)
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

func testAccCheckInstanceNotRecreated(t *testing.T, before, after *serverservice.ServerInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *before.ServerInstanceNo != *after.ServerInstanceNo {
			t.Fatalf("Ncloud Instance IDs have changed. Before %s. After %s", *before.ServerInstanceNo, *after.ServerInstanceNo)
		}
		return nil
	}
}

func testAccCheckServerDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, TestAccProvider)
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_server" {
			continue
		}
		instance, err := serverservice.GetServerInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("found unterminated instance: %s", *instance.ServerInstanceNo)
		}
	}

	return nil
}

func testAccServerImageNumberVpcConfig(testServerName, specCode string) string {
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

data "ncloud_server_image_numbers" "server_images" {
    filter {
        name = "name"
        values = ["ubuntu-22.04-base"]
    }
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
	server_spec_code = "%[2]s"
	login_key_name = ncloud_login_key.loginkey.key_name
}
`, testServerName, specCode)
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
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
	server_product_code = "%[2]s"
	login_key_name = ncloud_login_key.loginkey.key_name
}
`, testServerName, productCode)
}

func testAccServerVpcConfigNetworkInterface(testServerName, productCode string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "public_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-pub"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_subnet" "private_subnet" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s-priv"
	subnet             = "10.5.1.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PRIVATE"
	usage_type         = "GEN"
}

resource "ncloud_network_interface" "eth0" {
	name                  = "%[1]s-eth-0"
	subnet_no             = ncloud_subnet.public_subnet.id
	access_control_groups = [ncloud_vpc.test.default_access_control_group_no]
}

resource "ncloud_network_interface" "eth1" {
	name                  = "%[1]s-eth-1"
	subnet_no             = ncloud_subnet.private_subnet.id
	access_control_groups = [ncloud_vpc.test.default_access_control_group_no]
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.public_subnet.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
	server_product_code = "%[2]s"
	login_key_name = ncloud_login_key.loginkey.key_name
	network_interface {
		order = 0
		network_interface_no = ncloud_network_interface.eth0.id
	}
	
	network_interface {
		order = 1
		network_interface_no = ncloud_network_interface.eth1.id
	}
}
`, testServerName, productCode)
}
