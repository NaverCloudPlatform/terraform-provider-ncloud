---
subcategory: "VPC"
---


# Data Source: ncloud_network_acl_deny_allow_groups

This resource is useful for look up the list of Network ACL Deny-Allow Group in the region.

## Example Usage

### Retrieve by Deny-Allow Group ID

```hcl
data "ncloud_network_acl_deny_allow_groups" "deny_allow_groups" {
  network_acl_deny_allow_group_no_list = [ncloud_network_acl_deny_allow_group.deny_allow_group.id]
}

resource "ncloud_network_acl_deny_allow_group" "deny_allow_group" {
  vpc_no         = ncloud_vpc.vpc.id
  ip_list = ["10.0.0.1", "10.0.0.2"]
}
```

### Retrieve by Specific VPC and name

```hcl
data "ncloud_network_acl_deny_allow_groups" "deny_allow_groups" {
  vpc_no = ncloud_vpc.vpc.id
  name   = "deny-allow-test"
}
```

### Retrieve by filter

```hcl
data "ncloud_network_acl_deny_allow_groups" "deny_allow_groups" {
  filter {
    name = "name"
    values = ["deny-allow-test"]
    regex = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_deny_allow_group_no_list` - (Optional) List of Deny-Allow Group ID to retrieve.
* `vpc_no` - (Optional) The ID of the specific VPC to retrieve.
* `name` - (Optional) name of the specific Deny-Allow Group to retrieve.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attributes are exported:

* `network_acl_deny_allow_groups` - The list of Deny-Allow Group

### Network ACL Deny Allow Group Reference

`network_acl_deny_allow_groups` are also exported with the following attributes, where are relevant: Each element
supports the following:

* `id` - The ID of Deny-Allow Group.
* `network_acl_deny_allow_group_no` - The ID of Deny-Allow Group. (It is the same result as `id`)
* `vpc_no` - The ID of the associated VPC.
* `ip_list` - list of IP address that registered in the Deny-Allow Group.
* `name` - The name of Deny-Allow Group.
* `description` - Description of Deny-Allow Group.