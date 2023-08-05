---
subcategory: "Load Balancer"
---


# Data Source: ncloud_lb_listener

This module can be useful for getting detail of Load Balancer Listener created before.

## Example Usage

```hcl
variable "load_balancer_no" {}

variable "load_balancer_listener_no" {}

data "ncloud_lb_listener" "test" {
  id = var.load_balancer_listener_no
  load_balancer_no = var.load_balancer_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific listener to retrieve.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listener_no` - The ID of listener (It is the same result as id).
* `rule_no_list` - The list of listener rule number
* `load_balancer_no` - The ID of the load balancer.
* `target_group_no` - The ID of the target group.
* `port` - The port on which the load balancer is listening.
* `protocol` - The protocol type for the listener.
* `tls_min_version_type` - The TLS minimum supported version type code.
* `use_http2` - Whether to use HTTP/2 protocol.
* `ssl_certificate_no` - The ID of the SSL certificate.