---
subcategory: "Server"
---


# Data Source: ncloud_init_script

This module can be useful for getting detail of Init script created before.

## Example Usage

```hcl
variable "init_script_no" {}

data "ncloud_init_script" "init_script" {
  id = var.init_script_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific Init script to retrieve. 
* `name` - (Optional) The name of the specific Init script to retrieve. 
* `os_type` - (Optional) Type of O/S to apply server instance. Accepted values: `LNX` (LINUX) | `WND` (WINDOWS)
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `init_script_no` - The ID of Init script. (It is the same result as `id`)
* `content` - Initialization script content.
* `description` - Description of Init script