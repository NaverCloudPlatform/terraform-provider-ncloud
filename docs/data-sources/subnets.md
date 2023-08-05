---
subcategory: "VPC"
---


# Data Source: ncloud_subnets

This resource is useful for look up the list of Subnet in the region.

## Example Usage

In the example below, Get CIDR block information of the subnets by specific VPC ID.

```hcl
variable "vpc_no" {
}

data "ncloud_subnets" "list_subnet" {
  vpc_no = var.vpc_no
}

output "subnet_cidr_list_by_vpc_no" {
  value = {
    for subnet in data.ncloud_subnets.list_subnet.subnets:
    subnet.id => subnet.subnet
  }
}
```

## Argument Reference

The following arguments are supported:

* `subnet_no` - (Optional) The ID of the specific Subnet to retrieve.
* `vpc_no` - (Optional) The ID of the specific VPC to retrieve.
* `subnet` - (Optional) The CIDR block of subnet to retrieve. 
* `zone` - (Optional) Available zone where the subnet will be placed physically.
* `network_acl_no` - (Optional) The ID of Network ACL.
* `subnet_type` - (Optional) Internet connectivity. If you use `PUBLIC`, all VMs created within Subnet will be assigned a certified IP by default and will be able to communicate directly over the Internet. Considering the characteristics of Subnet, you can choose Subnet for the purpose of use. Accepted values: `PUBLIC` (Public) | `PRIVATE` (Private).
* `usage_type` - (Optional) Usage type, Accepted values: `GEN` (General) | `LOADB` (For LoadBalancer) | `BM` (For BareMetal) |`NATGW` (for NATGateway).
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attributes are exported:

* `subnets` - The list of Subnet

### Subnet Reference

`subnets` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - The ID of Subnet.
* `subnet_no` - The ID of Subnet. (It is the same result as `id`)
* `vpc_no` - The ID of the associated VPC.
* `name` - The name of subnet.
* `subnet` - The CIDR block of subnet. 
* `zone` - Available zone where the Subnet is placed.
* `network_acl_no` - The ID of Network ACL.
* `subnet_type` - Internet connectivity.
* `usage_type` - Usage type.
