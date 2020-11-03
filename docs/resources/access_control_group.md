---
layout: "ncloud"
page_title: "NCLOUD: ncloud_access_control_group"
sidebar_current: "docs-ncloud-resource-access-control-group"
description: |-
  Provides an ACG(Access Control Group) resource.
---

# Resource: ncloud_access_control_group

Provides an ACG(Access Control Group) resource.

~> **NOTE:** This resource only support VPC environment.

## Example Usage

```hcl
resource "ncloud_vpc" "vpc" {
  ipv4_cidr_block = "10.4.0.0/16"
}
resource "ncloud_access_control_group" "acg" {
  name        = "my-acg"
  description = "description"
  vpc_no      = ncloud_vpc.vpc.id

  inbound {
    protocol    = "TCP"
    ip_block    = "0.0.0.0/0"
    port_range  = "22"
  }

  inbound {
    protocol    = "TCP"
    ip_block    = "0.0.0.0/0"
    port_range  = "80"
  }

  outbound {
    protocol    = "TCP"
    ip_block    = "0.0.0.0/0" 
    port_range  = "1-65535"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) Indicates whether to get default group only
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

* `id` - The ID of ACG(Access Control Group)
* `access_control_group_no` - The ID of ACG(Access Control Group) (It is the same result as `id`)
* `is_default` - Whether is default or not by VPC creation.