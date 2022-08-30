# Data Source: ncloud_sourcebuild_project_runtime_version

-> **Note:** This data source is a beta release. Some features may change in the future.

This data source is useful for look up the list of Sourcebuild runtime version environment in the region.

## Example Usage

In the example below, Retrieves all Sourcebuild runtime version environments with the os name is "ubuntu" and runtime is "base" and runtime version name is "16.04-1.0.0".

```hcl
data "ncloud_sourcebuild_project_os" "os" {
  filter {
    name   = "name"
    values = ["ubuntu"]
  }
}

data "ncloud_sourcebuild_project_runtime" "runtime" {
  os_id = data.ncloud_sourcebuild_project_os.os.os[0].id

  filter {
    name   = "name"
    values = ["base"]
  }
}

data "ncloud_sourcebuild_project_runtime_version" "runtime_version" {
  os_id      = data.ncloud_sourcebuild_project_os.os.os[0].id
  runtime_id = data.ncloud_sourcebuild_project_runtime.runtime.runtime[0].id

  filter {
    name   = "name"
    values = ["16.04-1.0.0"]
  }
}

output "lookup-runtime_version-output" {
  value = data.ncloud_sourcebuild_project_runtime_version.runtime_version.runtime_version
}
```

## Argument Reference

The following arguments are supported:

* `os_id` - (Required) OS ID which runtime belongs.
    * [`ncloud_sourcebuild_project_os` data source](./data-sources/sourcebuild_project_os.md)
* `runtime_id` - (Required) Runtime ID which runtime version belongs.
    * [`ncloud_sourcebuild_project_runtime` data source](./data-sources/sourcebuild_project_runtime.md)
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `runtime_version` - Runtime versions available at Sourcebuild.

### Runtime Version Reference

`runtime_version` is also exported with the following attributes, where relevant: Each element supports the following:

* `id` - Runtime version ID.
* `name` - Runtime version name.
