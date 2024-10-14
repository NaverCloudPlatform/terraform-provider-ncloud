---
subcategory: "Server"
---


# Data Source: ncloud_block_storage

This module can be useful for getting detail of Block Storage created before.

## Example Usage

```terraform
variable "block_storage_no" {}

data "ncloud_block_storage" "storage" {
  id = var.block_storage_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific block storage to retrieve.
* `server_instance_no` - (Optional) The ID of the server instance associated with block storage to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

* `block_storage_no` - The ID of Block Storage. (It is the same result as `id`)
* `size` - The size of the Block Storage.
* `name` - The name of Block Storage.
* `description` - Description of Block Storage
* `disk_detail_type` - Type of Block Storage disk detail. 
* `server_name` - Server name.
* `type` - Block Storage type code.
* `device_name` - Device name.
* `product_code` - Block Storage product code.
* `status` - Block Storage instance status code.
* `disk_type` - Disk type code.
* `zone` - Available zone where the Block Storage placed.

~> **NOTE:** Arguments below only support Classic environment.

* `snapshot_no` - The ID of Block Storage Snapshot.

~> **NOTE:** Arguments below only support VPC environment.

* `max_iops` - Maximum IOPS.
* `is_encrypted_volume` - Volume encryption status. (`true` or `false`)
* `is_return_protection` - Enable return protection. (`true` or `false`)
* `hypervisor_type` - Hypervisor type. (`XEN` or `KVM`)
* `volume_type` - Volume type of the block storage. `XEN` type(` SSD` | `HDD`), `KVM`type(`FB1` | `CB1`)
