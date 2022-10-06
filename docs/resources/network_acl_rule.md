# Resource: ncloud_network_acl_rule

Provides a rule of Network ACL  resource.

~> **NOTE:** Do not create Network ACL Rule Resource more than once for the same Network ACL no. Doing so will cause a conflict of rule settings and occur an error, and rule settings will be deleted or overwritten.

~> **NOTE:** To performance, we recommend using one resource per Network ACL. If you use multiple resources in a single Network ACL, then can cause a slowdown.

## Example Usage

### Basic

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

### Network ACL Deny-Allow Group

```hcl
resource "ncloud_network_acl_deny_allow_group" "deny_allow_group" {
  vpc_no         = ncloud_vpc.vpc.id
  // below fields is optional
  name      = "deny-allow-group-test" 
  description = "by terraform"
  ip_list = ["10.0.0.1", "10.0.0.2"]
}

resource "ncloud_vpc" "vpc" {
   ipv4_cidr_block = "10.0.0.0/16"
 }
 
resource "ncloud_network_acl" "nacl" {
   vpc_no      = ncloud_vpc.vpc.id
}
 
resource "ncloud_network_acl_rule" "nacl_rule" {
  network_acl_no    = ncloud_network_acl.nacl.id

  inbound {
    priority    = 110
    protocol    = "TCP"
    rule_action = "ALLOW"
    deny_allow_group_no = ncloud_network_acl_deny_allow_group.deny_allow_group.id
    port_range  = "22"
  }
}

```

## Argument Reference

The following arguments are supported:

* `network_acl_no` - (Required) The ID of the Network ACL.
* `inbound` - (Optional) Specifies an Inbound(ingress) rules. Parameters defined below. This argument is processed
  in [attriutbe-as-blocks](https://www.terraform.io/docs/configuration/attr-as-blocks.html) mode.
* `outbound` - (Optional) Specifies an Outbound(egress) rules. Parameters defined below. This argument is processed
  in [attriutbe-as-blocks](https://www.terraform.io/docs/configuration/attr-as-blocks.html) mode.

### Network ACL Rule Reference

Both `inbound` and `outbound` support  following attributes:

* `priority` - (Required) Priority for rules, Used for ordering. Can be an integer from `1` to `199`.
* `protocol` - (Required) Select between TCP, UDP, and ICMP. Accepted values: `TCP` | `UDP` | `ICMP`
* `rule_action` - (Required) The action to take. Accepted values: `ALLOW` | `DROP`
* `ip_block` - (Optional, Required if `deny_allow_group_no` is not provided) The CIDR block to match. This must be a
  valid network mask.
* `deny_allow_group_no` - (Optional, Required if `ip_block` is not provided) The access source Deny-Allow Group number
  of network ACL rules. You can specify a Deny-Allow Group instead of an IP address block as the access
  source. `deny_allow_group_no` can be obtained through the Data source `ncloud_network_acl_deny_allow_group`
* `port_range` - (Optional) Range of ports to apply. You can enter from `1` to `65535`. e.g. set single port: `22` or
  set range port : `8000-9000`

~> **NOTE:** If the value of protocol is `ICMP`, the `port_range` values will be ignored and the rule will apply to all ports.

* `description` - (Optional) description to create.