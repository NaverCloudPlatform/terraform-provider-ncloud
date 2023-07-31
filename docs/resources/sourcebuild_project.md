---
subcategory: "Developer Tools"
---


# Resource: ncloud_sourcebuild_project

~> **Note** This resource only supports 'public' site.

~> **Note:** This resource is a beta release. Some features may change in the future.

Provides a Sourcebuild project resource.

## Example Usage

```hcl
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
      branch     = "master"
    }
  }
  env {
    compute {
      id = 1
    }
    platform {
      type = "SourceBuild"
      config {
        os {
          id = 1
        }
        runtime {
          id = 1
          version {
            id = 1
          }
        }
      }
    }
    timeout = 200
    docker_engine {
      use = true
      id  = 1
    }
    env_var {
      key   = "KEY"
      value = "VALUE"
    }
    env_var {
      key   = "KEY2"
      value = "VALUE2"
    }
  }
  build_command {
    pre_build   = ["pwd"]
    in_build = ["make"]
    docker_image_build {
      use        = true
      registry   = "test-registry"
      dockerfile = "/Dockerfile"
      image      = "custom-build-image"
      tag        = "1.0"
    }
  }
  linked {
    cloud_log_analytics = false
    file_safer          = true
  }
}
```

Create Sourcebuild project by referring to data sources (retrieve compute, os, runtime, runtime version and docker engine).

