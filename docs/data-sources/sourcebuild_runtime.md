# Data Source: ncloud_sourcebuild_runtime

This data source is useful for look up the list of Sourcebuild runtime environment in the region.

## Example Usage

In the example below, Retrieves all Sourcebuild runtime environments with the os name is "ubuntu" and runtime name is "base".

```hcl
data "ncloud_sourcebuild_os" "os" {
  filter {
    name   = "name"
    values = ["ubuntu"]
  }
}

data "ncloud_sourcebuild_runtime" "runtime" {
  os_id = data.ncloud_sourcebuild_os.os.os[0].id

  filter {
    name   = "name"
    values = ["base"]
  }
}

output "lookup-runtime-output" {
  value = data.ncloud_sourcebuild_runtime.runtime.runtime
}
```

## Argument Reference

The following arguments are supported:

* `os_id` - (Required) OS ID which runtime belongs.
    * [`ncloud_sourcebuild_os` data source](./data-sources/sourcebuild_os.md)
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `runtime` - Runtimes available at Sourcebuild.

### Runtime Reference

`runtime` are is exported with the following attributes, where relevant: Each element supports the following:

* `id` - Runtime ID.
* `name` - Runtime name.