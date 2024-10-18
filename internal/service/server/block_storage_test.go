package server_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccResourceNcloudBlockStorage_classic_basic(t *testing.T) {
	var storageInstance server.BlockStorage
	name := fmt.Sprintf("tf-storage-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage.storage"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckBlockStorageDestroyWithProvider(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(false)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-tf"),
					resource.TestCheckResourceAttr(resourceName, "status", "ATTAC"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "type", "SVRBS"),
					resource.TestCheckResourceAttr(resourceName, "disk_type", "NET"),
					resource.TestMatchResourceAttr(resourceName, "block_storage_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "server_instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "device_name", "/dev/xvdb"),
					resource.TestCheckResourceAttr(resourceName, "product_code", "SPBSTBSTAD000002"),
					resource.TestCheckResourceAttr(resourceName, "disk_detail_type", "HDD"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stop_instance_before_detaching"},
			},
		},
	})
}

func TestAccResourceNcloudBlockStorage_vpc_basic(t *testing.T) {
	var storageInstance server.BlockStorage
	name := fmt.Sprintf("tf-storage-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage.storage"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-tf"),
					resource.TestCheckResourceAttr(resourceName, "status", "ATTAC"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "type", "SVRBS"),
					resource.TestCheckResourceAttr(resourceName, "disk_type", "NET"),
					resource.TestMatchResourceAttr(resourceName, "block_storage_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "server_instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "product_code", "SPBSTBSTAD000006"),
					resource.TestCheckResourceAttr(resourceName, "disk_detail_type", "SSD"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stop_instance_before_detaching"},
			},
		},
	})
}

func TestAccResourceNcloudBlockStorage_vpc_kvm(t *testing.T) {
	var storageInstance server.BlockStorage
	name := fmt.Sprintf("tf-storage-kvm-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage.storage"
	zone := "KR-2"
	volumeType := "CB1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVpcConfigKvm(name, zone, volumeType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-tf"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "type", "SVRBS"),
					resource.TestCheckResourceAttr(resourceName, "disk_type", "NET"),
					resource.TestMatchResourceAttr(resourceName, "block_storage_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "server_instance_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "hypervisor_type", "KVM"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stop_instance_before_detaching"},
			},
		},
	})
}

func TestAccResourceNcloudBlockStorage_classic_ChangeServerInstance(t *testing.T) {
	var storageInstance server.BlockStorage
	name := fmt.Sprintf("tf-storage-update-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage.storage"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckBlockStorageDestroyWithProvider(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageClassicConfigUpdate(name, "ncloud_server.foo.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(false)),
				),
			},
			{
				Config: testAccBlockStorageClassicConfigUpdate(name, "ncloud_server.bar.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(false)),
				),
			},
		},
	})
}

func TestAccResourceNcloudBlockStorage_vpc_ChangeServerInstance(t *testing.T) {
	var storageInstance server.BlockStorage
	name := fmt.Sprintf("tf-storage-update-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage.storage"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVpcConfigUpdate(name, "ncloud_server.foo.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(true)),
				),
			},
			{
				Config: testAccBlockStorageVpcConfigUpdate(name, "ncloud_server.bar.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(true)),
				),
			},
		},
	})
}

func TestAccResourceNcloudBlockStorage_classic_size(t *testing.T) {
	var storageInstance server.BlockStorage
	name := fmt.Sprintf("tf-storage-size-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage.storage"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckBlockStorageDestroyWithProvider(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBlockStorageClassicConfigWithSize(name, 5),
				ExpectError: regexp.MustCompile(`expected size to be at least \(10\), got 5`),
			},
			{
				Config: testAccBlockStorageClassicConfigWithSize(name, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
				),
			},
			{
				Config: testAccBlockStorageClassicConfigWithSize(name, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
				),
			},
			{
				Config:      testAccBlockStorageClassicConfigWithSize(name, 10),
				ExpectError: regexp.MustCompile("The storage size is only expandable, not shrinking."),
			},
			{
				Config: testAccBlockStorageClassicConfigWithSize(name, 2000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceName, "size", "2000"),
				),
			},
		},
	})
}

