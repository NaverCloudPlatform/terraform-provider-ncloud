# Resource: ncloud_network_acl_deny_allow_group

Provides a rule of Network ACL Deny-Allow Group resource. You can manage list of IP using this resource, \
Network ACL Deny Allow Group can be added to the Network ACL rule(ncloud_network_acl_rule) using `deny_allow_group_no`.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_network_acl_deny_allow_group" "allow_group" {
  vpc_no         = ncloud_vpc.vpc.id
  // below fields is optional
  name      = "allow-group-test" 
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
    deny_allow_group_no = ncloud_network_acl_deny_allow_group.allow_group.id
    port_range  = "22"
  }
}

```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create
* `ip_list` - (Optional) Enter the IP address to be registered in the Deny-Allow Group. Up to 100 IPs can be registered
  and duplicate IP addresses are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Deny-Allow Group.
* `network_acl_deny_allow_group_no` - The ID of the Deny-Allow Group. (It is the same result as `id`)