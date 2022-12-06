package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourcePipelineProject_classic_basic(t *testing.T) {
	dataName := "data.ncloud_sourcepipeline_project.foo"
	resourceName := "ncloud_sourcepipeline_project.test-project"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccClassicProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineProjectConfig("test-project", "description test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudSourcePipelineProject_vpc_basic(t *testing.T) {
	dataName := "data.ncloud_sourcepipeline_project.foo"
	resourceName := "ncloud_sourcepipeline_project.test-project"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourcePipelineProjectConfig("test-project", "description test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourcePipelineProjectConfig(name, description string) string {
	return fmt.Sprintf(`
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
	name = "sourceCommit"
}

resource "ncloud_sourcebuild_project" "test-project" {
	name        = "sourceBuild"
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
			id = data.ncloud_sourcebuild_project_docker_engines.docker_engines.docker_engines[0].id
		}
		timeout = 500
		env_var {
			key   = "k1"
			value = "v1"
		}
	}
	build_command {
		pre_build   = ["pwd", "ls"]
		in_build = ["pwd", "ls"]
		post_build  = ["pwd", "ls"]
	}
}

resource "ncloud_sourcepipeline_project" "test-project" {
	name               = "%[1]s"
	description        = "%[2]s"
	task {
		name 		   = "task_name"
		type 		   = "SourceBuild"
		config {
			project_id   = ncloud_sourcebuild_project.test-project.id
		}
		linked_tasks   = []
	}
	triggers {
		schedule {
            day                       = ["MON", "TUE"]
            time                      = "13:01"
            timezone                  = "Asia/Seoul (UTC+09:00)"
            execute_only_with_change = false
        }
	}
}

data "ncloud_sourcepipeline_project" "foo" {
	id = ncloud_sourcepipeline_project.test-project.id
}
`, name, description)
}
