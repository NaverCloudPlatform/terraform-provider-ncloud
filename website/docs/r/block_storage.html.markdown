---
layout: "ncloud"
page_title: "NCLOUD: ncloud_block_storage"
sidebar_current: "docs-ncloud-resource-block-storage"
description: |-
  Provides a ncloud block storage resource.
---

# ncloud_block_storage

Provides a ncloud block storage resource.

## Example Usage

```hcl
resource "ncloud_block_storage" "storage" {
	"server_instance_no" = "812345"
	"name" = "tf-test-storage1"
	"size" = "10"
}
```

## Argument Reference

The following arguments are supported:

* `size` - (Required) Enter a block storage size to ceate. You can enter by the unit of GB. Up to 1000GB you can enter.
* `server_instance_no` - (Required) Server instance No. to attach. It is required and you can get a server instance No. by calling getServerInstanceList.
* `name` - (Optional) Block storage name to create default : Ncloud configures it by itself.
* `description` - (Optional) Block storage descriptions
* `disk_detail_type` - (Optional) You can choose a disk detail type code of HDD and SSD. default : HDD

## Attributes Reference

* `instance_no` - Block storage instance no
* `server_name` - Server name
* `type` - Block storage type code
* `device_name` - Device name
* `product_code` - Block storage product code
* `instance_status` - Block storage instance status code
* `instance_operation` - Block storage instance operation
* `instance_status_name` - Block storage instance status name
* `create_date` - Creation date of the block storage
* `disk_type` - Disk type code