func TestAccResourceNcloudBlockStorage_vpc_size(t *testing.T) {
	var storageInstance server.BlockStorage
	name := fmt.Sprintf("tf-storage-size-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage.storage"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBlockStorageVpcConfigWithSize(name+acctest.RandString(5), 5),
				ExpectError: regexp.MustCompile(`expected size to be at least \(10\), got 5`),
			},
			{
				Config: testAccBlockStorageVpcConfigWithSize(name, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
				),
			},
			{
				Config: testAccBlockStorageVpcConfigWithSize(name, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
				),
			},
			{
				Config:      testAccBlockStorageVpcConfigWithSize(name, 10),
				ExpectError: regexp.MustCompile("The storage size is only expandable, not shrinking."),
			},
			{
				Config: testAccBlockStorageVpcConfigWithSize(name+acctest.RandString(5), 2000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExistsWithProvider(resourceName, &storageInstance, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "size", "2000"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageExistsWithProvider(n string, i *server.BlockStorage, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		storage, err := server.GetBlockStorage(config, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if storage != nil {
			*i = *storage
			return nil
		}

		return fmt.Errorf("block storage not found")
	}
}

func testAccCheckBlockStorageDestroy(s *terraform.State) error {
	return testAccCheckBlockStorageDestroyWithProvider(s, GetTestProvider(true))
}

func testAccCheckBlockStorageDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_block_storage" {
			continue
		}
		blockStorage, err := server.GetBlockStorage(config, rs.Primary.ID)

		if blockStorage == nil {
			continue
		}
		if err != nil {
			return err
		}
		if blockStorage != nil && *blockStorage.Status != "ATTAC" {
			return fmt.Errorf("found attached block storage: %s", *blockStorage.BlockStorageInstanceNo)
		}
	}

	return nil
}

func testAccBlockStorageClassicConfigWithSize(name string, size int) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
    key_name = "%[1]s-key"
}

resource "ncloud_server" "server" {
    name = "%[1]s"
    server_image_product_code = "SPSW0LINUX000046"
    server_product_code = "SPSVRSTAND000004"
    login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_block_storage" "storage" {
    server_instance_no = ncloud_server.server.id
    name = "%[1]s-tf"
    size = "%[2]d"
}
`, name, size)
}

func testAccBlockStorageClassicConfig(name string) string {
	return testAccBlockStorageClassicConfigWithSize(name, 10)
}

func testAccBlockStorageVpcConfigWithSize(name string, size int) string {
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
        name = "hypervisor_type"
        values = ["XEN"]
    }
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
	server_spec_code = "s2-g2-s50"
	login_key_name = ncloud_login_key.loginkey.key_name
	delete_blockstorage_server_termination = true
}

resource "ncloud_block_storage" "storage" {
	server_instance_no = ncloud_server.server.id
	name = "%[1]s-tf"
	size = "%[2]d"
	hypervisor_type = "XEN"
	volume_type = "SSD"
}
`, name, size)
}

func testAccBlockStorageVpcConfig(name string) string {
	return testAccBlockStorageVpcConfigWithSize(name, 10)
}
func testAccBlockStorageClassicConfigUpdate(name, serverInstanceNo string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "foo" {
	name = "%[1]s-foo"
	server_image_product_code = "SPSW0LINUX000046"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_server" "bar" {
	name = "%[1]s-bar"
	server_image_product_code = "SPSW0LINUX000046"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_block_storage" "storage" {
	server_instance_no =  %[2]s
	name = "%[1]s-tf"
	size = "10"
}
`, name, serverInstanceNo)
}

func testAccBlockStorageVpcConfigUpdate(name, serverInstanceNo string) string {
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
        name = "hypervisor_type"
        values = ["XEN"]
    }
}

resource "ncloud_server" "foo" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s-foo"
	server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
	server_spec_code = "s2-g2-s50"
	login_key_name = ncloud_login_key.loginkey.key_name
	delete_blockstorage_server_termination = true
}

resource "ncloud_server" "bar" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s-bar"
	server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
	server_spec_code = "s2-g2-s50"
	login_key_name = ncloud_login_key.loginkey.key_name
	delete_blockstorage_server_termination = true
}

resource "ncloud_block_storage" "storage" {
	server_instance_no = %[2]s
	name = "%[1]s-tf"
	size = "10"
}
`, name, serverInstanceNo)
}

func testAccBlockStorageVpcConfigKvm(name string, zone string, volumeType string) string {
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
	zone               = "%[2]s"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

data "ncloud_server_image_numbers" "server_images" {
	hypervisor_type = "KVM"
	filter {
        name = "name"
        values = ["ubuntu-22.04-base"]
    }
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
	server_spec_code = "s2-g3"
	login_key_name = ncloud_login_key.loginkey.key_name
	delete_blockstorage_server_termination = true
}

resource "ncloud_block_storage" "storage" {
	server_instance_no = ncloud_server.server.id
	name = "%[1]s-tf"
	size = "10"
	hypervisor_type = "KVM"
	volume_type = "%[3]s"
	zone = "%[2]s"
}
`, name, zone, volumeType)
}
