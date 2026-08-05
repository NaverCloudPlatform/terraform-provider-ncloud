---
subcategory: "Load Balancer"
---


# Resource: ncloud_lb_target_group_attachment

Provides a Target Group Attachment resource.

~> **NOTE:** This resource only supports VPC environment.

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

* `id` - The ID of target group attachment. Format: `target_group_no:target_no[,target_no...]`

## Import

### `terraform import` command

* Load Balancer Target Group Attachment can be imported using the target group number and one or more target numbers separated by commas. For example:

```console
$ terraform import ncloud_lb_target_group_attachment.rsc_name 12345:23456,34567
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Load Balancer Target Group Attachment using the target group number and one or more target numbers separated by commas. For example:

```terraform
import {
  to = ncloud_lb_target_group_attachment.rsc_name
  id = "12345:23456,34567"
}
```
