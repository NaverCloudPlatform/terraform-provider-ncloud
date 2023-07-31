---
subcategory: "Server"
---


# Resource: ncloud_init_script

Provides an Init script resource.

## Example Usage

The example below shows how to create an initial script and apply a script when creating a server.

```hcl
variable "subnet_no" {}

resource "ncloud_init_script" "init" {
  name    = "ls-script"
  content = "#!/usr/bin/env\nls -al"
}

resource "ncloud_server" "server" {
  subnet_no                 = var.subnet_no
  server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
  init_script_no            = ncloud_init_script.init.id
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
