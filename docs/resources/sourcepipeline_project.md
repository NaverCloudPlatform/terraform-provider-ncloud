# Resource: ncloud_sourcepipeline_project

Provides a Sourcepipeline project resource.

## Example Usage

```hcl
resource "ncloud_sourcecommit_repository" "test-sourcecommit" {
	name = "sourceCommit"
}

resource "ncloud_sourcepipeline_project" "test-sourcepipeline" {
    name = "tf-sourcepipeline_project-test"
    description = ""
    tasks {
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
    tasks {
        name = "task_name_2"
        type = "SourceDeploy"
        config {
            project_id = 1234
            stage_id = 5678
            scenario_id = 9101
        }
        linked_tasks   = ["task_name_1"]
    }
    trigger {
        setting = true
        sourcecommit {
            repository = ncloud_sourcecommit_repository.test-sourcecommit.name
        }
    }
}
```

Create Sourcepipeline project by referring to data sources (retrieve sourcebuild_projects, sourcedeploy_projects, sourcedeploy_stages, sourcedeploy_sceanrioes)
```hcl
resource "ncloud_sourcecommit_repository" "test-sourcecommit" {
    name = "sourceCommit"
}

data "ncloud_sourcebuild_projects" "test-sourcebuild" {
}

data "ncloud_sourcedeploy_projects" "test-sourcedeploy_projects" {
}

data "ncloud_sourcedeploy_stages" "test-sourcedeploy_stages" {
    project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
}

data "ncloud_sourcedeploy_scenarioes" "test-sourcedeploy_scenarioes" {
    project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
    stage_id = data.ncloud_sourcedeploy_stages.test-sourcedeploy_stages.stages[0].id
}

resource "ncloud_sourcepipeline_project" "test-sourcepipeline" {
    name = "tf-sourcepipeline_project-test"
    description = ""
    tasks {
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
    tasks {
        name = "task_name_2"
        type = "SourceDeploy"
        config {
            project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
            stage_id = data.ncloud_sourcedeploy_stages.test-sourcedeploy_stages.stages[0].id
            scenario_id = data.ncloud_sourcedeploy_scenarioes.test-sourcedeploy_scenarioes.scenarioes[0].id
        }
        linked_tasks = ["task_name_1"]
    }
    trigger {
        setting = true
        sourcecommit {
            repository = ncloud_sourcecommit_repository.test-sourcecommit.name
        }
    }
}
```

## Argument Reference

The following arguments are supported:

*   `name` - (Required) The project name to create.
*   `description` - (Optional) The project description to create.
*   `tasks` - (Required) `tasks` block describes task information.
    *   `name` - (Required) Task name.
    *   `type` - (Required) Task type. Select between SourceBuild and SourceDeploy. Accepted values: `SourceBuild` | `SourceDeploy` (`SourceDeploy` is available only in VPC environment).
    *   `config` - (Required) `config` block describes task configuration.
        *   `project_id` - (Required) Project Id of a task. Get avaliable values using the datasource `ncloud_sourcebuild_projects` or `ncloud_sourcedeploy_projects`
        *   `stage_id` - (Optional, Required if `tasks.type` value is SourceDeploy) Stage Id of a task. Get avaliable values using the datasource `ncloud_sourcedeploy_stages`
        *   `scenario_id` - (Optional, Required if `tasks.type` value is SourceDeploy) Scenario Id of a task. Get avaliable values using the datasource `ncloud_sourcedeploy_scenarioes`
        *   `target`- (Optional) Target of a task job.
            *   `repository_branch` - (Optional) Target repository branch of SourceBuild task. Default : main branch of target repository
*   `trigger` - (Required) `trigger` block describes trigger configuration.
    *   `setting` - (Required) Trigger setting option. You can decide whether to set trigger or not.
    *   `sourcecommit` - (Optional)
        *   `repository` - (Optional, Required if `trigger.setting` value is true) Name of sourcecommit repository to trigger execution of pipeline
        *   `branch` - (Optional, Required if `trigger.setting` value is true) Name of a repository branch to trigger execution of pipeline.

## Attributes Reference

*   `id` - The ID of Sourcepipeline project.
*   `tasks`
    *   `config`
        *   `target`
            *   `type` - Target type of a task. Accepted values: `SourceCommit` | `GitHub` | `BitBucket` | `SourceBuild` | `ObjectStorage` | `KubernetesService`
            *   `repository` - Target source repository of the Sourcebuild task. It is set only when `tasks.type` is SourceBuild
            *   `project_name` - Target Sourcebuild project name of the Sourcedeploy task. It is set only when `tasks.type` is SourceDeploy and `tasks.config.target.type` is SourceBuild.
            *   `file` - Target file name of the Sourcedeploy task. It is set only when `tasks.type` is SourceDeploy and `tasks.config.target.type` is ObjectStorage.
            *   `manifest` - Target manifest file name of the Sourcedeploy task. It is set only when `tasks.type` is SourceDeploy and `tasks.config.target.type` is KubernetesService.
            *   `full_manifest` - List of target manifest files name. It is set only when `tasks.type` is SourceDeploy and `tasks.config.target.type` is KubernetesService.
