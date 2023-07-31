---
subcategory: "VPC"
---


# Resource: ncloud_route

Provides a Route resource.

## Example Usage

### Usage with NAT Gateway

```hcl
resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_route_table" "route_table" {
  vpc_no                = ncloud_vpc.vpc.id  
  supported_subnet_type = "PUBLIC"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no = ncloud_vpc.vpc.id
  zone   = "KR-2"
}

resource "ncloud_route" "foo" {
  route_table_no         = ncloud_route_table.route_table.id
  destination_cidr_block = "0.0.0.0/0"
  target_type            = "NATGW"  // NATGW (NAT Gateway) | VPCPEERING (VPC Peering) | VGW (Virtual Private Gateway).
  target_name            = ncloud_nat_gateway.nat_gateway.name
  target_no              = ncloud_nat_gateway.nat_gateway.id
}
```

## Argument Reference

The following arguments are supported:

* `route_table_no` - (Required) The ID of the Route table.
* `destination_cidr_block` - (Required) Destination CIDR block, Set the destination IP address range for the route you want to add. (e.g. 0.0.0.0/0, 100.10.20.0/24) 
* `target_type` - (Required) Destination target type, Select the destination type of the route to add. Accepted values: `NATGW` (NAT Gateway) | `VPCPEERING` (VPC Peering) | `VGW` (Virtual Private Gateway).
* `target_no` - (Required) Set the destination identification number for the destination type.
* `target_name` - (Required) Set the destination name for the destination type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `is_default` - Whether is default or not by Route table creation.
* `vpc_no` - The ID of the associated VPC.

## Import

Individual routes can be imported using ROUTE_TABLE_NO:DESTINATION_CIDR. For example, import a route in a route table `57039` with an IPv4 destination CIDR `0.0.0.0/0` like this:

``` hcl
$ terraform import ncloud_route.my_route 57039:0.0.0.0/0
```