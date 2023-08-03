---
subcategory: "VPC"
---


# Resource: ncloud_vpc

Provides a VPC resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "vpc" {
 ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
 vpc_no = ncloud_vpc.vpc.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `ipv4_cidr_block` - (Required) The CIDR block of the VPC. The range must be between /16 and/28 within the private band (10.0.0/8,172.16.0.0/12,192.168.0.0/16).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of VPC.
* `vpc_no` - The ID of VPC. (It is the same result as `id`)
* `default_network_acl_no` - The ID of the network ACL created by default on VPC creation.
* `default_access_control_group_no` - The ID of the ACG created by default on VPC creation.
* `default_public_route_table_no` - The ID of the Public Route Table created by default on VPC creation.
* `default_private_route_table_no` - The ID of the Private Route Table created by default on VPC creation.