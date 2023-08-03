---
subcategory: "Developer Tools"
---


# Resource : ncloud_sourcedeploy_project_stage

~> **Note** This resource only supports 'public' site.

-> **Note:** This resource is a beta release. Some features may change in the future.

Provides a Sourcedeploy stage resource.

## Example Usage

```hcl
resource "ncloud_sourcedeploy_project" "project" {
  name = "test-deploy-project"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
  project_id   = ncloud_sourcedeploy_project.project.id
  name         = "test-deploy-stage"
  target_type  = "Server"
  config {
    server {
      id = 1234
    } 
  }
}

```

Create Sourcedeploy stage by referring to data sources (retrieve server).

```hcl
data "ncloud_server" "server" {
  filter {
    name   = "name"
    values = ["test-server"]
  }
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
  project_id  = ncloud_sourcedeploy_project.project.id
  name        = "test-deploy-stage"
  target_type = "Server"
  config {
    server {
      id = data.ncloud_server.server.id
    } 
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of Sourcedeploy project.
* `name` - (Required) The name of stage.
* `target_type` - (Required) The type of deploy target. Accepted values: `Server`, `AutoScalingGroup`, `KubernetesService`, `ObjectStorage`.
* `config` - (Required) The configuration of deploy target.
    * `server` - server 
        * `id` - (Optional, Required If type=`Server`) The no of server. [`ncloud_server` data source](../data-sources/server.md)
    * `auto_scaling_group_no` - (Optional, Required If type=`AutoScalingGroup`) The ID of Auto Scaling Group.  [`ncloud_auto_scaling_group` data source](../data-sources/auto_scaling_group.md)
    * `cluster_uuid` - (Optional, Required If type=`KubernetesService`) The uuid of Kubernetes Service Cluster.  [`ncloud_nks_cluster` data source](../data-sources/nks_cluster.md)
    * `bucket_name` - (Optional, Required If type=`ObjectStorage`) The name of ObjectStorage bucket.


## Attributes Reference

* `id` - The ID of stage.
* `config` - (Required) The configuration of deploy target.
    * `server` - server 
        * `name` - The name of server.
    * `auto_scaling_group_name` - The name of Auto Scaling Group.
    * `cluster_name` - The name of Kubernetes Service Cluster.

## Import

SourceDeploy stage can be imported using the project_id and stage_id separated by a colon (:), e.g.,

$ terraform import ncloud_sourcedeploy_project_stage.my_stage project_id:stage_id