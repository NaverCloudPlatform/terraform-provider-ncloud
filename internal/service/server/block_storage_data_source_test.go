package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudBlockStorage_classic_basic(t *testing.T) {
	resourceName := "ncloud_block_storage.storage"
	dataName := "data.ncloud_block_storage.by_id"
	name := fmt.Sprintf("tf-ds-storage-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ClassicProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ComposeConfig(
					testAccBlockStorageClassicConfig(name),
					testAccDataSourceNcloudBlockStorageConfig,
				),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "product_code", resourceName, "product_code"),
					resource.TestCheckResourceAttrPair(dataName, "size", resourceName, "size"),
					resource.TestCheckResourceAttrPair(dataName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataName, "disk_detail_type", resourceName, "disk_detail_type"),
					resource.TestCheckResourceAttrPair(dataName, "disk_type", resourceName, "disk_type"),
					resource.TestCheckResourceAttrPair(dataName, "block_storage_no", resourceName, "block_storage_no"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "device_name", resourceName, "device_name"),
					TestAccCheckDataSourceID("data.ncloud_block_storage.by_filter"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudBlockStorage_vpc_basic(t *testing.T) {
	resourceName := "ncloud_block_storage.storage"
	dataName := "data.ncloud_block_storage.by_id"
	name := fmt.Sprintf("tf-ds-storage-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ComposeConfig(
					testAccBlockStorageVpcConfig(name),
					testAccDataSourceNcloudBlockStorageConfig,
				),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "status", resourceName, "status"),
					resource.TestCheckResourceAttrPair(dataName, "product_code", resourceName, "product_code"),
					resource.TestCheckResourceAttrPair(dataName, "size", resourceName, "size"),
					resource.TestCheckResourceAttrPair(dataName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataName, "disk_detail_type", resourceName, "disk_detail_type"),
					resource.TestCheckResourceAttrPair(dataName, "disk_type", resourceName, "disk_type"),
					resource.TestCheckResourceAttrPair(dataName, "block_storage_no", resourceName, "block_storage_no"),
					resource.TestCheckResourceAttrPair(dataName, "server_instance_no", resourceName, "server_instance_no"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "device_name", resourceName, "device_name"),
					TestAccCheckDataSourceID("data.ncloud_block_storage.by_filter"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudBlockStorageConfig = `
data "ncloud_block_storage" "by_id" {
	id = ncloud_block_storage.storage.id
}

data "ncloud_block_storage" "by_filter" {
	filter {
		name = "block_storage_no"
		values = [ncloud_block_storage.storage.id]
	}
}
`
