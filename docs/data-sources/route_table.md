---
subcategory: "VPC"
---


# Data Source: ncloud_route_table

This module can be useful for getting detail of Route Table created before.

## Example Usage

### Basic Usage

```hcl
variable "route_table_no" {}

data "ncloud_route_table" "selected" {
  id = var.route_table_no
}
```

### Usage of using filter

The example below is an example of using filters to select a default private route table and to connect with a NAT gateway.

```hcl
variable "vpc_no" {}

data "ncloud_route_table" "selected" {
  vpc_no                = var.vpc_no
  supported_subnet_type = "PUBLIC"
  filter {
    name = "is_default"
    values = ["true"]
  }
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no = var.vpc_no
  zone   = "KR-2"
}

resource "ncloud_route" "foo" {
  route_table_no         = data.ncloud_route_table.selected.id
  destination_cidr_block = "0.0.0.0/0"
  target_type            = "NATGW"
  target_name            = ncloud_nat_gateway.nat_gateway.name
  target_no              = ncloud_nat_gateway.nat_gateway.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific Route Table to retrieve.
* `vpc_no` - (Optional) The ID of the specific VPC to retrieve.
* `supported_subnet_type` - (Optional) Subnet type. Accepted values : `PUBLIC` (Public) | `PRIVATE` (Private). 
* `name` - (Optional) name of the specific Route Table to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `route_table_no` - The ID of Route Table. (It is the same result as `id`)
* `description` - Description of Route Table.
* `is_default` - Whether is default or not by VPC creation.