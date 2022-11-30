provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

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
  stage_id   = data.ncloud_sourcedeploy_project_stages.test-sourcedeploy_stages.stages[0].id
}

data "ncloud_sourcepipeline_projects" "test-sourcepipeline" {
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
      project_id  = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
      stage_id    = data.ncloud_sourcedeploy_project_stages.test-sourcedeploy_stages.stages[0].id
      scenario_id = data.ncloud_sourcedeploy_project_stage_scenarios.test-sourcedeploy_scenarios.scenarios[0].id
    }
    linked_tasks = ["task_name_1"]
  }
  triggers {
    repository {
      type   = "sourcecommit"
      name   = ncloud_sourcecommit_repository.test-sourcecommit.name
      branch = "master"
    }
    schedule {
      day                       = ["MON", "TUE"]
      time                      = "13:01"
      timezone                  = "Asia/Seoul (UTC+09:00)"
      execute_only_with_change = false
    }
    sourcepipeline {
      id = data.ncloud_sourcepipeline_projects.test-sourcepipeline.projects[0].id
    }
  }
}
