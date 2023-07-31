---
subcategory: "VPC"
---


# Resource: ncloud_network_acl_deny_allow_group

Provides a rule of Network ACL Deny-Allow Group resource. You can manage list of IP using this resource, \
Network ACL Deny-Allow Group can be added to the Network ACL Rule(`ncloud_network_acl_rule`) using `deny_allow_group_no`.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_network_acl_deny_allow_group" "deny_allow_group" {
  vpc_no      = ncloud_vpc.vpc.id
  // below fields is optional
  name        = "deny-allow-group-test" 
  description = "by terraform"
  ip_list     = ["10.0.0.1", "10.0.0.2"]
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
    priority            = 110
    protocol            = "TCP"
    rule_action         = "ALLOW"
    deny_allow_group_no = ncloud_network_acl_deny_allow_group.deny_allow_group.id
    port_range          = "22"
  }
}

```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `ip_list` - (Required) Enter the IP addresses as list to be registered in the Deny-Allow Group.
  Up to 100 IPs can be registered. Duplicate IP addresses are not allowed.
* `name` - (Optional) The name to create. If omitted, terraform will assign a random, unique name.
* `description` - (Optional) Description to create

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Deny-Allow Group.
* `network_acl_deny_allow_group_no` - The ID of the Deny-Allow Group. (It is the same result as `id`)