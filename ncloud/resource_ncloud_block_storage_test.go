package ncloud

import (
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudBlockStorageBasic(t *testing.T) {
	var storageInstance server.BlockStorageInstance
	prefix := getTestPrefix()
	testServerInstanceName := prefix + "-vm"
	testBlockStorageName := prefix + "-storage"
	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if *storageInstance.BlockStorageName != testBlockStorageName {
				return fmt.Errorf("not found: %s", testBlockStorageName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageConfig(testServerInstanceName, testBlockStorageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageExists(
						"ncloud_block_storage.storage", &storageInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_block_storage.storage",
						"name",
						testBlockStorageName),
					resource.TestCheckResourceAttr(
						"ncloud_block_storage.storage",
						"instance_status",
						"ATTAC"),
				),
			},
			{
				ResourceName:      "ncloud_block_storage.storage",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBlockStorageExists(n string, i *server.BlockStorageInstance) resource.TestCheckFunc {
	return testAccCheckBlockStorageExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckBlockStorageExistsWithProvider(n string, i *server.BlockStorageInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		client := provider.Meta().(*ProviderConfig).Client
		storage, err := getBlockStorageInstance(client, rs.Primary.ID)
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
	return testAccCheckBlockStorageDestroyWithProvider(s, testAccProvider)
}

func testAccCheckBlockStorageDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*ProviderConfig).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_block_storage" {
			continue
		}
		blockStorage, err := getBlockStorageInstance(client, rs.Primary.ID)

		if blockStorage == nil {
			continue
		}
		if err != nil {
			return err
		}
		if blockStorage != nil && *blockStorage.BlockStorageInstanceStatus.Code != "ATTAC" {
			return fmt.Errorf("found attached block storage: %s", *blockStorage.BlockStorageInstanceNo)
		}
	}

	return nil
}

func testAccBlockStorageConfig(serverInstanceName string, blockStorageName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%s-key"
}

resource "ncloud_server" "server" {
	name = "%s"
	server_image_product_code = "SPSW0LINUX000032"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_block_storage" "storage" {
	server_instance_no = "${ncloud_server.server.id}"
	name = "%s"
	size = "10"
}
`, serverInstanceName, serverInstanceName, blockStorageName)
}
