---
subcategory: "Developer Tools"
---


# Data Source: ncloud_sourcebuild_project

~> **Note** This data source only supports 'public' site.

~> **Note:** This data source is a beta release. Some features may change in the future.

This data source is useful for look up Sourcebuild project detail in the region.

## Example Usage

In the example below, Retrieves Sourcebuild project detail with the project id is '1234'.

```hcl
data "ncloud_sourcebuild_project" "build_project" {
  id = 1234
}

output "lookup-build_project-output" {
  value = data.ncloud_sourcebuild_project.build_project
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) Sourcebuild Project ID.

## Attributes Reference

* `project_no` - Sourcebuild Project ID. (It is the same result as `id`)
* `name` - Name of the Sourcebuild Project.
* `description` - Sourcebuild project description.
* `source` - Build target's type and config.
    * `type` - Build target type.
    * `config` - Build target config.
        * `repository_name` - Repository name to build.
        * `branch` - Branch to build.
* `env` - Build environment.
    * `compute` - Computing environment to build.
        * `id` - Computing type id.
        * `cpu` - CPU of build environment.
        * `mem` - Memory of build environment.
    * `platform` - Information about the build environment image.
        * `type` - Build environment image type.
        * `config` - Build environment image config.
            * `os` - OS config.
                * `id` - OS id.
                * `name` - OS name.
                * `version` - OS version.
                * `archi` - OS architecture.
            * `runtime` - Runtime config.
                * `id` - runtime id.
                * `name` - runtime name.
                * `version` - runtime version.
                    * `id` - runtime version id.
                    * `name` - runtime version name.
            * `registry` - Registry name of NCP Container Registry where the image to build is located.
            * `image` - Container image name to build.
            * `tag` - Container image tag to build.
    * `docker_engine` - Docker engine to use when building docker image.
        * `use` - Whether or not to use of docker engine.
        * `id` - Docker engine id.
        * `name` - Docker engine name.
    * `timeout` - Build timeout (in Minutes).
    * `env_var` - Environment variable to use for build.
        * `key` - Key of environment variable.
        * `value` - Value of environment variable.
* `build_command` - Commands to execute in build.
    * `pre_build` - Commands before build.
    * `in_build` - Commands during build.
    * `post_build` - Commands after build.
    * `docker_image_build` - Docker image build config.
        * `use` - Whether or not to use of dockerbuild.
        * `dockerfile` - Dockerfile path in build source folder.
        * `registry` - Registry name of NCP Container Registry to store the image.
        * `image` - Image name to upload to registry.
        * `tag` - Image tag to upload to registry.
        * `latest` - Save status of the latest tag.
* `artifact` - Artifact to save build results.
    * `use` - Whether or not to save build results.
    * `path` - Location to save build results.
    * `object_storage_to_upload` - Object Storage to save build results.
        * `bucket` - Bucket name of NCP Object Storage to save build results.
        * `path` - path in the NCP Object Storage bucket to save build results.
        * `filename` - File name to save build results.
    * `backup` - Whether or not to backup build results.
* `build_image_upload` - Save build environment after completing this build.
    * `use` - Whether or not to save build environment.
    * `container_registry_name` - Registry name of NCP Container Registry to store the image of the build environment after completing the build.
    * `image_name` - Image name to upload to registry.
    * `tag` - Image tag to upload to registry.
    * `latest` -  Save status of the latest tag.
* `linked` - Set up linkage with other services related this build.
    * `cloud_log_analytics` - Whether or not to save build log in the NCP Cloud Log Analytics.
    * `cloud_log_analytics` - Whether or not to check safety using NCP File Safer.
* `lastbuild` - Information of last build.
    * `id` - ID of last build.
    * `status` - Status of last build.
    * `timestamp` - Time of last build.
* `created` - Information about creating a Sourcebuild project.
    * `user` - Created user
    * `timestamp` - Created time
