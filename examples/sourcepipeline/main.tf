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

data "ncloud_sourcedeploy_stages" "test-sourcedeploy_stages" {
  project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
}

data "ncloud_sourcedeploy_scenarioes" "test-sourcedeploy_scenarioes" {
  project_id = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
  stage_id   = data.ncloud_sourcedeploy_stages.test-sourcedeploy_stages.stages[0].id
}

resource "ncloud_sourcepipeline_project" "test-sourcepipeline" {
  name        = "tf-sourcepipeline_project-test"
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
      project_id  = data.ncloud_sourcedeploy_projects.test-sourcedeploy_projects.projects[0].id
      stage_id    = data.ncloud_sourcedeploy_stages.test-sourcedeploy_stages.stages[0].id
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
