package server_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
	serverservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccResourceNcloudBlockStorageSnapshot_vpc_basic(t *testing.T) {
	name := fmt.Sprintf("tf-snap-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage_snapshot.snapshot"
	hypervisorType := "KVM"
	serverSpec := "s2-g3"
	zone := "KR-2"
	volumeType := "CB1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVpcBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotVpcConfig(name, hypervisorType, zone, serverSpec, volumeType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-tf"),
					resource.TestCheckResourceAttr(resourceName, "hypervisor_type", hypervisorType),
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

func testAccCheckVpcBlockStorageSnapshotDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_block_storage_snapshot" {
			continue
		}
		snapshot, err := server.GetVpcBlockStorageSnapshotDetail(config, rs.Primary.ID)

		if err != nil {
			return err
		}
		if snapshot != nil {
			return fmt.Errorf("unterminated snapshot : %s", *snapshot.BlockStorageSnapshotInstanceNo)
		}
	}

	return nil
}

// TODO: Fix TestAcc ErrorTestAccResourceNcloudBlockStorageBasic
//
//nolint:unused
func ignore_TestAccResourceNcloudBlockStorageSnapshotBasic(t *testing.T) {
	var snapshotInstance *serverservice.BlockStorageSnapshot
	prefix := GetTestPrefix()
	testLoginKeyName := prefix + "-key"
	testServerInstanceName := prefix + "-vm"
	testBlockStorageName := prefix + "-storage"
	testSnapshotName := prefix + "-snapshot"
	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if *snapshotInstance.BlockStorageSnapshotName != testSnapshotName {
				return fmt.Errorf("not found: %s", testSnapshotName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotConfig(testLoginKeyName, testServerInstanceName, testBlockStorageName, testSnapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageSnapshotExists(
						"ncloud_block_storage_snapshot.ss", snapshotInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_block_storage_snapshot.ss",
						"name",
						testSnapshotName),
				),
			},
			{
				ResourceName:      "ncloud_block_storage_snapshot.ss",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//nolint:unused
func testAccCheckBlockStorageSnapshotExists(n string, i *serverservice.BlockStorageSnapshot) resource.TestCheckFunc {
	return testAccCheckBlockStorageSnapshotExistsWithProvider(n, i, func() *schema.Provider { return GetTestProvider(false) })
}

//nolint:unused
func testAccCheckBlockStorageSnapshotExistsWithProvider(n string, i *serverservice.BlockStorageSnapshot, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		config := provider.Meta().(*conn.ProviderConfig)
		snapshot, err := serverservice.GetClassicBlockStorageSnapshotInstance(config, rs.Primary.ID)
		log.Printf("[DEBUG] testAccCheckBlockStorageSnapshotExistsWithProvider snapshot %#v", snapshot)

		if err != nil {
			return nil
		}

		if snapshot != nil {
			*i = *snapshot
			return nil
		}

		return fmt.Errorf("block storage snapshot is not found")
	}
}

//nolint:unused
func testAccCheckBlockStorageSnapshotDestroy(s *terraform.State) error {
	return testAccCheckBlockStorageSnapshotDestroyWithProvider(s, GetTestProvider(false))
}

//nolint:unused
func testAccCheckBlockStorageSnapshotDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_block_storage_snapshot" {
			continue
		}
		log.Printf("[DEBUG] testAccCheckBlockStorageSnapshotDestroyWithProvider getBlockStorageSnapshotInstance %s", rs.Primary.ID)
		snapshot, err := serverservice.GetClassicBlockStorageSnapshotInstance(config, rs.Primary.ID)
		log.Printf("[DEBUG] testAccCheckBlockStorageSnapshotDestroyWithProvider snapshot %#v", snapshot)
		if snapshot == nil {
			return nil
		}
		if err != nil {
			log.Printf("[ERROR] testAccCheckBlockStorageSnapshotDestroyWithProvider err: %s", err.Error())
			return err
		}
		log.Printf("[DEBUG] testAccCheckBlockStorageSnapshotDestroyWithProvider ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceStatus.Code) %s", *snapshot.Status)
		if *snapshot.Status != "TERMINATED" {
			return fmt.Errorf("found block storage snapshot: %s", ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceNo))
		}
	}

	return nil
}

//nolint:unused
func testAccBlockStorageSnapshotConfig(testLoginKeyName string, serverInstanceName string, blockStorageName string, snapshotName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "key" {
	key_name = "%s"
}

resource "ncloud_server" "vm" {
	name = "%s"
	server_image_product_code = "SPSW0LINUX000046"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.key.key_name}"
}

resource "ncloud_block_storage" "bs" {
	server_instance_no = "${ncloud_server.vm.id}"
	name = "%s"
	size = "10"
}

resource "ncloud_block_storage_snapshot" "ss" {
	block_storage_instance_no = "${ncloud_block_storage.bs.id}"
	name = "%s"
}
`, testLoginKeyName, serverInstanceName, blockStorageName, snapshotName)
}

func testAccBlockStorageSnapshotVpcConfig(name string, hypervisorType string, zone string, serverSpec string, volumeType string) string {
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
	zone               = "%[3]s"
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
	server_spec_code = "%[4]s"
	login_key_name = ncloud_login_key.loginkey.key_name
	delete_blockstorage_server_termination = true
}

resource "ncloud_block_storage" "storage" {
	server_instance_no = ncloud_server.server.id
	name = "%[1]s-tf"
	size = "10"
	hypervisor_type = "%[2]s"
	volume_type = "%[5]s"
	zone = "%[3]s"
}

resource "ncloud_block_storage_snapshot" "snapshot" {
    block_storage_instance_no = ncloud_block_storage.storage.id
	name = "%[1]s-tf"
    description = "Terraform test snapshot"
}
`, name, hypervisorType, zone, serverSpec, volumeType)
}
