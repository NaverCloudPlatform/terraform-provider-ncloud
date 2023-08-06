---
subcategory: "Server"
---


# Resource: ncloud_access_control_group_rule

Provides an rule of ACG(Access Control Group) resource.

~> **NOTE:** This resource only supports VPC environment.

~> **NOTE:** Do not create multiple ACG(Access Control Group) Rule resources and set them to a single ACG, as only one ACG Rule will be applied to a single ACG and may behave differently than expected, causing the rule to be overwritten.

## Example Usage

```hcl
resource "ncloud_vpc" "vpc" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_access_control_group" "acg" {
  name        = "my-acg"
  description = "description"
  vpc_no      = ncloud_vpc.vpc.id
}

resource "ncloud_access_control_group_rule" "acg-rule" {
  access_control_group_no = ncloud_access_control_group.acg.id
  
  inbound {
    protocol    = "TCP"
    ip_block    = "0.0.0.0/0"
    port_range  = "22"
    description = "accept 22 port"
  }

  inbound {
    protocol    = "TCP"
    ip_block    = "0.0.0.0/0"
    port_range  = "80"
    description = "accept 80 port"
  }

  outbound {
    protocol    = "TCP"
    ip_block    = "0.0.0.0/0" 
    port_range  = "1-65535"
    description = "accept 1-65535 port"
  }
}
```

## Argument Reference

~> **NOTE:** One of either `ip_block` or `source_access_control_group_no` is required.

The following arguments are supported:

* `access_control_group_no` - (Required) The ID of the ACG.
* `inbound` - (Optional) Specifies an Inbound(ingress) rules. Parameters defined below. This argument is processed in [attriutbe-as-blocks](https://www.terraform.io/docs/configuration/attr-as-blocks.html) mode.
* `outbound` - (Optional) Specifies an Outbound(egress) rules. Parameters defined below. This argument is processed in [attriutbe-as-blocks](https://www.terraform.io/docs/configuration/attr-as-blocks.html) mode.

### Access Control Group Rule Reference

Both `inbound` and `outbound` support  following attributes:

* `protocol` - (Required) Select between TCP, UDP, and ICMP. Accepted values: `TCP` | `UDP` | `ICMP`
* `ip_block` - (Optional) The CIDR block to match. This must be a valid network mask. Cannot be specified with `source_access_control_group_no`.
* `source_access_control_group_no` - (Optional) The ID of specific ACG to apply this rule to. Cannot be specified with `ip_block`.
* `port_range` - (Optional) Range of ports to apply. You can enter from `1` to `65535`. e.g. set single port: `22` or set range port : `8000-9000`

~> **NOTE:** If the value of protocol is `ICMP`, the `port_range` values will be ignored and the rule will apply to all ports.

* `description` - (Optional) description to create.

## Attributes Reference

* `id` - The ID of ACG(Access Control Group) rule
