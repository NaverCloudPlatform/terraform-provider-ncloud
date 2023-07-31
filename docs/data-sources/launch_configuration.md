---
subcategory: "Auto Scaling"
---


# Data Source: ncloud_launch_configuration

This module can be useful for getting detail of Launch Configuration created before.

## Example Usage

```hcl
variable "launch_configuration_no" {}

data "ncloud_launch_configuration" "example" {
  id = var.launch_configuration_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific launch configuration to retrieve.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `launch_configuration_no` - The ID of Launch Configuration. (It is the same result as `id`)
* `name` - The name of Launch Configuration.
* `server_image_product_code` - Server image product code.
* `server_product_code` - Server product code.
* `member_server_image_no` - The ID of Member server image.
* `login_key_name` - The login key name to encrypt with the public key.
* `init_script_no` - Set init script ID, The server can run a user-set initialization script at first boot.

~> **NOTE:** Below attributes only support Classic environment.

* `user_data` - defined actionable scripts.

~> **NOTE:** Below attributes only support VPC environment.

* `is_encrypted_volume` - Whether to encrypt basic block storage if server image is RHV.