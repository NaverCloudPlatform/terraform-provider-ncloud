# Resource: ncloud_sourcepipeline_project

~> **Note** This resource only supports 'public' site.

~> **Note:** This resource is a beta release. Some features may change in the future.

Provides a Sourcepipeline project resource.

## Example Usage

```hcl
resource "ncloud_sourcecommit_repository" "test-sourcecommit" {
	name = "sourceCommit"
}

resource "ncloud_sourcepipeline_project" "test-sourcepipeline" {
    name = "tf-sourcepipeline_project-test"
    task {
        name = "task_name_1"
        type = "SourceBuild"
        config {
	    project_id   = 1234
            target {
                repository_branch = "master"
            }
        }
        linked_tasks   = []
    }
    task {
        name = "task_name_2"
        type = "SourceDeploy"
        config {
            project_id = 1234
            stage_id = 5678
            scenario_id = 9101
        }
        linked_tasks   = ["task_name_1"]
    }
    triggers {
        sourcecommit {
            repository_name = ncloud_sourcecommit_repository.test-sourcecommit.name
            branch = "master"
        }
    }
}
```

Create Sourcepipeline project by referring to data sources (retrieve sourcebuild_projects, sourcedeploy_projects, sourcedeploy_project_stages, sourcedeploy_project_sceanrios)

```hcl
resource "ncloud_sourcecommit_repository" "test-sourcecommit" {
    name = "sourceCommit"
}

data "ncloud_sourcebuild_projects" "test-sourcebuild" {
}

data "ncloud_sourcedeploy_projects" "test-sourcedeploy_projects" {
}

data "ncloud_sourcedeploy_project_stages" "test-sourcedeploy_stages" {
    project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
}

data "ncloud_sourcedeploy_project_stage_scenarios" "test-sourcedeploy_scenarios" {
    project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
    stage_id = data.ncloud_sourcedeploy_project_stages.test-sourcedeploy_stages.stages[0].id
}

resource "ncloud_sourcepipeline_project" "test-sourcepipeline" {
    name = "tf-sourcepipeline_project-test"
    task {
        name = "task_name_1"
        type = "SourceBuild"
        config {
            project_id = data.ncloud_sourcebuild_projects.test-sourcebuild.projects[0].id
            target {
                repository_branch = "master"
            }
        }
        linked_tasks = []
    }
    task {
        name = "task_name_2"
        type = "SourceDeploy"
        config {
            project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
            stage_id = data.ncloud_sourcedeploy_project_stages.test-sourcedeploy_stages.stages[0].id
            scenario_id = data.ncloud_sourcedeploy_project_stage_scenarios.test-sourcedeploy_scenarios.scenarios[0].id
        }
        linked_tasks = ["task_name_1"]
    }
    triggers {
        sourcecommit {
            repository_name = ncloud_sourcecommit_repository.test-sourcecommit.name
            branch = "master"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

*   `name` - (Required) The project name to create.
*   `description` - (Optional) The project description to create.
*   `task` - (Required) `task` block describes task information.
    *   `name` - (Required) Task name.
    *   `type` - (Required) Task type. Select between SourceBuild and SourceDeploy. Accepted values: `SourceBuild` | `SourceDeploy` (`SourceDeploy` is available only in VPC environment).
    *   `config` - (Required) `config` block describes task configuration.
        *   `project_id` - (Required) Project Id of a task. Get avaliable values using the datasource `ncloud_sourcebuild_projects` or `ncloud_sourcedeploy_projects`
        *   `stage_id` - (Optional, Required if `task.type` value is SourceDeploy) Stage Id of a task. Get avaliable values using the datasource `ncloud_sourcedeploy_project_stages`
        *   `scenario_id` - (Optional, Required if `task.type` value is SourceDeploy) Scenario Id of a task. Get avaliable values using the datasource `ncloud_sourcedeploy_project_stage_scenarios`
        *   `target`- (Optional) Target of a task job.
            *   `repository_branch` - (Optional) Target repository branch of SourceBuild task. Default : main branch of target repository
    *   `linked_tasks` - (Required) Linked tasks which has to be executed before.
*   `triggers` - (Required) `triggers` block describes trigger configuration.
    *   `sourcecommit` - (Optional)
        *   `repository_name` - (Required) Name of sourcecommit repository to trigger execution of pipeline
        *   `branch` - (Required) Name of a repository branch to trigger execution of pipeline.

## Attributes Reference

*   `id` - The ID of Sourcepipeline project.
*   `task`
    *   `config`
        *   `target`
            *   `type` - Target type of a task. Accepted values: `SourceCommit` | `GitHub` | `BitBucket` | `SourceBuild` | `ObjectStorage` | `KubernetesService`
            *   `repository_name` - Target source repository of the Sourcebuild task. It is set only when `task.type` is SourceBuild
            *   `project_name` - Target Sourcebuild project name of the Sourcedeploy task. It is set only when `task.type` is SourceDeploy and `task.config.target.type` is SourceBuild.
            *   `file` - Target file name of the Sourcedeploy task. It is set only when `task.type` is SourceDeploy and `task.config.target.type` is ObjectStorage.
            *   `manifest` - Target manifest file name of the Sourcedeploy task. It is set only when `task.type` is SourceDeploy and `task.config.target.type` is KubernetesService.
            *   `full_manifest` - List of target manifest files name. It is set only when `task.type` is SourceDeploy and `task.config.target.type` is KubernetesService.
