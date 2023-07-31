---
subcategory: "VPC"
---


# Data Source: ncloud_vpc

This module can be useful for getting detail of VPC created before, such as determining the CIDR block of that VPC.

## Example Usage

The following example shows how to take a VPC ID as a variable and obtain the data needed to create a subnet using this data source.

```hcl
variable "vpc_no" {}

data "ncloud_vpc" "selected" {
  id = var.vpc_no
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = data.ncloud_vpc.selected.id
  subnet         = cidrsubnet(data.ncloud_vpc.selected.ipv4_cidr_block, 8, 1)
  zone           = "KR-2"
  network_acl_no = data.ncloud_vpc.selected.default_network_acl_no
  subnet_type    = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific VPC to retrieve.
* `name` - (Optional) name of the specific VPC to retrieve
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_no` - The ID of VPC. (It is the same result as `id`)
* `ipv4_cidr_block` - The CIDR block for the association.
* `default_network_acl_no` - The ID of the network ACL created by default on VPC creation.
* `default_access_control_group_no` - The ID of the ACG created by default on VPC creation.
* `default_public_route_table_no` - The ID of the Public Route Table created by default on VPC creation.
* `default_private_route_table_no` - The ID of the Private Route Table created by default on VPC creation.
