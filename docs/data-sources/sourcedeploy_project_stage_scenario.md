# Data Source: ncloud_sourcedeploy_project_stage_scenario

~> **Note** This data source only supports 'public' site.

-> **Note:** This data source is a beta release. Some features may change in the future.

This resource is useful for look up the list of Sourcedeploy scenario in the region.

## Example Usage

In the example below, specific Sourcedeploy scenario name.

```hcl
data "ncloud_sourcedeploy_project_stage_scenario" "deploy_scenario" {
  project_id = 1234
  stage_id   = 1234
  id         = 1234
}

output "output_scenario"{
  value = data.ncloud_sourcedeploy_project_stage_scenario.deploy_scenario
}
```


## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of Sourcedeploy project.
* `stage_id` - (Required) The ID of Sourcedeploy stage.
* `id` - (Required) The ID of Sourcedeploy scenario.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of scenario.
* `description` - Sourcedeploy project description.
    * `config` - scenario config.
        * `strategy` - Deployment strategy.
        * `file` - Deployment file.
            * `type` - File type.
            * `object_storage` - Objectstorage config.
                * `bucket` - The Name of ObjectStorage bucket.
                * `object` - ObjectStorage object .
            * `source_build` - Sourcebuild config.
                * `id` - The ID of SourceBiuld project.
                * `name` - The name of SourceBuild project.
        * `rollboack` - Rollback on deployment failure.
        * `deploy_command` - Commands to execute in deploy.
            * `pre_deploy` - Commands before deploy.
                * `user` - Running Account.
                * `command` - Run Command.
            * `path` - Deploy file.
                * `source_path` - Source file path.
                * `deploy_path` - Deploy Path.
            * `post_deploy` - Commands after deploy.
                * `user` - Running Account.
                * `command` - Run Command.
        * `load_balancer` - Loadbalancer target group for blue-green deployment.
            * `load_balancer_target_group_no` - Loadbalancer Target Group no.
            * `load_balancer_target_group_name` - The name of Loadbalancer Target Group.
            * `delete_server` - Whether to delete Servers in the auto scaling group.
        * `manifest` - Manifest file for Kubernetesservice deployment.
            * `type` - Repository type.
            * `repository_name` - The name of repository.
            * `branch` - The name of repository branch.
            * `path` - File path.
        * `canary_config` - config when deploying Kubernetesservice canary.
            * `canary_count` - Number of baseline and canary pod.
            * `analysis_type` - Canary analysis method.
            * `timeout` - Maximum time of deployment/cancellation.
            * `prometheus` - Prometheus Url.
            * `env` - Analysis environment.
                * `baseline` - Analysis environment variable > baseline.
                * `canary` - Analysis environment variable > canary.
            * `metrics` - Metric.
                * `name` - Metric name.
                * `success_criteria` - Success criteria.
                * `query_type` - Query type.
                * `weight` - Weight.
                * `metric` - Metric.
                * `filter` - Filter.
                * `query` - Query.
            * `analysis_config` - Analysis config.
                * `duration` - Analysis time.
                * `delay` - Analysis delay time.
                * `interval` - Analysis cycle.
                * `step` - Metric collection cycle.
            * `pass_score` - Analysis success score.
        * `path` - Deploy file.
            * `source_path` - Source file path.
            * `deploy_path` - Deploy Path.
