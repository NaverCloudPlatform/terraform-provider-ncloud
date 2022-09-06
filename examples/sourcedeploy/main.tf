provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}


resource "ncloud_sourcedeploy_project" "sd-project" {
	name = "test-deploy-project"
}


data "ncloud_server" "server" {
	filter {
		name    = "name"
		values  = ["terraform-test"]
	}
}


resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	name    							  = "test-deploy-stage"
	target_type    							  = "Server"
	config {
		server_ids  						= [data.ncloud_server.server.id]
	}
}

data "ncloud_sourcebuild_projects" "test-sourcebuild" {
}

resource "ncloud_sourcedeploy_project_stage_scenario" "server_normal" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.svr_stage.id
	name    								  = "test-deploy-scenario"
	description   	 					= "test"
	config {
		strategy 							  = "normal"
		file {
			type     						  = "SourceBuild"
			source_build {
				id 							    = data.ncloud_sourcebuild_projects.test-sourcebuild.projects[0].id
			}
		}
		rollback 							  = true
		deploy_command {
			pre_deploy {
				user  						  = "root"
				command   						  = "echo pre"
			}
			path {
				source_path 				= "/"
				deploy_path 				= "/test"
			}
			post_deploy {
				user  						  = "root"
				command   						  = "echo post"
			}
		}
	}
}
