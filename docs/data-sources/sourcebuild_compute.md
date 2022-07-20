# Data Source: ncloud_sourcebuild_compute

This data source is useful for look up the list of Sourcebuild compute in the region.

## Example Usage

In the example below, Retrieves all compute environments with the number of cpu is '2'.

```hcl
data "ncloud_sourcebuild_compute" "compute" {
  filter {
    name   = "cpu"
    values = [2]
  }
}

output "lookup-compute-output" {
  value = data.ncloud_sourcebuild_compute.compute.compute
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

* `compute` - Computing environments available at Sourcebuild.

### Compute Reference

`compute` is also exported with the following attributes, where relevant: Each element supports the following:

* `id` - Compute ID.
* `cpu` - CPU of build environment.
* `mem` - Memory of build environment.