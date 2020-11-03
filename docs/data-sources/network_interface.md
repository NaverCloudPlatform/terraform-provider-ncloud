# Data Source: ncloud_network_interface

This module can provide useful for get detail of Network Interface created before.

~> **NOTE:** This resource only support VPC environment.

## Example Usage

```hcl
variable "network_interface_no" {}

data "ncloud_network_interface_no" "nic" {
  id = var.network_interface_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific Network Interface to retrieve.
* `name` - (Optional) The name of the specific Network Interface to retrieve.
* `private_ip` - (Optional) The Private IP of the specific Network Interface to retrieve.  
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `network_interface_no` - The ID of Network Interface. (It is the same result as `id`)
* `subnet_no` - The ID of the associated Subnet.
* `description` - Description of Network Interface.
* `access_control_groups` - List of ACG ID applied network interfaces.
* `server_instance_no` - The ID of server instance assigned network interface.
* `status` - The status of Network Interface.
* `instance_type` - Type of server instance.
* `is_default` - Whether is default or not by Server instance creation.