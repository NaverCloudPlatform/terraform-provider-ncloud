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
					testAccCheckDataSourceID("data.ncloud_sourcedeploy_project_stage_scenario.scenario"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceDeployScenarioConfig(scenarioNameSvr string) string {
	return fmt.Sprintf(`
data "ncloud_server" "server" {
	filter {
		name   = "name"
		values = ["%[1]s"]
	}
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

resource "ncloud_sourcecommit_repository" "test-repo" {
	name = "tf-test-repository"
}

resource "ncloud_sourcebuild_project" "test-build-project" {
	name        = "tf-test-project"
	description = "my build project"
	source {
		type = "SourceCommit"
		config {
			repository_name = ncloud_sourcecommit_repository.test-repo.name
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
		timeout = 500
		env_var {
			key   = "k1"
			value = "v1"
		}
	}
	build_command {
		pre_build  = ["pwd", "ls"]
		in_build   = ["pwd", "ls"]
		post_build = ["pwd", "ls"]
	}
}

resource "ncloud_sourcedeploy_project" "project" {
	name = "tf-test-project"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id = ncloud_sourcedeploy_project.project.id
	name       = "svr"
	target_type = "Server"
	config {
		server{
			id = data.ncloud_server.server.id
		} 
	}
}

resource "ncloud_sourcedeploy_project_stage_scenario" "server_normal" {
	project_id  = ncloud_sourcedeploy_project.project.id
	stage_id    = ncloud_sourcedeploy_project_stage.svr_stage.id
	name        = "%[2]s"
	description = "test"
	config {
		strategy = "normal"
		file {
			type = "SourceBuild"
			source_build {
				id = ncloud_sourcebuild_project.test-build-project.id
			}
		}
		rollback = true
		deploy_command {
			pre_deploy {
				user     = "root"
				command  = "echo pre"
			}
			path {
				source_path = "/"
				deploy_path = "/test"
			}
			post_deploy {
				user    = "root"
				command = "echo post"
			}
		}
	}
}

data "ncloud_sourcedeploy_project_stage_scenario" "scenario"{
	project_id = ncloud_sourcedeploy_project.project.id
	stage_id   = ncloud_sourcedeploy_project_stage.svr_stage.id
	id         = ncloud_sourcedeploy_project_stage_scenario.server_normal.id
}

`, TF_TEST_SD_SERVER_NAME, scenarioNameSvr)
}
