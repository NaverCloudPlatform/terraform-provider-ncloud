---
layout: "ncloud"
page_title: "NCLOUD: ncloud_network_acl"
sidebar_current: "docs-ncloud-resource-network-acl"
description: |-
  Provides a Network ACL resource.
---

# Resource: ncloud_network_acl

Provides a Network ACL resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "vpc" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
  vpc_no      = ncloud_vpc.vpc.id
  // below fields is optional
  name        = "main"
  description = "for test"

  inbound {
    priority    = 100
    protocol    = "TCP"
    rule_action = "ALLOW"
    ip_block    = "0.0.0.0/0"
    port_range  = "22"
  }

  inbound {
    priority    = 110
    protocol    = "TCP"
    rule_action = "ALLOW"
    ip_block    = "0.0.0.0/0"
    port_range  = "80"
  }

  outbound {
    priority    = 100
    protocol    = "TCP"
    rule_action = "ALLOW"
    ip_block    = "0.0.0.0/0"
    port_range  = "1-65535"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create
* `inbound` - (Optional) Set Inbound(ingress) rules
  * `priority` - (Required) Priority for rules, Used for ordering. Can be an integer from `1` to `199`.
  * `protocol` - (Required) Select between TCP, UDP, and ICMP. Accepted values: `TCP` | `UDP` | `ICMP`
  * `rule_action` - (Required) The action to take. Accepted values: `ALLOW` | `DROP`
  * `ip_block` - (Required) The CIDR block to match. This must be a valid network mask.
  * `port_range` - (Optional) Range of ports to apply. You can enter from `1` to `65535`. e.g. set single port: `22` or set range port : `8000-9000`
* `outbound` - (Optional) Set Outbound(egress) rules
  * `priority` - (Required) Priority for rules, Used for ordering. Can be an integer from `1` to `199`.
  * `protocol` - (Required) Select between TCP, UDP, and ICMP. Accepted values: `TCP` | `UDP` | `ICMP`
  * `rule_action` - (Required) The action to take. Accepted values: `ALLOW` | `DROP`
  * `ip_block` - (Required) The CIDR block to match. This must be a valid network mask.
  * `port_range` - (Optional) Range of ports to apply. You can enter from `1` to `65535`. e.g. set single port: `22` or set range port : `8000-9000`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Network ACL.
* `network_acl_no` - The ID of the Network ACL. (It is the same result as `id`)
* `is_default` - Whether is default or not by VPC creation.