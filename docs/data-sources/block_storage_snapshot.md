---
subcategory: "Server"
---


# Data Source: ncloud_block_storage_snapshot

This module can be useful for getting detail of Snapshot (Block Storage) created before.

## Example Usage

```hcl
variable "snapshot_no" {}

data "ncloud_block_storage_snapshot" "snapshot" {
  id = var.snapshot_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific Snapshot to retrieve.
* `block_storage_no` - (Optional) The ID of the specific Block storage to retrieve. 
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshot_no` - The ID of Snapshot. (It is the same result as `id`)
* `name` - The name of snapshot.
* `volume_size` - The size of snapshot volume.
* `description` - Description of snapshot.