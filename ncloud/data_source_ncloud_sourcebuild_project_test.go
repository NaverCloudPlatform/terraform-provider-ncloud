package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceBuildProject(t *testing.T) {
	name := fmt.Sprintf("test-sourcebuild-project-name-%s", acctest.RandString(5))
	repoName := fmt.Sprintf("test-repo-basic-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceBuildProjectConfig(name, repoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID("data.ncloud_sourcebuild_project.project"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceBuildProjectConfig(name string, repoName string) string {
	return fmt.Sprintf(`
data "ncloud_sourcebuild_project_compute" "compute" {
}

data "ncloud_sourcebuild_project_os" "os" {
}

data "ncloud_sourcebuild_project_runtime" "runtime" {
	os_id = data.ncloud_sourcebuild_project_os.os.os[0].id
}

data "ncloud_sourcebuild_project_runtime_version" "runtime_version" {
	os_id      = data.ncloud_sourcebuild_project_os.os.os[0].id
	runtime_id = data.ncloud_sourcebuild_project_runtime.runtime.runtime[0].id
}

data "ncloud_sourcebuild_project_docker" "docker" {
}
	  
resource "ncloud_sourcecommit_repository" "test-repo" {
	name = "%[1]s"
}

resource "ncloud_sourcebuild_project" "test-project" {
	name        = "%[2]s"
	description = "my build project"
	source {
		type = "SourceCommit"
		config {
			repository = ncloud_sourcecommit_repository.test-repo.name
			branch     = "master"
		}
	}
	env {
		compute {
			id = data.ncloud_sourcebuild_project_compute.compute.compute[0].id
		}
		platform {
			type = "SourceBuild"
			config {
				os {
					id = data.ncloud_sourcebuild_project_os.os.os[0].id
				}
				runtime {
					id = data.ncloud_sourcebuild_project_runtime.runtime.runtime[0].id
					version {
						id = data.ncloud_sourcebuild_project_runtime_version.runtime_version.runtime_version[0].id
					}
				}
			}
		}
	}
}

data "ncloud_sourcebuild_project" "project" {
	id = ncloud_sourcebuild_project.test-project.id
}
`, repoName, name)
}
