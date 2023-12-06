---
subcategory: "Server"
---


# Resource: ncloud_access_control_group

Provides an ACG(Access Control Group) resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```hcl
resource "ncloud_vpc" "vpc" {
  ipv4_cidr_block = "10.4.0.0/16"
}
resource "ncloud_access_control_group" "acg" {
  name        = "my-acg"
  description = "description"
  vpc_no      = ncloud_vpc.vpc.id
}
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated VPC.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) Indicates whether to get default group only.

## Attributes Reference

* `id` - The ID of ACG(Access Control Group)
* `access_control_group_no` - The ID of ACG(Access Control Group) (It is the same result as `id`)
* `is_default` - Whether is default or not by VPC creation.

## Import

### `terraform import` command

* Access Control Group can be imported using the `id`. For example:

```console
$ terraform import ncloud_access_control_group.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Access Control Group using the `id`. For example:

```terraform
import {
  to = ncloud_access_control_group.rsc_name
  id = "12345"
}
```
