# Data Source: ncloud_sourcepipieline_project

~> **Note** This data source only supports 'public' site.

~> **Note:** This data source is a beta release. Some features may change in the future.

This module can be useful for getting detail of Sourcepipeline project created before.

## Example Usage

In the example below, Retrieves Sourcepipeline project detail with the project id is '1234'.

```hcl

data "ncloud_sourcepipeline_project" "pipeline_project" {
    id = 1234
}

output "lookup-pipeline_project-output" {
    value = data.ncloud_sourcepipeline_project.pipeline_project
}

```

## Argument Reference

The following arguments are supported:

*   `id` - (Required) Sourcepipeline project id.

## Attributes Reference

The following attributes are exported:

*   `name` - The project name.
*   `description` - Description of the project.
*   `task`
    *   `name` - Task name.
    *   `type` - Task type. Accepted values: `SourceBuild` | `SourceDeploy`
    *   `config`
        *   `project_id` - Project Id of a task.
        *   `stage_id` - Stage Id of a task.
        *   `scenario_id` - Scenario Id of a task.
        *   `target`
            *   `type` - Target type of a task. Accepted values: `SourceCommit` | `GitHub` | `BitBucket` | `SourceBuild` | `ObjectStorage` | `KubernetesService`.
            *   `repository_name` - Target source repository of the Sourcebuild task.
            *   `repository_branch` - Target repository branch of the Sourcebuild task.
            *   `project_name` - Target sourcebuild project name of the Sourcedeploy task.
            *   `file` - Target file name of the Sourcedeploy task.
            *   `manifest` - Target manifest file name of the Sourcedeploy task.
            *   `full_manifest` - List of target manifest files name of the Sourcedeploy task.
    *   `linked_tasks` - List of linked tasks.
*   `triggers`
    *   `repository` - Repository trigger.
        *   `type` - Type of the repository.
        *   `name` - Name of the repository.
        *   `branch` - Name of a repository branch.
    *   `schedule` - Schedule trigger.
        *   `day` - List of day of week .
        *   `time` - Time to trigger.
        *   `timezone` - Timezone for trigger.
        *   `execute_only_with_change` -  Schedule trigger option. Schedule trigger always execute in time, if option is false. Schedule trigger execute when Sourcepipeline project configuration or Sourcecommit repository has changed, if option is true.
    *   `sourcepipeline` - Sourcepipeline trigger.
        *   `id` - Id of the sourcepipeline project to trigger execution of pipeline.
        *   `name` - Name of the sourcepipeline project to trigger execution of pipeline.
