# Resource : ncloud_sourcedeploy_project_stage

Provides a Sourcedeploy stage resource.

## Example Usage

```hcl
resource "ncloud_sourcedeploy_project" "project" {
	name    							= "test-deploy-project"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	name    							  = "test-deploy-stage"
	type    							  = "Server"
	config {
		server_no  						= [1234]
	}
}

```

```hcl
Create Sourcedeploy stage by referring to data sources (retrieve server).

data "ncloud_server" "server" {
	filter {
		name = "name"
		values = ["test-server"]
	}
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	name    							  = "test-deploy-stage"
	type    							  = "Server"
	config {
		server_no  						= [data.ncloud_server.server.id]
	}
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of Sourcedeploy project.
* `name` - (Required) The name of stage.
* `type` - (Required) The type of deploy target.
* `config` - (Required) The configuration of deploy target.
    * `server_no` - (Optional, Required If type=`Server`) The no of server. [`ncloud_server` data source](../data-source/server.md)
    * `auto_scaling_group_no` - (Optional, Required If type=`AutoScalingGroup`) The ID of Auto Scaling Group.  [`ncloud_auto_scaling_group` data source](../data-source/auto_scaling_group.md)
    * `cluster_uuid` - (Optional, Required If type=`KubernetesService`) The uuid of Kubernetes Service Cluster.  [`ncloud_nks_cluster` data source](../data-source/nks_cluster.md)
    * `bucket_name` - (Optional, Required If type=`ObjectStorage`) The name of ObjectStorage bucket.


## Attributes Reference

* `id` - The ID of stage.
* `name` - The name of stage.
