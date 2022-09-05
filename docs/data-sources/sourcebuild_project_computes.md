# Data Source: ncloud_sourcebuild_project_computes

~> **Note** This data source only supports 'public' site.

~> **Note:** This data source is a beta release. Some features may change in the future.

This data source is useful for look up the list of Sourcebuild compute in the region.

## Example Usage

In the example below, Retrieves all compute environments with the number of cpu is '2'.

```hcl
data "ncloud_sourcebuild_project_computes" "computes" {
  filter {
    name   = "cpu"
    values = [2]
  }
}

output "lookup-computes-output" {
  value = data.ncloud_sourcebuild_project_computes.computes.computes
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `computes` - Computing environments available at Sourcebuild.

### Computes Reference

`computes` is also exported with the following attributes, where relevant: Each element supports the following:

* `id` - Compute ID.
* `cpu` - CPU of build environment.
* `mem` - Memory of build environment.
