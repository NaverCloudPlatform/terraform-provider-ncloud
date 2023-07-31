---
subcategory: "Server"
---


# Resource: ncloud_block_storage

Provides a Block Storage resource.

## Example Usage

```hcl
resource "ncloud_block_storage" "storage" {
	server_instance_no = "812345"
	name = "tf-test-storage1"
	size = "10"
}
```

## Argument Reference

The following arguments are supported:

* `size` - (Required) The size of the block storage to create. It is automatically set when you take a snapshot.
* `server_instance_no` - **(Required) When first created**. (Optional) After creation. Server instance ID to which you want to assign the block storage.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create.
* `disk_detail_type` - (Optional) Type of block storage disk detail to create. Default `SSD`. Accepted values: `SSD` | `HDD` 
* `stop_instance_before_detaching` - (Optional, Boolean) Set this to true to ensure that the target instance is stopped before trying to detach the block storage. It stops the instance, if it is not already stopped.
	> If `stop_instance_before_detaching` is `true`, server will be stopped and **will not start automatically**. User must start server instance manually via NCLOUD console or API.

~> **NOTE:** Below arguments only support Classic environment.

* `zone` - (Optional) The availability zone in which the block storage instance will be created.
* `snapshot_no` - (Optional) Create the block storage from the snapshots you take.

## Attributes Reference

* `id` - The ID of Block storage instance.
* `block_storage_no` - The ID of Block storage instance. (It is the same result as `id`)
* `server_name` - Server name.
* `type` - Block storage type code.
* `device_name` - Device name.
* `product_code` - Block storage product code.
* `status` - Block storage instance status code.
* `disk_type` - Disk type code.