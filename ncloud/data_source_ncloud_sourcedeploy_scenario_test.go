package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceDeployScenario(t *testing.T) {
	stageNameSvr := getTestSourceDeployScenarioName() + "svr"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceDeployScenarioConfig(stageNameSvr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcedeploy_scenario.scenario"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployScenarioConfig(scenarioNameSvr string) string {
	return fmt.Sprintf(`
data "ncloud_server" "server" {
	filter {
		name = "name"
		values = ["%[1]s"]
	}
}
data "ncloud_sourcebuild_project_compute" "compute" {
}

data "ncloud_sourcebuild_project_os" "os" {
}

data "ncloud_sourcebuild_project_runtime" "runtime" {
	os_id 					= data.ncloud_sourcebuild_project_os.os.os[0].id
}

data "ncloud_sourcebuild_project_runtime_version" "runtime_version" {
	os_id      				= data.ncloud_sourcebuild_project_os.os.os[0].id
	runtime_id 				= data.ncloud_sourcebuild_project_runtime.runtime.runtime[0].id
}

data "ncloud_sourcebuild_project_docker" "docker" {
}

resource "ncloud_sourcecommit_repository" "test-repo" {
	name 					= "tf-test-repository"
}

resource "ncloud_sourcebuild_project" "test-build-project" {
	name        					= "tf-test-project"
	description 					= "my build project"
	source {
		type 						= "SourceCommit"
		config {
			repository 				= ncloud_sourcecommit_repository.test-repo.name
			branch     				= "master"
		}
	}
	env {
		compute {
			id 						= data.ncloud_sourcebuild_project_compute.compute.compute[0].id
		}
		platform {
			type 					= "SourceBuild"
			config {
				os {
					id 				= data.ncloud_sourcebuild_project_os.os.os[0].id
				}
				runtime {
					id 				= data.ncloud_sourcebuild_project_runtime.runtime.runtime[0].id
					version {
						id 			= data.ncloud_sourcebuild_project_runtime_version.runtime_version.runtime_version[0].id
					}
				}
			}
		}
		docker {
			use 					= true
			id 						= data.ncloud_sourcebuild_project_docker.docker.docker[0].id
		}
		timeout 					= 500
		env_vars {
			key   					= "k1"
			value 					= "v1"
		}
	}
	cmd {
		pre  						= ["pwd", "ls"]
		build 						= ["pwd", "ls"]
		post						= ["pwd", "ls"]
	}
}

resource "ncloud_sourcedeploy_project" "project" {
	name    								= "tf-test-project"
}

resource "ncloud_sourcedeploy_stage" "svr_stage" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	name    								= "svr"
	type    								= "Server"
	config {
		server_no  							= [data.ncloud_server.server.id]
	}
}

resource "ncloud_sourcedeploy_scenario" "server_normal" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_stage.svr_stage.id
	name    								= "%[2]s"
	description   	 						= "test"
	config {
		strategy 							= "normal"
		file {
			type     						= "SourceBuild"
			source_build {
				id 							= ncloud_sourcebuild_project.test-build-project.id
			}
		}
		rollback 							= true
		cmd {
			pre {
				user  						= "root"
				cmd   						= "echo pre"
			}
			deploy {
				source_path 				= "/"
				deploy_path 				= "/test"
			}
			post {
				user  						= "root"
				cmd   						= "echo post"
			}
		}
	}
}

data "ncloud_sourcedeploy_scenario" "scenario"{
	project_id		= ncloud_sourcedeploy_project.project.id
	stage_id		= ncloud_sourcedeploy_stage.svr_stage.id
	id				= ncloud_sourcedeploy_scenario.server_normal.id
}

`, TF_TEST_SD_SERVER_NAME, scenarioNameSvr)
}
