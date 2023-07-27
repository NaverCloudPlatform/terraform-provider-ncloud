---
subcategory: "VPC"
---


# Resource: ncloud_network_acl

Provides a Network ACL resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_vpc" "vpc" {
   ipv4_cidr_block = "10.0.0.0/16"
 }
 
resource "ncloud_network_acl" "nacl" {
   vpc_no      = ncloud_vpc.vpc.id
  // below fields is optional
   name        = "main"
   description = "for test"
 }
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Network ACL.
* `network_acl_no` - The ID of the Network ACL. (It is the same result as `id`)
* `is_default` - Whether is default or not by VPC creation.

## Import

Network ACL can be imported using the id, e.g.,

$ terraform import ncloud_network_acl.my_nacl id