```hcl
data "ncloud_sourcebuild_project_computes" "computes" {
}

data "ncloud_sourcebuild_project_os" "os" {
  filter {
    name   = "name"
    values = ["ubuntu"]
  }
}

data "ncloud_sourcebuild_project_os_runtimes" "runtimes" {
  os_id = data.ncloud_sourcebuild_project_os.os.os[0].id
}

data "ncloud_sourcebuild_project_os_runtime_versions" "runtime_versions" {
  os_id      = data.ncloud_sourcebuild_project_os.os.os[0].id
  runtime_id = data.ncloud_sourcebuild_project_os_runtimes.runtimes.runtimes[0].id
}

data "ncloud_sourcebuild_project_docker_engines" "docker_engines" {
  filter {
    name   = "name"
    values = ["Docker:18.09.1"]
  }
}

resource "ncloud_sourcebuild_project" "test-build-project" {
  name        = "test-build-project"
  description = "my build project"
  source {
    type = "SourceCommit"
    config {
      repository_name = "test-repository"
      branch     = "master"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Sourcebuild Project. Specify a name that is only English letters, numbers, and special characters (-, _).
* `description` - (Optional) Sourcebuild project description.
* `source` - (Required) Build target's type and config.
    * `type` - (Required) Build target type. Accepted values: `SourceCommit`. (Other repository types are not supported yet.)
        * [`ncloud_sourcecommit_repository` data source](../data-sources/sourcecommit_repository.md)
    * `config` - (Required) Build target config.
        * `repository_name` - (Required) Repository name to build.
        * `branch` - (Required) Branch to build.
* `env` - (Required) Build environment.
    * `compute` - (Required) Computing environment to build.
        * [`ncloud_sourcebuild_project_compute` data source](../data-sources/sourcebuild_project_computes.md)
        * `id` - (Required) Computing type id.
    * `platform` - (Required) Information about the build environment image.
        * `type` - (Required) Build environment image type. Accepted values: `SourceBuild`, `ContainerRegistry` `PublicRegistry`.
        * `config` - (Required) Build environment image config.
            * `os` - (Optional, Required if `env.platform.type` is set to `SourceBuild`) OS config.
                * [`ncloud_sourcebuild_project_os` data source](../data-sources/sourcebuild_project_os.md)
                * `id` - (Required) OS id.
            * `runtime` - (Optional, Required if `env.platform.type` is set to `SourceBuild`) Runtime config.
                * [`ncloud_sourcebuild_project_os_runtimes` data source](../data-sources/sourcebuild_project_os_runtimes.md)
                * `id` - (Required) runtime id.
                * `version` - (Required) runtime version.
                    * [`ncloud_sourcebuild_project_os_runtime_versions` data source](../data-sources/sourcebuild_project_os_runtime_versions.md)
                    * `id` - (Required) runtime version id.
            * `registry` - (Optional, Required if `env.platform.type` is set to `ContainerRegistry`) Registry name of NCP Container Registry where the image to build is located.
            * `image` - (Optional, Required if `env.platform.type` is set to `ContainerRegistry` or `PublicRegistry`) Container image name to build.
            * `tag` - (Optional, Required if `env.platform.type` is set to `ContainerRegistry` or `PublicRegistry`) Container image tag to build.
    * `docker_engine` - (Optional) Docker engine to use when building docker image.
        * [`ncloud_sourcebuild_project_docker` data source](../data-sources/sourcebuild_project_docker_engines.md)
        * `use` - (Required) Whether or not to use of docker engine. (Default `false`)
        * `id` - (Optional) Docker engine id.
    * `timeout` - (Optional) Build timeout (in Minutes). Specify it between `5` and `540`. Default `60`.
    * `env_var` - (Optional) Environment variables to use for build.
        * `key` - (Required) Key of environment variable.
        * `value` - (Required) Value of environment variable.
* `build_command` - (Optional) Commands to execute in build.
    * `pre_build` - (Optional) Commands before build.
    * `in_build` - (Optional) Commands during build.
    * `post_build` - (Optional) Commands after build.
    * `docker_image_build` - (Optional) Docker image build config.
        * `use` - (Optional) Whether or not to use of dockerbuild. (Default `false`)
        * `dockerfile` - (Optional, Required if `build_command.docker_image_build.use` is set to `true`) Dockerfile path in build source folder.
        * `registry` - (Optional, Required if `build_command.docker_image_build.use` is set to `true`) Registry name of NCP Container Registry to store the image.
        * `image` - (Optional, Required if `build_command.docker_image_build.use` is set to `true`) Image name to upload to registry.
        * `tag` - (Optional, Required if `build_command.docker_image_build.use` is set to `true`) Image tag to upload to registry.
        * `latest` - (Optional) Save status of the latest tag. (Default `false`)
* `artifact` - (Optional) Artifact to save build results.
    * `use` - (Optional) Whether or not to save build results. (Default `false`)
    * `path` - (Optional, Required if `artifact.use` is set to `true`) Location to save build results.
    * `object_storage_to_upload` - (Optional, Required if `artifact.use` is set to `true`) Object Storage to save build results.
        * `bucket` - (Required) Bucket name of NCP Object Storage to save build results.
        * `path` - (Required) path in the NCP Object Storage bucket to save build results.
        * `filename` - (Required) File name to save build results.
    * `backup` - (Optional) Whether or not to backup build results.
* `build_image_upload` - (Optional) Save build environment after completing this build.
    * `use` - (Optional) Whether or not to save build environment. (Default `false`)
    * `container_registry_name` - (Optional, Required if `build_image_upload.use` is set to `true`) Registry name of NCP Container Registry to store the image of the build environment after completing the build.
    * `image_name` - (Optional, Required if `build_image_upload.use` is set to `true`) Image name to upload to registry.
    * `tag` - (Optional, Required if `build_image_upload.use` is set to `true`) Image tag to upload to registry.
    * `latest` - (Optional)  Save status of the latest tag. (Default `false`)
* `linked` - (Optional) Set up linkage with other services related this build.
    * `cloud_log_analytics` - (Optional) Whether or not to save build log in the NCP Cloud Log Analytics. (Default `false`)
    * `file_safer` - (Optional) Whether or not to check safety using NCP File Safer. (Default `false`)

## Attributes Reference

* `id` - Sourcebuild Project ID.
* `project_no` - Sourcebuild Project ID.
* `env`
    * `compute`
        * `cpu` - CPU of build environment.
        * `mem` - Memory of build environment.
    * `platform`
        * `config`
            * `os`
                * `name` - OS name of build environment.
                * `version` - OS version of build environment.
                * `archi` - OS architecture of build environment.
            * `runtime`
                * `name` - Runtime name of build environment.
                * `version`
                    * `name` - Runtime version name of build environment.
    * `docker_engine`
        * `name` - Docker engine name.
* `lastbuild` - Information of last build.
    * `id` - ID of last build.
    * `status` - Status of last build.
    * `timestamp` - Timestamp of last build.
* `created` - Information about creating a Sourcebuild project.
    * `user` - Created user
    * `timestamp` - Created timestamp
