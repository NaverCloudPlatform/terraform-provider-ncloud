package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccResourceNcloudBlockStorageSnapshotBasic(t *testing.T) {
	var snapshotInstance server.BlockStorageSnapshotInstance
	prefix := getTestPrefix()
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
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_block_storage_snapshot.snapshot",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotConfig(testServerInstanceName, testBlockStorageName, testSnapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageSnapshotExists(
						"ncloud_block_storage_snapshot.snapshot", &snapshotInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_block_storage_snapshot.snapshot",
						"block_storage_snapshot_name",
						testSnapshotName),
					resource.TestCheckResourceAttr(
						"ncloud_block_storage_snapshot.snapshot",
						"block_storage_instance_snapshot_status.code",
						"INIT"),
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
		storage, err := getBlockStorageSnapshotInstance(client, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if storage != nil {
			*i = *storage
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
		snapshot, err := getBlockStorageSnapshotInstance(client, rs.Primary.ID)

		if snapshot == nil {
			break
		}
		if err != nil {
			return err
		}
		if snapshot != nil && *snapshot.BlockStorageSnapshotInstanceStatus.Code != "TERMT" {
			return fmt.Errorf("found block storage snapshot: %s", *snapshot.BlockStorageSnapshotInstanceNo)
		}
	}

	return nil
}

func testAccBlockStorageSnapshotConfig(serverInstanceName string, blockStorageName string, snapshotName string) string {
	return fmt.Sprintf(`
resource "ncloud_server" "server" {
	"server_name" = "%s"
	"server_image_product_code" = "SPSW0LINUX000032"
	"server_product_code" = "SPSVRSTAND000004"
}

resource "ncloud_block_storage" "storage" {
	"server_instance_no" = "${ncloud_server.server.id}"
	"block_storage_name" = "%s"
	"block_storage_size_gb" = "10"
}

resource "ncloud_block_storage_snapshot" "snapshot" {
	"block_storage_instance_no" = "${ncloud_block_storage.storage.id}"
	"block_storage_snapshot_name" = "%s"
}
`, serverInstanceName, blockStorageName, snapshotName)
}
