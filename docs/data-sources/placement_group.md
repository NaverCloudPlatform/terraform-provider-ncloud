---
subcategory: "Server"
---


# Data Source: ncloud_placement_group

This module can be useful for getting detail of Placement group created before.

## Example Usage

```hcl
variable "placement_group_no" {}

data "ncloud_placement_group" "group-a" {
  id = var.placement_group_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of specific Placement group to retrieve.
* `name` - (Optional) The name of specific Placement group to retrieve.
* `placement_group_type` - (Optional) Type of placement group to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
 
* `placement_group_no` - The ID of Placement group. (It is the same result as `id`)
