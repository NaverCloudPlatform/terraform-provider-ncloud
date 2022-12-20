package ncloud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcebuild"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNcloudSourcebuildProject_basic(t *testing.T) {
	var project sourcebuild.GetProjectDetailResponse
	name := fmt.Sprintf("test-sourcebuild-project-basic-%s", acctest.RandString(5))
	repoName := fmt.Sprintf("test-repo-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcebuild_project.test-project"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourcebuildProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcebuildConfig(name, repoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcebuildProjectExists(resourceName, &project),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudSourcebuildProject_update(t *testing.T) {
	var project sourcebuild.GetProjectDetailResponse
	name := fmt.Sprintf("test-sourcebuild-project-name-%s", acctest.RandString(5))
	repoName := fmt.Sprintf("test-repo-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcebuild_project.test-project"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourcebuildProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcebuildConfig(name, repoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcebuildProjectExists(resourceName, &project),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccResourceNcloudSourcebuildUpdatedConfig(name, repoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcebuildProjectExists(resourceName, &project),
					resource.TestCheckResourceAttr(resourceName, "env.0.timeout", "100"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceNcloudSourcebuildConfig(name string, repoName string) string {
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
	name = "%[1]s"
}

resource "ncloud_sourcebuild_project" "test-project" {
	name        = "%[2]s"
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
}`, repoName, name)
}

func testAccResourceNcloudSourcebuildUpdatedConfig(name string, repoName string) string {
	return fmt.Sprintf(`
data "ncloud_sourcebuild_project_computes" "computes" {
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
			repository_name = ncloud_sourcecommit_repository.test-repo.name
			branch     = "master"
		}
	}
	env {
		compute {
			id = data.ncloud_sourcebuild_project_computes.computes.computes[0].id
		}
		platform {
			type = "PublicRegistry"
			config {
				image    = "ubuntu"
				tag      = "latest"
			}
		}
		timeout = 100
		env_var {
			key   = "k2"
			value = "v2"
		}
	}
	build_command {
		pre_build   = [""]
		in_build = ["pwd", "ls"]
		post_build  = ["pwd", "ls"]
	}
}`, repoName, name)
}

func testAccCheckSourcebuildProjectExists(n string, project *sourcebuild.GetProjectDetailResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project no is set")
		}

		buildProject, err := getSourceBuildProject(config, &rs.Primary.ID)
		if err != nil {
			return err
		}

		project = buildProject

		return nil
	}
}

func testAccCheckSourcebuildProjectDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_sourcebuild_project" {
			continue
		}

		buildProject, err := getSourceBuildProject(config, &rs.Primary.ID)

		if err != nil {
			return err
		}

		if buildProject != nil {
			return errors.New("project still exists")
		}
	}

	return nil
}

func getSourceBuildProject(config *ProviderConfig, id *string) (*sourcebuild.GetProjectDetailResponse, error) {
	logCommonRequest("getProjectDetail", id)
	//This api throws an error when the resource cannot be found.
	resp, err := config.Client.sourcebuild.V1Api.GetProject(context.Background(), id)

	//Don't throw an error when the error is 'resource not found' to continue executing business logic.
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		} else {
			logErrorResponse("getProjectDetail", err, id)
			return nil, err
		}
	}

	logResponse("getProjectDetail", resp)

	return resp, nil
}
