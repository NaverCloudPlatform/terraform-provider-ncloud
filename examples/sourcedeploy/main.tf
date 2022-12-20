provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_server" "server" {
  filter {
    name   = "name"
    values = ["server_name"]
  }
}

data "ncloud_auto_scaling_group" "asg" {
  filter {
    name   = "name"
    values = ["asg_name"]
  }
}

data "ncloud_sourcecommit_repositories" "test-repo" {
  filter {
    name   = "name"
    values = ["repo_name"]
    regex  = true
  }
}

data "ncloud_sourcebuild_projects" "test-sourcebuild" {
  filter {
    name   = "name"
    values = ["build_name"]
  }
}

resource "ncloud_sourcedeploy_project" "sd-project" {
  name = "test-deploy-project"
}

resource "ncloud_sourcedeploy_project_stage" "test-stage-svr" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  name        = "svr"
  target_type = "Server"
  config {
    server{
      id = data.ncloud_server.server.id
    } 
  }
}
resource "ncloud_sourcedeploy_project_stage" "test-stage-asg" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  name        = "asg"
  target_type = "AutoScalingGroup"
  config {
    auto_scaling_group_no = data.ncloud_auto_scaling_group.asg.id
  }
}
resource "ncloud_sourcedeploy_project_stage" "test-stage-nks" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  name        = "nks"
  target_type = "KubernetesService"
  config {
    cluster_uuid = "cluster_uuid"
  }
}
resource "ncloud_sourcedeploy_project_stage" "test-stage-obj" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  name        = "obj"
  target_type = "ObjectStorage"
  config {
    bucket_name = "bucket_name"
  }
}


resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-server-normal" {
  project_id  = ncloud_sourcedeploy_project.project.id
  stage_id    = ncloud_sourcedeploy_project_stage.svr_stage.id
  name        = "server_normal"
  description = "test"
  config {
    strategy = "normal"
    file {
      type = "SourceBuild"
      source_build {
        id = data.ncloud_sourcebuild_projects.test-sourcebuild.projects[0].id
      }
    }
    rollback = true
    deploy_command {
      pre_deploy {
        user    = "root"
        command = "echo pre"
      }
      path {
        source_path = "/"
        deploy_path = "/test"
      }
      post_deploy {
        user    = "root"
        command = "echo post"
      }
    }
  }
}



resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-asg-normal" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  stage_id    = ncloud_sourcedeploy_project_stage.test-stage-asg.id
  name        = "asg_normal"
  description = "test"
  config {
    strategy = "normal"
    file {
      type = "SourceBuild"
      source_build {
        id = ncloud_sourcebuild_project.test-build-project.id
      }
    }
    rollback = true
    deploy_command {
      pre_deploy {
        user    = "root"
        command = "echo pre"
      }
      path {
        source_path = "/"
        deploy_path = "/test"
      }
      post_deploy {
        user    = "root"
        command = "echo post"
      }
    }
  }
}

resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-asg-bg" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  stage_id    = ncloud_sourcedeploy_project_stage.test-stage-asg.id
  name        = "asg_bg"
  description = "test"
  config {
    strategy = "blueGreen"
    file {
      type = "SourceBuild"
      source_build {
        id = ncloud_sourcebuild_project.test-build-project.id
      }
    }
    rollback = true
    deploy_command {
      pre_deploy {
        user    = "root"
        command = "echo pre"
      }
      path {
        source_path = "/"
        deploy_path = "/test"
      }
      post_deploy {
        user    = "root"
        command = "echo post"
      }
    }
    load_balancer {
      load_balancer_target_group_no = "lb_target_group_no"
      delete_server                 = true
    }
  }
}

resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-nks-rolling" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  stage_id    = ncloud_sourcedeploy_project_stage.test-stage-nks.id
  name        = "nks_rolling"
  description = "test"
  config {
    strategy = "rolling"
    manifest {
      type            = "SourceCommit"
      repository_name = ncloud_sourcecommit_repository.test-repo.name
      branch          = "master"
      path            = ["/deployment/prod.yaml"]
    }
  }
}

resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-nks-bg" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  stage_id    = ncloud_sourcedeploy_project_stage.test-stage-nks.id
  name        = "nks_bg"
  description = "test"
  config {
    strategy = "blueGreen"
    manifest {
      type            = "SourceCommit"
      repository_name = ncloud_sourcecommit_repository.test-repo.name
      branch          = "master"
      path            = ["/deployment/canary.yaml"]
    }
  }
}

resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-nks-canary-manual" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  stage_id    = ncloud_sourcedeploy_project_stage.test-stage-nks.id
  name        = "nks_canary_manual"
  description = "test"
  config {
    strategy = "canary"
    manifest {
      type            = "SourceCommit"
      repository_name = ncloud_sourcecommit_repository.test-repo.name
      branch          = "master"
      path            = ["/deployment/canary.yaml"]
    }
    canary_config {
      analysis_type = "manual"
      timeout       = 10
      canary_count  = 1
    }
  }
}

resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-nks-canary-auto" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  stage_id    = ncloud_sourcedeploy_project_stage.test-stage-nks.id
  name        = "nks_canary_auto"
  description = "test"
  config {
    strategy = "canary"
    manifest {
      type            = "SourceCommit"
      repository_name = ncloud_sourcecommit_repository.test-repo.name
      branch          = "master"
      path            = ["test.yaml"]
    }
    canary_config {
      analysis_type = "auto"
      canary_count  = 1
      prometheus    = "prometheus_url"
      env {
        baseline = "baselineenv"
        canary   = "canaryenv"
      }
      metrics {
        name             = "success_rate"
        success_criteria = "base"
        query_type       = "promQL"
        weight           = 100
        query            = "test"
      }
      analysis_config {
        duration = 10
        delay    = 1
        interval = 1
        step     = 10
      }
      pass_score = 90
    }
  }
}

resource "ncloud_sourcedeploy_project_stage_scenario" "test-scenario-obj-normal" {
  project_id  = ncloud_sourcedeploy_project.test-project.id
  stage_id    = ncloud_sourcedeploy_project_stage.test-stage-obj.id
  name        = "obj_normal"
  description = "test"
  config {
    file {
      type = "SourceBuild"
      source_build {
        id = ncloud_sourcebuild_project.test-build-project.id
      }
    }
    path {
      source_path = "/"
      deploy_path = "/terraform"
    }
  }
}