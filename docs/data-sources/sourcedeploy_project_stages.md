# Data Source: ncloud_sourcedeploy_project_stages

~> **Note** This data source only supports 'public' site.

-> **Note:** This data source is a beta release. Some features may change in the future.

This resource is useful for look up the list of Sourcedeploy stage in the region.

## Example Usage

In the example below, specific Sourcedeploy stages name.

```hcl
data "ncloud_sourcedeploy_project_stages" "deploy_stages" {
  project_id = 1234
  filter {
    name    = "name"
    values  = ["asg"]
  }
}

output "output_stages"{
  value = data.ncloud_sourcedeploy_project_stages.deploy_stages.stages
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of Sourcedeploy project.

* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attributes are exported:

* `stages` - The list of Sourcedeploy stage.

### Stage Reference

`stages` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - The ID of Sourcedeploy stage.
* `name` - The name of Sourcedeploy stage.
