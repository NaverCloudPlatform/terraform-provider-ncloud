---
subcategory: "VPC"
---


# Data Source: ncloud_route_tables

This resource is useful for look up the list of Route table in the region.

## Example Usage

The example below shows how to make multiple private default route tables.

```hcl
data "ncloud_route_tables" "private_route_tables" {
  supported_subnet_type = "PRIVATE"
  filter {
    name = "is_default"
    values = ["true"]
  }
}

resource "ncloud_nat_gateway" "nat_gw" {
  count       = length(data.ncloud_route_tables.private_route_tables.route_tables)
  vpc_no      = data.ncloud_route_tables.private_route_tables.route_tables[count.index].vpc_no
  zone        = "KR-2"
}

resource "ncloud_route" "foo" {
  count                  = length(ncloud_nat_gateway.nat_gw)
  route_table_no         = data.ncloud_route_tables.private_route_tables.route_tables[count.index].id
  destination_cidr_block = "0.0.0.0/0"
  target_type            = "NATGW"
  target_name            = ncloud_nat_gateway.nat_gw[count.index].name
  target_no              = ncloud_nat_gateway.nat_gw[count.index].id
}

```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Optional) The ID of the specific VPC to retrieve.
* `supported_subnet_type` - (Optional) Subnet type. Accepted values : `PUBLIC` (Public) | `PRIVATE` (Private). 
* `name` - (Optional) The name of the specific Route Table to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attributes are exported:

* `route_tables` - The list of Route Tables

### Route Table Reference

`route_tables` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - The ID of Route Table.
* `route_table_no` - The ID of Route Table. (It is the same result as `id`)
* `vpc_no` - The ID of the associated VPC.
* `supported_subnet_type` - Subnet type. Accepted values : `PUBLIC` (Public) | `PRIVATE` (Private). 
* `name` - The name of Route Table.
* `description` - Description of Route Table.