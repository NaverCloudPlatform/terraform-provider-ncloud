provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_sourcebuild_compute" "compute" {
}

data "ncloud_sourcebuild_os" "os" {
}

data "ncloud_sourcebuild_runtime" "runtime" {
  os_id = data.ncloud_sourcebuild_os.os.os[0].id
}

data "ncloud_sourcebuild_runtime_version" "runtime_version" {
  os_id      = data.ncloud_sourcebuild_os.os.os[0].id
  runtime_id = data.ncloud_sourcebuild_runtime.runtime.runtime[0].id
}

data "ncloud_sourcebuild_docker" "docker" {
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
      repository = ncloud_sourcecommit_repository.test-repository.name
      branch     = "master"
    }
  }
  env {
    compute {
      id = data.ncloud_sourcebuild_compute.compute.compute[0].id
    }
    platform {
      type = "SourceBuild"
      config {
        os {
          id = data.ncloud_sourcebuild_os.os.os[0].id
        }
        runtime {
          id = data.ncloud_sourcebuild_runtime.runtime.runtime[0].id
          version {
            id = data.ncloud_sourcebuild_runtime_version.runtime_version.runtime_version[0].id
          }
        }
      }
    }
    docker {
      use = true
      id  = data.ncloud_sourcebuild_docker.docker.docker[0].id
    }
  }
}
