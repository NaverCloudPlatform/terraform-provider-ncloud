---
subcategory: "VPC"
---


# Data Source: ncloud_nat_gateway

This module can provide useful for get detail of NAT Gateway created before.

## Example Usage

The example below is how to register a rule on Route table using an existing NAT Gateway.

```hcl
variable "nat_gateway_no" {}

data "ncloud_nat_gateway" "nat_gateway" {
  id = var.nat_gateway_no
}

resource "ncloud_route_table" "route_table" {
  vpc_no                = data.ncloud_nat_gateway.nat_gateway.vpc_no  
  supported_subnet_type = "PUBLIC"
}

resource "ncloud_route" "foo" {
  route_table_no         = data.ncloud_nat_gateway.nat_gateway.id
  destination_cidr_block = "0.0.0.0/0"
  target_type            = "NATGW"
  target_name            = data.ncloud_nat_gateway.nat_gateway.name
  target_no              = data.ncloud_nat_gateway.nat_gateway.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific NAT gateway to retrieve.
* `name` - (Optional) The name of the specific NAT gateway to retrieve.
* `vpc_name` - (Optional) name of the specific associated VPC to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nat_gateway_no` - The ID of NAT gateway. (It is the same result as `id`)
* `vpc_no` - The ID of the associated VPC.
* `subnet_no` - The ID of the associated Subnet.
* `subnet_name` - The name of the associated Subnet.
* `zone` - Available zone where the NAT gateway placed.
* `public_ip` - Public IP on NAT Gateway created.
* `public_ip_no` - The ID of the associated Public IP.
* `private_ip` - Private IP on NAT Gateway created.
* `description` - Description of NAT Gateway.
