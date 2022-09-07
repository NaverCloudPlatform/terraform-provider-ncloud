# Data Source: ncloud_sourcedeploy_project_stage_scenarios

~> **Note** This data source only supports 'public' site.

-> **Note:** This data source is a beta release. Some features may change in the future.

This resource is useful for look up the list of Sourcedeploy scenario in the region.

## Example Usage

In the example below, specific Sourcedeploy scenario name.

```hcl
data "ncloud_sourcedeploy_project_stage_scenarios" "deploy_scenarios" {
  project_id = 1234
  stage_id   = 1234
  filter{
    name   = "name"
    values = ["test"]
  }
}

output "deploy_scenarios_output" {
  value = data.ncloud_sourcedeploy_project_stage_scenarios.deploy_scenarios.scenarios
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of Sourcedeploy project.
* `stage_id` - (Required) The ID of Sourcedeploy stage.

* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.


## Attributes Reference

The following attributes are exported:

* `scenarios` - The list of Sourcedeploy scenario.

### Scenarios Reference

`scenarios` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - The ID of Sourcedeploy scenario.
* `name` - The name of Sourcedeploy scenario.
