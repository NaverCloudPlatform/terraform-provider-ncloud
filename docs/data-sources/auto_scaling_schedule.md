---
subcategory: "Auto Scaling"
---


# Data Source: ncloud_auto_scaling_policy

This module can be useful for getting detail of Auto Scaling Schedule created before.

## Example Usage

```hcl
variable "schedule_name" {}
variable "auto_scaling_group_no" {}

data "ncloud_auto_scaling_schedule" "example" {
  id = var.policy_name
  auto_scaling_group_no = var.auto_scaling_group_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific auto scaling schedule to retrieve.
* `auto_scaling_group_no` (Required) The ID of Auto Scaling Group.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of Auto Scaling Schedule.
* `desired_capacity` - The number of servers is adjusted according to the desired capacity value.
* `min_size` - The minimum size of the Auto Scaling Group.
* `max_size` - The maximum size of the Auto Scaling Group.
* `start_time` - the date and time when the schedule first starts.
* `end_time` - the date and time when the schedule end.
* `recurrence` - Repeat Settings.
* `auto_scaling_group_no` - The number of the auto scaling group.

~> **NOTE:** Below attributes only support VPC environment.

* `time_zone` - the time band for the repeat settings.