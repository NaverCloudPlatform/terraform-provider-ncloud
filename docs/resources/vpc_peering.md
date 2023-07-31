---
subcategory: "VPC"
---


# Resource: ncloud_vpc_peering

Provides a VPC Peering resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "main" {
  name            = "vpc-main"
  ipv4_cidr_block = "10.4.0.0/16"
}

resource "ncloud_vpc" "peer" {
  name            = "vpc-peer"
  ipv4_cidr_block = "10.5.0.0/16"
}

resource "ncloud_vpc_peering" "foo" {
  name          = "vpc_peering_example"
  source_vpc_no = ncloud_vpc.main.id
  target_vpc_no = ncloud_vpc.peer.id
}
```

## Argument Reference

The following arguments are supported:

* `source_vpc_no` - (Required) The ID of VPC from which the request is sent.
* `target_vpc_no `- (Required) The ID of VPC to receive requests.
* `target_vpc_name `- (Optional) The name of the VPC that receives the request.
* `target_vpc_login_id `- (Optional) VPC Owner ID to receive requests (If the account receiving the request is different from the account you send, you must enter the account receiving the request. Must match E-mail format).
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of VPC peering.
* `vpc_peering_no` - The ID of VPC peering. (It is the same result as `id`)
* `has_reverse_vpc_peering` - Reverse VPC Peering exists.
* `is_between_accounts` - VPC Peering Between Accounts.