---
layout: "ncloud"
page_title: "NCLOUD: ncloud_nat_gateway"
sidebar_current: "docs-ncloud-resource-nat-gateway"
description: |-
  Provides a NAT gateway resource.
---

# Resource: ncloud_nat_gateway

Provides a NAT gateway resource.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `zone` - (Required) Available zone where the subnet will be placed physically.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the NAT Gateway.
* `nat_gateway_no` - The ID of the NAT Gateway. (It is the same result as `id`) 
* `public_ip` - Public IP on created NAT Gateway.