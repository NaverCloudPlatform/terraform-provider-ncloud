provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_sourcebuild_project_computes" "computes" {
}

data "ncloud_sourcebuild_project_os" "os" {
}

data "ncloud_sourcebuild_project_os_runtimes" "runtimes" {
  os_id = data.ncloud_sourcebuild_project_os.os.os[0].id
}

data "ncloud_sourcebuild_project_os_runtime_versions" "runtime_versions" {
  os_id      = data.ncloud_sourcebuild_project_os.os.os[0].id
  runtime_id = data.ncloud_sourcebuild_project_os_runtimes.runtimes.runtimes[0].id
}

data "ncloud_sourcebuild_project_docker_engines" "docker_engines" {
}

resource "ncloud_sourcecommit_repository" "test-repository" {
  name = "test-repository"
}

resource "ncloud_sourcebuild_project" "test-build-project" {
  name        = "test-build-project"
  description = "my build project"
  source {
    type = "SourceCommit"
    config {
      repository_name = ncloud_sourcecommit_repository.test-repository.name
      branch          = "master"
    }
  }
  env {
    compute {
      id = data.ncloud_sourcebuild_project_computes.computes.computes[0].id
    }
    platform {
      type = "SourceBuild"
      config {
        os {
          id = data.ncloud_sourcebuild_project_os.os.os[0].id
        }
        runtime {
          id = data.ncloud_sourcebuild_project_os_runtimes.runtimes.runtimes[0].id
          version {
            id = data.ncloud_sourcebuild_project_os_runtime_versions.runtime_versions.runtime_versions[0].id
          }
        }
      }
    }
    docker_engine {
      use = true
      id  = data.ncloud_sourcebuild_project_docker_engines.docker_engines.docker_engines[0].id
    }
  }
}
