# Resource: ncloud_nat_gateway

Provides a NAT gateway resource.

## Example Usage

### Old Version(Subnet not associated) Usage

```hcl
resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.id
  zone        = "KR-2"
  // below fields is optional
  name        = "nat-gw"
  description = "description"
}

```

### New Version(Subnet associated) Usage

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
  usage_type     = "NATGW"
}

resource "ncloud_nat_gateway" "nat_gateway" {
  vpc_no      = ncloud_vpc.vpc.id
  subnet_no   = ncloud_subnet.subnet.id
  zone        = "KR-2"
  // below fields is optional
  name        = "nat-gw"
  description = "description"
}

```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `zone` - (Required) Available zone where the subnet will be placed physically.
* `subnet_no` - (Conditional) The ID of the associated SUBNET. This is required when creating a new one. The subnet type determines whether the NATGateway type is public or private. 
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `public_ip_no` - (Optional) The ID of the associated Public IP. If omitted, will auto create. should only be set if public type.
* `private ip` - (Optional) Private IP on created NAT Gateway. If omitted, will auto create.
* `description` - (Optional) description to create.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the NAT Gateway.
* `nat_gateway_no` - The ID of the NAT Gateway. (It is the same result as `id`) 
* `public_ip` - Public IP on created NAT Gateway.
* `subnet_name` - Subnet name on created NAT Gateway.

## Import

NAT Gateway can be imported using the id, e.g.,

$ terraform import ncloud_nat_gateway.my_nat_gateway id
