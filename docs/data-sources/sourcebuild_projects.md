# Data Source: ncloud_sourcebuild_projects

~> **Note** This data source only supports 'public' site.

~> **Note:** This data source is a beta release. Some features may change in the future.

This data source is useful for look up the list of Sourcebuild project in the region.

## Example Usage

In the example below, Retrieves all Sourcebuild projects with "Allow" in their permissions.

```hcl
data "ncloud_sourcebuild_projects" "build_projects" {
  filter {
    name   = "permission"
    values = ["Allow"]
  }
}

output "lookup-build_projects-output" {
  value = data.ncloud_sourcebuild_projects.build_projects.projects
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Optional) Search by project name (project including string).
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `projects` - Sourcebuild projects.

### Project Reference

`projects` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - Sourcebuild project ID.
* `name` - Sourcebuild project Name.
* `action_name` - Permission status for searching details.
* `permission` - Permission name for searching details. (`Allow` or `Deny`)
