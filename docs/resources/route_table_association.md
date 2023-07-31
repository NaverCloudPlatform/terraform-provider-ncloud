---
subcategory: "VPC"
---


# Resource: ncloud_route_table_association

Provide a resource to create association between route table and a subnet.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "vpc" {
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.id
	subnet             = "10.3.1.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}

resource "ncloud_route_table" "route_table" {
	vpc_no                = ncloud_vpc.vpc.id
	description           = "for test"
	supported_subnet_type = "PUBLIC"
}

resource "ncloud_route_table_association" "route_table_subnet" {
	route_table_no        = ncloud_route_table.route_table.id
	subnet_no             = ncloud_subnet.subnet.id
}
```

## Argument Reference

The following arguments are supported:

* `route_table_no` - (Required) The ID of the Route Table.
* `subnet_no` - (Required) The ID of Subnet to create association.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the route table association (`route_table_no`:`subnet_no`)
