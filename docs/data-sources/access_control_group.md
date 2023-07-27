---
subcategory: "Server"
---


# Data Source: ncloud_access_control_group

When creating a server instance (VM), you can add an ACG(Access Control Group) that you specified to set firewalls. `ncloud_access_control_group` provides details about a specific ACG(Access Control Group) information.

## Example Usage

### Basic usage

```hcl
variable "access_control_group_no" {}

data "ncloud_access_control_group" "selected" {
	id = var.access_control_group_no
}
```

### Search by ACG name

```hcl
data "ncloud_access_control_group" "selected" {
  name = "acg-name"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) List of ACG ID you want to get
* `name` - (Optional) Name of the ACG you want to get
* `vpc_no` - (Optional) The ID of the specific VPC to retrieve.
* `is_default` - (Optional) Indicates whether to get default group only
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `access_control_group_no` - The ID of ACG. (It is the same result as `id`)
* `description` - ACG description
