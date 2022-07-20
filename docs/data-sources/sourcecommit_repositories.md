# Data Source: ncloud_sourcecommit_repositories

This resource is useful for look up the list of Sourcecommit repository in the region.

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

* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
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
* `name` - Sourcecommit repository Name.
* `action_name`- Permission status for searching details.
* `permission`- Permission name for searching details. (`Allow` or `Deny`)