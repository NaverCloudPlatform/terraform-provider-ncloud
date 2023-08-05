---
subcategory: "Developer Tools"
---


# Data Source: ncloud_sourcedeploy_project_stage

~> **Note** This data source only supports 'public' site.

-> **Note:** This data source is a beta release. Some features may change in the future.

This resource is useful for look up the list of Sourcedeploy stage detail in the region.

## Example Usage

In the example below, Retrieves Sourcedeploy stage detail with the project id is '1234'.

```hcl
data "ncloud_sourcedeploy_project_stage" "deploy_project_stage" {
  project_id    = 1234
  id            = 1234
}

output "output_project_stage"{
  value = data.ncloud_sourcedeploy_project_stage.deploy_project_stage
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of Sourcedeploy project.
* `id` - (Required) The ID of Sourcedeploy stage.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of stage.
* `target_type` - The type of deploy target.
* `config` - The configuration of deploy target.
    * `server` - server
        * `id` - The id of server.
        * `name` - The name of server.
    * `auto_scaling_group_no` - The ID of Auto Scaling Group.
    * `auto_scaling_group_name` - The name of Auto Scaling Group.
    * `cluster_uuid` - The uuid of Kubernetes Service Cluster.
    * `cluster_name` - The name of Kubernetes Service Cluster.
    * `bucket_name` - The name of ObjectStorage bucket.
