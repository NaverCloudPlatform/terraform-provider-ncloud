---
subcategory: "Load Balancer"
---


# Resource: ncloud_lb_target_group_attachment

Provides a Target Group Attachment resource.

## Example Usage
```hcl
resource "ncloud_server" "test" {
  # ...
}

resource "ncloud_lb_target_group" "test" {
  # ...
}

resource "ncloud_lb_target_group_attachment" "test" {
  target_group_no = ncloud_lb_target_group.test.target_group_no
  target_no_list = [ncloud_server.test.instance_no]
}
```

## Argument Reference

The following arguments are supported:

* `target_group_no` - (Required) The ID of target group.
* `target_no_list` - (Required) The List of server instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of target.