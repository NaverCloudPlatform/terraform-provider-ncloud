---
subcategory: "Load Balancer"
---


# Data Source: ncloud_lb

This module can be useful for getting detail of Load Balancer created before.

## Example Usage

```hcl
variable "load_balancer_no" {}

data "ncloud_lb" "test" {
	id = var.load_balancer_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific load balancer to retrieve.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `load_balancer_no` - The ID of load balancer (It is the same result as id).
* `name` - The name of the load balancer.
* `description` - The description of the load balancer.
* `network_type` - The network type of load balancer.
* `idle_timeout` - The time in seconds that the idle timeout.
* `type` - The type of load balancer.
* `throughput_type` - The performance type code of load balancer.
* `subnet_no_list` - A list of IDs in the associated Subnets.
* `domain` - Domain name of load balancer.
* `vpc_no` - The ID of the associated VPC.
* `ip_list` - A list of IP address of load balancer.