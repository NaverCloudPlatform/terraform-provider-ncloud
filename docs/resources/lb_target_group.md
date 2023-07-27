---
subcategory: "Load Balancer"
---


# Resource: ncloud_lb_target_group

Provides a Target Group resource.

## Example Usage
```hcl
resource "ncloud_lb_target_group" "test" {
  vpc_no   = ncloud_vpc.test.vpc_no
  protocol = "HTTP"
  target_type = "VSVR"
  port        = 8080
  description = "for test"
  health_check {
    protocol = "HTTP"
    http_method = "GET"
    port           = 8080
    url_path       = "/monitor/l7check"
    cycle          = 30
    up_threshold   = 2
    down_threshold = 2
  }
  algorithm_type = "RR"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the target group.
* `port` - (Optional) The port on which targets receive traffic. Default: `80`. Valid from `1` to `65534`.
* `protocol` - (Required) The protocol to use for routing traffic to the targets. Accepted values: `TCP` | `PROXY_TCP` | `HTTP` | `HTTPS`. The protocol you use determines which type of load balancer is applicable. `APPLICATION` Load Balancer Accepted values: `HTTP` | `HTTPS`, `NETWORK` Load Balancer Accepted values : `TCP`, `NETWORK_PROXY` Load Balancer Accepted values : `PROXY_TCP`.
* `description` - (Optional) The description of the target group.
* `health_check` - (Optional) The health check to check the health of the target.
    * `cycle` - (Optional) The number of health check cycle. Default: `30`. Valid from `5` to `300`.
    * `down_threshold` - (Optional) The number of health check failure threshold. You can determine the number of consecutive health check failures that are required before a health check is considered a failed state. Default: `2`. Valid from `2` to `10`.
    * `up_threshold` - (Optional) The number of health check normal threshold. You can determine the number of consecutive health checks that are required before health checks are considered success state. Default: `2`.  Valid from `2` to `10`.
    * `http_method` - (Optional) The HTTP method for the health check. You can determine which HTTP method to use for health checks. If the health check protocol type is `HTTP` or `HTTPS`, be sure to enter it. Accepted values: `HEAD` | `GET`.
    * `port` - (Optional) The port to use for health checks. Default: 80. Valid from `1` to `65534`.
    * `protocol` - (Required) The type of protocol to use for health checks. If the target group protocol type is `TCP` or `PROXY_TCP`, Heal Check Protocol is only valid for `TCP`. If the target group protocol type is `HTTP` or `HTTPS`, Heal Check Protocol is valid only for `HTTP` and `HTTPS`.
    * `url_path` - (Optional) The URL path of the health check. Valid only if Health Check protocol type is `HTTP` or `HTTPS`. URL path must begin with `/`.
* `target_type` - (Optional) The type of target to be added to the target group.
* `vpc_no` - (Required) The ID of the VPC in to create the target group.
* `use_sticky_session` - (Optional) Whether to use session specific access. 
* `use_proxy_protocol` - (Optional) Whether to use a proxy protocol. Valid only available if the target group type selected is `TCP` | `HTTP` | `HTTPS`.
* `algorithm_type` - (Optional) The type of algorithm to use for load balancing. Accepted values: `RR`(Round Robin) | `SIPHS`(Source IP Hash) | `LC`(Least Connection) | `MH`(Maglev Hash). `RR` | `SIPHS` | `LC` are valid only if the target group type is `PROXY_TCP`, `HTTP` or `HTTPS`. `MH` | `RR` are valid only if the target group type is `TCP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of target group.
* `target_group_no` - The ID of target group (It is the same result as id).
* `load_balancer_instance_no` - The ID of the Load Balancer associated with the Target Group.
* `target_no_list` - The list of target number to bind to the target group.
