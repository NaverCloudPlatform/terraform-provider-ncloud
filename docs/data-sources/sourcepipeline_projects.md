# Data Source: ncloud_sourcepipieline_projects

This data source is useful for look up the list of Sourcepipeline projects in the region.

## Example Usage

In the example below, Retrieves all Sourcepipeline projects.

```hcl
data "ncloud_sourcepipeline_projects" "list_sourcepipeline" {
}

output "sourcepipeline_list" {
    value = {
        for pipeline in data.ncloud_sourcepipeline_projects.list_sourcepipeline.projects :
        pipeline.id => pipeline.name
    }
}
```

## Argument Reference

The following arguments are supported:

*   `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
*   `filter` - (Optional) Custom filter block as described below.
    *   `name` - (Required) The name of the field to filter by.
    *   `values` - (Required) Set of values that are accepted for the given field.
    *   `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attributes are exported:

*   `projects` - The list of Sourcepipeline project

### Project Reference

`projects` are also exported with the following attributes, where relevant: Each element supports the following:

*   `id` - The ID of Sourcepipeline project.
*   `name` - The name Sourcepipeline project.
