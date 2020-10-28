---
layout: "ncloud"
page_title: "NCLOUD: ncloud_network_acl_rule"
sidebar_current: "docs-ncloud-resource-network-acl-rule"
description: |-
  Provides a network acl rule resource.
---

# Resource: ncloud_network_acl_rule

Provides a rule of Network ACL  resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "vpc" {
   ipv4_cidr_block = "10.0.0.0/16"
 }
 
resource "ncloud_network_acl" "nacl" {
   vpc_no      = ncloud_vpc.vpc.id
   name        = "main"
   description = "for test"
 }

resource "ncloud_network_acl_rule" "nacl_rule" {
  network_acl_no    = ncloud_network_acl.nacl.id
  network_rule_type = "INBND" // INBND | OTBND
  priority          = 100     // 1 to 199
  protocol          = "TCP"   // TCP | UDP | ICMP
  rule_action       = "ALLOW" // ALLOW | DROP
  ip_block          = "0.0.0.0/0"
  port_range        = "22"    // 1-65535
  // below fields is optional
  description       = "allow ssh port"
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_no` - (Required) The ID of the Network ACL.
* `priority` - (Required) Priority for rules, Used for ordering. Can be an integer from `1` to `199`.
* `protocol` - (Required) Select between TCP, UDP, and ICMP. Accepted values: `TCP` | `UDP` | `ICMP`
* `rule_action` - (Required) The action to take. Accepted values: `ALLOW` | `DROP`
* `ip_block` - (Required) The CIDR block to match. This must be a valid network mask.
* `rule_type` - (Required) Specifies an inbound(INBND) or outbound(OTBND) rule. Accepted values: `INBND` (inbound) | `OTBND` (outboud)
* `port_range` - (Optional) Range of ports to apply. You can enter from `1` to `65535`. e.g. set single port: `22` or set range port : `8000-9000`

~> **NOTE:** If the value of protocol is `ICMP`, the `port_range` values will be ignored and the rule will apply to all ports.

* `description` - (Optional) description to create