---
subcategory: "Server"
---


# Resource: ncloud_placement_group

Provides a Placement Group resource.

## Example Usage

The example below shows how to create a placement group and apply a script when creating a server.

```terraform
variable "subnet_no" {}

resource "ncloud_placement_group" "group-a" {
  name = "plc-group-a"
}

data "ncloud_server_image_numbers" "kvm-image" {
  server_image_name = "rocky-8.10-base"
  filter {
    name   = "hypervisor_type"
    values = ["KVM"]
  }
}

data "ncloud_server_specs" "kvm-spec" {
  filter {
    name   = "server_spec_code"
    values = ["s2-g3"]
  }
}

resource "ncloud_server" "server" {
  subnet_no           = var.subnet_no
  server_image_number = data.ncloud_server_image_numbers.kvm-image.image_number_list.0.server_image_number
  server_spec_code    = data.ncloud_server_specs.kvm-spec.server_spec_list.0.server_spec_code
  placement_group_no  = ncloud_placement_group.group-a.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `placement_group_type` - (Optional) Type of placement group. Default `AA`. Accepted values: `AA`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Placement group.
* `placement_group_no` - The ID of the Placement group. (It is the same result as `id`)

## Import

### `terraform import` command

* Placement Group can be imported using the `id`. For example:

```console
$ terraform import ncloud_placement_group.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Placement Group using the `id`. For example:

```terraform
import {
  to = ncloud_placement_group.rsc_name
  id = "12345"
}
```
