---
subcategory: "Developer Tools"
---


# Data Source: ncloud_sourcebuild_project_docker_engines

~> **Note** This data source only supports 'public' site.

~> **Note:** This data source is a beta release. Some features may change in the future.

This data source is useful for look up the list of Sourcebuild docker engine in the region.

## Example Usage

In the example below, Retrieves all docker engines with "Docker:18.09.1" in their names.

```hcl
data "ncloud_sourcebuild_project_docker_engines" "docker_engines" {
  filter {
    name   = "name"
    values = ["Docker:18.09.1"]
  }
}

output "lookup-docker_engines-output" {
  value = data.ncloud_sourcebuild_project_docker_engines.docker_engines.docker_engines
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `docker_engines` - Docker engines available at Sourcebuild.

### Docker Engines Reference

`docker_engines` is also exported with the following attributes, where relevant: Each element supports the following:

* `id` - Docker engine ID.
* `name` - Docker engine name.
