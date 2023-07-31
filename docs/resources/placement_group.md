---
subcategory: "Server"
---


# Resource: ncloud_placement_group

Provides a Placement group resource.

## Example Usage

The example below shows how to create a placement group and apply a script when creating a server.

```hcl
variable "subnet_no" {}

resource "ncloud_placement_group" "group-a" {
  name = "plc-group-a"
}

resource "ncloud_server" "server" {
  subnet_no                 = var.subnet_no
  server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
  placement_group_no        = ncloud_placement_group.group-a.id
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