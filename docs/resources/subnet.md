---
subcategory: "VPC"
---


# Resource: ncloud_subnet

Provides a Subnet resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.1.0/24"
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC" // PUBLIC(Public) | PRIVATE(Private)
  // below fields is optional
  name           = "subnet-01"
  usage_type     = "GEN"    // GEN(General) | LOADB(For load balancer)
}
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the VPC where you want to place the Subnet.
* `subnet` - (Required) assign some subnet address ranges within the range of VPC addresses, must be between /16 and/28 within the private band (10.0.0/8,172.16.0.0/12,192.168.0.0/16).
* `zone` - (Required) Available zone where the subnet will be placed physically.
* `network_acl_no` - (Required) The ID of Network ACL.
* `subnet_type` - (Required) Internet connectivity. If you use `PUBLIC` all VMs created within Subnet will be assigned a certified IP by default and will be able to communicate directly over the Internet. Considering the characteristics of Subnet, you can choose Subnet for the purpose of use. Accepted values: `PUBLIC` (Public) | `PRIVATE` (Private).
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `usage_type` - (Optional) Usage type, Default `GEN`. Accepted values: `GEN` (General) | `LOADB` (For LoadBalancer) | `BM` (For BareMetal) |`NATGW` (for NATGateway).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of Subnet.
* `subnet_no` - The ID of the Subnet. (It is the same result as `id`)
* `vpc_no` - The ID of VPC. 
