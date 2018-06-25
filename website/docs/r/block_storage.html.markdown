---
layout: "ncloud"
page_title: "NCLOUD: ncloud_block_storage"
sidebar_current: "docs-ncloud-resource-block-storage"
description: |-
  Provides an ncloud block storage resource.
---

# ncloud_block_storage

Provides an ncloud block storage resource.

## Example Usage

```hcl
resource "ncloud_block_storage" "storage" {
	"server_instance_no" = "${var.server_instance_no}"
	"block_storage_name" = "tf-test-storage1"
	"block_storage_size_gb" = "10"
}
```

## Argument Reference

The following arguments are supported:

* `block_storage_name` - (Optional) Block storage name to create default : Ncloud configures it by itself.
* `block_storage_size_gb` - (Required) Enter a block storage size to ceate. You can enter by the unit of GB. Up to 1000GB you can enter.
* `block_storage_description` - (Optional) Block storage descriptions
* `server_instance_no` - (Required) Server instance No. to attach. It is required and you can get a server instance No. by calling getServerInstanceList.
* `disk_detail_type_code` - (Optional) You can choose a disk detail type code of HDD and SSD. default : HDD

## Attributes Reference

* `block_storage_instance_no` - Block storage instance no
* `block_storage_size` - Block storage size in bytes
* `server_name` - Server name
* `block_storage_type`
    * `code` - Block storage type code
    * `code_name` - Block storage type name
* `device_name` - Device name
* `block_storage_product_code` - Block storage product code
* `block_storage_instance_status`
    * `code` - Block storage instance status code
    * `code_name` - Block storage instance status code name
* `block_storage_instance_operation` - Block storage instance operation
* `block_storage_instance_status_name` - Block storage instance status name
* `create_date` - Creation date of the block storage
* `disk_type`
    * `code` - Disk type code
    * `code_name` - Disk type name
* `disk_detail_type`
    * `code` - Disk detail code
    * `code_name` - Disk detail name
