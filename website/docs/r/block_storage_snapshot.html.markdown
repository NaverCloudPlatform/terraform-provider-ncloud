---
layout: "ncloud"
page_title: "NCLOUD: ncloud_block_storage_snapshot"
sidebar_current: "docs-ncloud-resource-block-storage-snapshot"
description: |-
  Provides a ncloud block storage snapshot resource.
---

# ncloud_block_storage

Provides a ncloud block storage snapshot resource.

## Example Usage

```hcl
resource "ncloud_block_storage_snapshot" "snapshot" {
	"block_storage_instance_no" = "812345"
	"block_storage_snapshot_name" = "tf-test-snapshot1"
	"block_storage_snapshot_description" = "Terraform test snapshot1"
}
```

## Argument Reference

The following arguments are supported:

* `block_storage_instance_no` - (Required) Block storage instance No for creating snapshot.
* `block_storage_snapshot_name` - (Optional) Block storage snapshot name to create. default : Ncloud assigns default values.
* `block_storage_snapshot_description` - (Optional) Descriptions on a snapshot to create.

## Attributes Reference

* `block_storage_snapshot_instance_no` - Block Storage Snapshot Instance Number
* `block_storage_snapshot_volume_size` - Block Storage Snapshot Volume Size
* `original_block_storage_instance_no` - Original Block Storage Instance Number
* `original_block_storage_name` - Original Block Storage Name
* `block_storage_snapshot_instance_status`
    * `code` - Block Storage Snapshot Instance Status code
    * `code_name` - Block Storage Snapshot Instance Status name
* `block_storage_snapshot_instance_operation`
    * `code` - Block Storage Snapshot Instance Operation code
    * `code_name` - Block Storage Snapshot Instance Operation name
* `block_storage_snapshot_instance_status_name` - Block Storage Snapshot Instance Status Name
* `create_date` - Creation date of the block storage snapshot instance
* `block_storage_snapshot_instance_description` - Block Storage Snapshot Instance Description
* `server_image_product_code` - Server Image Product Code
* `os_information` - OS Information