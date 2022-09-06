# Data Source: ncloud_sourcedeploy_projects

~> **Note** This data source only supports 'public' site.

-> **Note:** This data source is a beta release. Some features may change in the future.

This data source is useful for look up the list of Sourcedeploy project in the region.

## Example Usage

In the example below, Retrieves all Sourcedeploy projects.

```hcl
data "ncloud_sourcedeploy_projects" "deploy_projects"{
  filter {
    name    = "name"
    values  = ["test"]
  }
}

output "deploy_projects_output"{
  value = data.ncloud_Sourcedeploy_projects.data_projects.projects
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

* `projects` - The list of Sourcedeploy project.

### Project Reference

`projects` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - The ID of Sourcedeploy project.
* `name` - The name of Sourcedeploy project.
