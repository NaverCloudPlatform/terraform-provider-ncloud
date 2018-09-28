package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

// TODO: Fix TestAcc ErrorTestAccResourceNcloudBlockStorageBasic
//
func ignore_TestAccResourceNcloudBlockStorageSnapshotBasic(t *testing.T) {
	var snapshotInstance server.BlockStorageSnapshotInstance
	prefix := getTestPrefix()
	testLoginKeyName := prefix + "-key"
	testServerInstanceName := prefix + "-vm"
	testBlockStorageName := prefix + "-storage"
	testSnapshotName := prefix + "-snapshot"
	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if ncloud.StringValue(snapshotInstance.BlockStorageSnapshotName) != testSnapshotName {
				return fmt.Errorf("not found: %s", testSnapshotName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_block_storage_snapshot.ss",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotConfig(testLoginKeyName, testServerInstanceName, testBlockStorageName, testSnapshotName),
				//SkipFunc: func() (bool, error) { return true, nil }, // TODO: Fix TestAcc ErrorTestAccResourceNcloudBlockStorageBasic
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageSnapshotExists(
						"ncloud_block_storage_snapshot.ss", &snapshotInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_block_storage_snapshot.ss",
						"block_storage_snapshot_name",
						testSnapshotName),
				),
			},
		},
	})
}

func testAccCheckBlockStorageSnapshotExists(n string, i *server.BlockStorageSnapshotInstance) resource.TestCheckFunc {
	return testAccCheckBlockStorageSnapshotExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckBlockStorageSnapshotExistsWithProvider(n string, i *server.BlockStorageSnapshotInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		client := provider.Meta().(*NcloudAPIClient)
		snapshot, err := getBlockStorageSnapshotInstance(client, rs.Primary.ID)
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

func testAccCheckBlockStorageSnapshotDestroy(s *terraform.State) error {
	return testAccCheckBlockStorageSnapshotDestroyWithProvider(s, testAccProvider)
}

func testAccCheckBlockStorageSnapshotDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*NcloudAPIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_block_storage_snapshot" {
			continue
		}
		log.Printf("[DEBUG] testAccCheckBlockStorageSnapshotDestroyWithProvider getBlockStorageSnapshotInstance %s", rs.Primary.ID)
		snapshot, err := getBlockStorageSnapshotInstance(client, rs.Primary.ID)
		log.Printf("[DEBUG] testAccCheckBlockStorageSnapshotDestroyWithProvider snapshot %#v", snapshot)
		if snapshot == nil {
			return nil
		}
		if err != nil {
			log.Printf("[ERROR] testAccCheckBlockStorageSnapshotDestroyWithProvider err: %s", err.Error())
			return err
		}
		log.Printf("[DEBUG] testAccCheckBlockStorageSnapshotDestroyWithProvider ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceStatus.Code) %s", ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceStatus.Code))
		if ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceStatus.Code) != "TERMT" {
			return fmt.Errorf("found block storage snapshot: %s", ncloud.StringValue(snapshot.BlockStorageSnapshotInstanceNo))
		}
	}

	return nil
}

func testAccBlockStorageSnapshotConfig(testLoginKeyName string, serverInstanceName string, blockStorageName string, snapshotName string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "key" {
	"key_name" = "%s"
}

resource "ncloud_server" "vm" {
	"server_name" = "%s"
	"server_image_product_code" = "SPSW0LINUX000032"
	"server_product_code" = "SPSVRSTAND000004"
	"login_key_name" = "${ncloud_login_key.key.key_name}"
}

resource "ncloud_block_storage" "bs" {
	"server_instance_no" = "${ncloud_server.vm.id}"
	"block_storage_name" = "%s"
	"block_storage_size_gb" = "10"
}

resource "ncloud_block_storage_snapshot" "ss" {
	"block_storage_instance_no" = "${ncloud_block_storage.bs.id}"
	"block_storage_snapshot_name" = "%s"
}
`, testLoginKeyName, serverInstanceName, blockStorageName, snapshotName)
}
