# Data Source: ncloud_sourcebuild_os

This data source is useful for look up the list of Sourcebuild os in the region.

## Example Usage

In the example below, Retrieves all os environments with "ubuntu" in their names.

```hcl
data "ncloud_sourcebuild_os" "os" {
  filter {
    name   = "name"
    values = ["ubuntu"]
  }
}

output "lookup-os-output" {
  value = data.ncloud_sourcebuild_os.os.os
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

* `os` - OS available at Sourcebuild.

### OS Reference

`os` is also exported with the following attributes, where relevant: Each element supports the following:

* `id` - OS ID.
* `name` - OS name.
* `archi` - OS architecture.
* `version` - OS version.