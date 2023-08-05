---
subcategory: "Developer Tools"
---


# Data Source: ncloud_sourcecommit_repositories

~> **Note:** This data source only supports 'public' site.

~> **Note:** This data source is a beta release. Some features may change in the future.

This data source is useful for look up the list of Sourcecommit repository in the region.

## Example Usage

In the example below, Retrieves all repositories with "test" in their names.

```hcl
data "ncloud_sourcecommit_repositories" "lookup-repos" {
  filter {
    name = "name"
    values = ["test"]
    regex = true
  }
}

output "lookup-repos-output" {
  value = data.ncloud_sourcecommit_repositories.lookup-repos.repositories
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.



## Attributes Reference

The following attributes are exported:

* `repositories` - The list of repositories.

### Repository Reference

`repositories` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - Sourcecommit repository ID.
* `repository_no` - Sourcecommit repository ID.
* `name` - Sourcecommit repository Name.
* `action_name`- Permission status for searching details.
* `permission`- Permission name for searching details. (`Allow` or `Deny`)