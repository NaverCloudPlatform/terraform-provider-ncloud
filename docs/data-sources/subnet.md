---
subcategory: "VPC"
---


# Data Source: ncloud_subnet

This module can be useful for getting detail of Subnet created before. for example, determine the CIDR block of that Subnet

## Example Usage

The following example shows how one might accept a subnet id as a variable and use this data source to obtain the data necessary to create a ACG that allows connections from hosts in that subnet.

```hcl
variable "subnet_no" {}

data "ncloud_subnet" "selected" {
  id = var.subnet_no
}

resource "ncloud_access_control_group" "acg" {
 vpc_no = data.ncloud_subnet.selected.vpc_no
}

resource "ncloud_access_control_group_rule" "subnet-inbound-tcp-80" {
  access_control_group_no = ncloud_access_control_group.acg.id
  rule_type               = "INBND"
  protocol                = "TCP"
  ip_block                = ncloud_subnet.selected.subnet
  port_range              = "80"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific Subnet to retrieve.
* `vpc_no` - (Optional) The ID of the specific VPC to retrieve.
* `subnet` - (Optional) The CIDR block of Subnet to retrieve. 
* `zone` - (Optional) Available zone where the subnet will be placed physically.
* `network_acl_no` - (Optional) The ID of Network ACL.
* `subnet_type` - (Optional) Internet connectivity. If you use `PUBLIC`, all VMs created within Subnet will be assigned a certified IP by default and will be able to communicate directly over the Internet. Considering the characteristics of Subnet, you can choose Subnet for the purpose of use. Accepted values: `PUBLIC` (Public) | `PRIVATE` (Private).
* `usage_type` - (Optional) Usage type, Accepted values: `GEN` (General) | `LOADB` (For LoadBalancer) | `BM` (For BareMetal) |`NATGW` (for NATGateway).
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `subnet_no` - The ID of Subnet. (It is the same result as `id`)
* `name` - The name of Subnet.
