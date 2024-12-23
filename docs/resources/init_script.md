---
subcategory: "Server"
---


# Resource: ncloud_init_script

Provides an Init Script resource.

## Example Usage

The example below shows how to create an initial script and apply a script when creating a server.

```terraform
variable "subnet_no" {}

resource "ncloud_init_script" "init" {
  name    = "ls-script"
  content = "#!/usr/bin/env\nls -al"
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
  init_script_no      = ncloud_init_script.init.id
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required) Initialization script content. Scripts such as Python, Perl, Shell are available for Linux environments. However, on the first line, you must specify the script path you want to run in the form of `#!/usr/bin/env` python, `#!/bin/perl`, `#!/bin/bash`, etc. Windows environments can only write Visual Basic scripts.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create.
* `os_type` - (Optional) Type of O/S to apply server instance. Default `LNX`. Accepted values: `LNX` (LINUX) | `WND` (WINDOWS)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Init script.
* `init_script_no` - The ID of the Init script. (It is the same result as `id`)

## Import

### `terraform import` command

* Init Script can be imported using the `id`. For example:

```console
$ terraform import ncloud_init_script.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Init Script using the `id`. For example:

```terraform
import {
  to = ncloud_init_script.rsc_name
  id = "12345"
}
```
