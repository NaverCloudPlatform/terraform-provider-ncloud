---
subcategory: "Load Balancer"
---


# Data Source: ncloud_lb_target_group

This module can be useful for getting detail of Load Balancer Target Group created before.

## Example Usage

```hcl
variable "target_group_no" {}

data "ncloud_lb_target_group" "test" {
  id = var.target_group_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific target group to retrieve.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `target_group_no` - The ID of target group (It is the same result as id).
* `load_balancer_instance_no` - The ID of the Load Balancer associated with the Target Group.
* `name` - The name of the target group.
* `port` - The port on which targets receive traffic.
* `protocol` - The protocol to use for routing traffic to the targets.
* `description` - The description of the target group.
* `health_check` - The health check to check the health of the target.
    * `cycle` - The number of health check cycle.
    * `down_threshold` - The number of health check failure threshold. You can determine the number of consecutive health check failures that are required before a health check is considered a failed state.
    * `up_threshold` - The number of health check normal threshold. You can determine the number of consecutive health checks that are required before health checks are considered success state.
    * `http_method` - The HTTP method for the health check. You can determine which HTTP method to use for health checks.
    * `port` - The port to use for health checks.
    * `protocol` - The type of protocol to use for health checks.
    * `url_path` - The URL path of the health check.
* `target_no_list` - The list of target number to bind to the target group.
* `target_type` - The type of target to be added to the target group.
* `vpc_no` - The ID of the VPC in to create the target group.
* `use_sticky_session` - Whether to use session specific access.
* `use_proxy_protocol` - Whether to use a proxy protocol.
* `algorithm_type` - The type of algorithm to use for load balancing.