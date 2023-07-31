---
subcategory: "VPC"
---


# Resource: ncloud_route_table

Provides a Route Table resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_route_table" "route_table" {
  vpc_no                = ncloud_vpc.vpc.id  
  supported_subnet_type = "PUBLIC" // PUBLIC | PRIVATE
  // below fields is optional
  name                  = "route-table"
  description           = "for test"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `supported_subnet_type` - (Required) Subnet type. Accepted values : `PUBLIC` (Public) | `PRIVATE` (Private). 
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Route table.
* `route_table_no` - The ID of the Route table. (It is the same result as `id`)
* `is_default` - Whether is default or not by VPC creation.