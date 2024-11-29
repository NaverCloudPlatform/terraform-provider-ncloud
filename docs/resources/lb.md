---
subcategory: "Load Balancer"
---


# Resource: ncloud_lb

Provides a Load Balancer resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
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
* `type` - (Required) The type of load balancer to create. Accepted values: `APPLICATION` | `NETWORK` | `NETWORK_PROXY`.
* `subnet_no_list` - (Required) A list of IDs in the associated Subnets.
* `network_type` - (Optional) The network type of load balancer to create. Accepted values: `PUBLIC` | `PRIVATE`. Default: `PUBLIC`.
* `idle_timeout` - (Optional) The time in seconds that the idle timeout. Valid only if the load balancer type is not `NETWORK`. Default: 60.
* `throughput_type` - (Optional) The performance type code of load balancer. `SMALL` | `MEDIUM` | `LARGE` | `DYNAMIC` | `XLARGE`. If the `type` is `APPLICATION` or `NETWORK_PROXY` Options : `SMALL` | `MEDIUM` | `LARGE` | `XLARGE`, Default : `SMALL`. If the `type` is `NETWORK` Options : `DYNAMIC`, Default : `DYNAMIC`.
* `description` - (Optional) The description of the load balancer.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of load balancer.
* `load_balancer_no` - The ID of load balancer (It is the same result as id).
* `domain` - Domain name of load balancer.
* `vpc_no` - The ID of the associated VPC.
* `ip_list` - A list of IP address of load balancer.

## Import

### `terraform import` command

* Load Balancer can be imported using the `id`. For example:

```console
$ terraform import ncloud_lb.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Load Balancer using the `id`. For example:

```terraform
import {
  to = ncloud_lb.rsc_name
  id = "12345"
}
```
