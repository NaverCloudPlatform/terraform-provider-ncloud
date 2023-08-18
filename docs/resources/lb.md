---
subcategory: "Load Balancer"
---


# Resource: ncloud_lb

Provides a Load Balancer resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage
```hcl
resource "ncloud_lb" "test" {
  name = "tf-lb-test"
  network_type = "PUBLIC"
  type = "APPLICATION"
  subnet_no_list = [ ncloud_subnet.test.subnet_no ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the load balancer.
* `description` - (Optional) The description of the load balancer.
* `network_type` - (Optional) The network type of load balancer to create. Accepted values: `PUBLIC` | `PRIVATE`. Default: `PUBLIC`.
* `idle_timeout` - (Optional) The time in seconds that the idle timeout. Valid only if the load balancer type is not `NETWORK`. Default: 60.
* `type` - (Required) The type of load balancer to create. Accepted values: `APPLICATION` | `NETWORK` | `NETWORK_PROXY`.
* `throughput_type` - (Optional) The performance type code of load balancer. Accepted values: `SMALL` | `MEDIUM` | `LARGE`. If the load balancer type is `NETWORK` and the load balancer network type is `PRIVATE`, only `SMALL` can be selected. Default: `SMALL`.
* `subnet_no_list` - (Required) A list of IDs in the associated Subnets.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of load balancer.
* `load_balancer_no` - The ID of load balancer (It is the same result as id).
* `domain` - Domain name of load balancer.
* `vpc_no` - The ID of the associated VPC.
* `ip_list` - A list of IP address of load balancer.

## Import

Individual load balancer can be imported using `LOAD_BALANCER_NO`.
For example, import a load balancer `17019658` like this:

```bash
$ terraform import ncloud_lb.my_lb 17019658
```