---
subcategory: "Auto Scaling"
---


# Data Source: ncloud_auto_scaling_group

This module can be useful for getting detail of Auto Scaling Group created before.

## Example Usage

```hcl
variable "auto_scaling_group_no" {}

data "ncloud_auto_scaling_group" "example" {
  id = var.auto_scaling_group_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific auto scaling group to retrieve.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_scaling_group_no` - The ID of Auto Scaling Group. (It is the same result as `id`)
* `name` - The name of Auto Scaling Group.
* `launch_configuration_no` - The number of the associated launch configuration.
* `desired_capacity` - The number of servers is adjusted according to the desired capacity value.
* `min_size` - The minimum size of the Auto Scaling Group.
* `max_size` - The maximum size of the Auto Scaling Group.
* `default_cooldown` - The cooldown time is the period set to ignore even if the monitoring event alarm occurs after the actual scaling is being performed or is completed.
* `health_check_type_code` - `SVR` or `LOADB`. Controls how health checking is done.
* `wait_for_capacity_timeout` - The maximum amount of time Terraform should wait for an ASG instance to become healthy.
* `health_check_grace_period` - time to hold health check after the server instance is put into the service with the health check hold period.
* `server_instance_no_list` - List of server instances belonging to Auto Scaling Group.

~> **NOTE:** Below attributes only support Classic environment.

* `zone_no_list` - the list of zone numbers where server instances belonging to this group will exist.

~> **NOTE:** Below attributes only support VPC environment.

* `subnet_no` - The ID of the associated Subnet.
* `vpc_no` - The ID of the associated VPC.
* `access_control_group_no_list` - The ID of the ACG.
* `target_group_list` - Target Group number list of Load Balancer.
* `server_name_prefix` - Create name beginning with the specified prefix.