---
subcategory: "Auto Scaling"
---


# Data Source: ncloud_auto_scaling_policy

This module can be useful for getting detail of Auto Scaling Policy created before.

## Example Usage

```hcl
variable "policy_name" {}
variable "auto_scaling_group_no" {}

data "ncloud_auto_scaling_policy" "example" {
  id = var.policy_name
  auto_scaling_group_no = var.auto_scaling_group_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific auto scaling policy to retrieve.
* `auto_scaling_group_no` (Required) The ID of Auto Scaling Group.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of Auto Scaling Policy.
* `adjustment_type_code` - how the number of servers is scaled when the scaling policy is performed.
* `scaling_adjustment` - Specify the adjustment value for the adjustment type.
* `cooldown` - The cooldown time is the period set to ignore even if the monitoring event alarm occurs after the actual scaling is being performed or is completed.
* `min_adjustment_step` - Change the number of server instances by the minimum adjustment width.