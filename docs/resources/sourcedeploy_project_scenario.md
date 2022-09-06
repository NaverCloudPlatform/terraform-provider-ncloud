# Resource : ncloud_sourcedeploy_project_stage_scenario

-> **Note:** This resource is a beta release. Some features may change in the future.

This resource is useful for look up the list of Sourcedeploy scenario in the region.

## Example Usage

In the example below, specific Sourcedeploy scenario name.

```hcl
resource "ncloud_sourcedeploy_project" "project" {
	name    								= "test-deploy-project"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	name    								= "test-deploy-stage"
	target_type    								= "Server"
	config {
		server_ids  							= [1234]
	}
}

resource "ncloud_sourcedeploy_project_stage_scenario" "server_normal" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.svr_stage.id
	name    								= "test-deploy-scenario"
	description   	 						= "test"
	config {
		strategy 							= "normal"
		file {
			type     						= "SourceBuild"
			source_build {
				id 							= 1234
			}
		}
		rollback 							= true
		deploy_command {
			pre_deploy {
				user  						= "root"
				command   						= "echo pre"
			}
			path {
				source_path 				= "/"
				deploy_path 				= "/test"
			}
			post_deploy {
				user  						= "root"
				command   						= "echo post"
			}
		}
	}
}


```

Create Sourcedeploy scenario by referring to data sources (retrieve sourcebuild_project).

```hcl

data "ncloud_sourcebuild_projects" "test-sourcebuild" {
}

resource "ncloud_sourcedeploy_project_stage_scenario" "server_normal" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.svr_stage.id
	name    								= "test-deploy-scenario"
	description   	 						= "test"
	config {
		strategy 							= "normal"
		file {
			type     						= "SourceBuild"
			source_build {
				id 							= data.ncloud_sourcebuild_projects.test-sourcebuild.projects[0].id
			}
		}
		rollback 							= true
		deploy_command {
			pre_deploy {
				user  						= "root"
				command   						= "echo pre"
			}
			path {
				source_path 				= "/"
				deploy_path 				= "/test"
			}
			post_deploy {
				user  						= "root"
				command   						= "echo post"
			}
		}
	}
}


```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of Sourcedeploy project.
* `stage_id` - (Required) The ID of Sourcedeploy stage.

* `name` - (Required) The name of scenario.
* `description` - (Optional) Sourcedeploy project description.
    * `config` - (Required) scenario config.
        * `strategy` - (Required) Deployment strategy.
        * `file` - (Optional, Required If stage type is set to `Server` or `AutoScalingGroup` or `ObjectStorage`) Deployment file.
            * `type` - (Required) File type.
            * `object_storage` - (Optional, Required if file.type is set to `ObjectStorage` ) Objectstorage config.
                * `bucket` - (Required) The Name of ObjectStorage bucket.
                * `object` - (Required) ObjectStorage object .
            * `source_build` - (Optional, Required if file.type is set to`SourceBuild` ) Sourcebuild config.
                * `id` - (Required) The ID of SourceBiuld project. [`ncloud_sourcebuild_project` data source](../data-sources/sourcebuild_project.md)
        * `rollboack` - (Optional,  Required If stage type is set to `Server` or `AutoScalingGroup` ) Rollback on deployment failure.
        * `deploy_command` - (Optional) Commands to execute in deploy.
            * `pre_deploy` - (Optional) Commands before deploy.
                * `user` - (Required) Running Account.
                * `command` - (Required) Run Command.
            * `path` - (Optional) Deploy file.
                * `source_path` - (Required) Source file path.
                * `deploy_path` - (Required) Deploy Path.
            * `post_deploy` - (Optional) Commands after deploy.
                * `user` - (Required) Running Account.
                * `command` - (Required) Run Command.
        * `load_balancer` - (Optional, Required If stage type is set to `AutoScalingGroup` & strategy is set to `blueGreen`) Loadbalancer target group for blue-green deployment. 
            * `load_balancer_target_group_no` - (Required) Loadbalancer Target Group no. [`ncloud_lb_target_group` data source](../data-sources/lb_target_group.md)
            * `delete_server` - (Required) Whether to delete Servers in the auto scaling group.
        * `manifest` - (Optional, Required If stage type is set to `KubernetesService`) Manifest file for Kubernetesservice deployment.
            * `type` - (Required) Repository type.
          	* `repository_name` - (Required) The name of repository.
            * `branch` - (Required) The name of repository branch.
              * `path` - (Required) File path.
        * `canary_config` - (Optional, Required If stage type is set to `KubernetesService` &  strategy is set to `canary` ) config when deploying Kubernetesservice canary.
			* `analysis_type` - (Required) Canary analysis method.
            * `canary_count` - (Required) Number of baseline and canary pod.
            * `timeout` - (Optional,  Required if canaryConfig.analysisType=`manual`) Maximum time of deployment/cancellation.
            * `prometheus` - (Optional, Required if canaryConfig.analysisType=`auto`) Prometheus Url.
            * `env` - (Optional,  Required if canaryConfig.analysisType=`auto`) Analysis environment.
                * `baseline` - (Required) Analysis environment variable > baseline.
                * `canary` - (Required) Analysis environment variable > canary.
            * `metrics` - (Optional, Required if canaryConfig.analysisType=`auto`) Metric.
                * `name` - (Required) Metric name.
                * `success_criteria` - (Required) Success criteria.
                * `weight` - (Required) Weight.
				* `query_type` - (Required) Query type.
                * `metric` - (Optional, Required if canaryConfig.query_type is set to `default`  ) Metric.
                * `filter` - (Optional,  Required if canaryConfig.query_type is set to`default` ) Filter.
                * `query` - (Optional,  Required if canaryConfig.query_type is set to `promQL` ) Query.
            * `analysis_config` - (Optional, Required if canaryConfig.analysisType is set to `auto` ) Analysis config.
                * `duration` - (Required) Analysis time.
                * `delay` - (Required) Analysis delay time.
                * `interval` - (Required) Analysis cycle.
                * `step` - (Required) Metric collection cycle.
            * `pass_score` - (Optional, Required if canaryConfig.analysisType=`auto`) Analysis success score.
        * `path` - (Optional, Required If stage type is set to `ObjectStorage`) Deploy file.
            * `source_path` - (Required) Source file path.
            * `deploy_path` - (Required) Deploy Path.


## Attributes Reference

* `id` - The ID of scenario.
* `name` - The name of scenario.