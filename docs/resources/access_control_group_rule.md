---
layout: "ncloud"
page_title: "NCLOUD: ncloud_access_control_group"
sidebar_current: "docs-ncloud-resource-access-control-group"
description: |-
  Provides an rule of ACG(Access Control Group) resource.
---

# Resource: ncloud_access_control_group_rule

Provides an rule of ACG(Access Control Group) resource.

~> **NOTE:** This resource only support VPC environment.

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

resource "ncloud_access_control_group_rule" "inbound-tcp-22" {
  access_control_group_no = ncloud_access_control_group.acg.id
  description             = "accept 22 port"
  rule_type               = "INBND"
  protocol                = "TCP"
  ip_block                = "0.0.0.0/0"
  port_range              = "22"
}

resource "ncloud_access_control_group_rule" "inbound-tcp-8080" {
  access_control_group_no = ncloud_access_control_group.acg.id
  description             = "accept 8080 port"
  rule_type               = "INBND"
  protocol                = "TCP"
  ip_block                = "0.0.0.0/0"
  port_range              = "8080"
}
```

## Argument Reference

~> **NOTE:** One of either `ip_block` or `source_access_control_group_no` is required.

The following arguments are supported:

* `access_control_group_no` - (Required) The ID of the ACG.
* `rule_type` - (Required) Specifies an inbound(INBND) or outbound(OTBND) rule. Accepted values: `INBND` (inbound) | `OTBND` (outboud).
* `protocol` - (Required) Select between TCP, UDP, and ICMP. Accepted values: `TCP` | `UDP` | `ICMP`
* `ip_block` - (Optional) The CIDR block to match. This must be a valid network mask. Cannot be specified with `source_access_control_group_no`.
* `source_access_control_group_no` - (Optional) The ID of specific ACG to apply this rule to. Cannot be specified with `ip_block`.
* `port_range` - (Optional) Range of ports to apply. You can enter from `1` to `65535`. e.g. set single port: `22` or set range port : `8000-9000`

~> **NOTE:** If the value of protocol is `ICMP`, the `port_range` values will be ignored and the rule will apply to all ports.

* `description` - (Optional) description to create.