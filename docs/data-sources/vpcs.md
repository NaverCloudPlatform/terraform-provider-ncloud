---
subcategory: "VPC"
---


# Data Source: ncloud_vpcs

This resource is useful for look up the list of VPC in the region.

## Example Usage

```hcl
data "ncloud_vpcs" "list_vpc" {
}

output "cidr_list" {
  value = {
    for vpc in data.ncloud_vpcs.list_vpc.vpcs:
    vpc.id => vpc.ipv4_cidr_block
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - The ID of the specific VPC to retrieve.
* `name` - (Optional) name of the specific VPC to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attributes are exported:

* `vpcs` - The list of vpcs

### VPC Reference

`vpcs` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - The ID of VPC.
* `vpc_no` - The ID of VPC. (It is the same result as `id`)
* `name` - The name of VPC.
* `ipv4_cidr_block` - The CIDR block for the association.
* `default_network_acl_no` - The ID of the network ACL created by default on VPC creation.
* `default_access_control_group_no` - The ID of the ACG created by default on VPC creation.