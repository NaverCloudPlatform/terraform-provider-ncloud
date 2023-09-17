package devtools_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/devtools"
)

func TestAccResourceNcloudSourcePipelineProject_classic_basic(t *testing.T) {
	var project devtools.PipelineProject
	name := fmt.Sprintf("test-pipeline-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcepipeline_project.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckSourcePipelineProjectDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcePipelineProjectClassicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(false)),
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

func TestAccResourceNcloudSourcePipelineProject_classic_updateTaskName(t *testing.T) {
	var project devtools.PipelineProject
	name := fmt.Sprintf("test-pipeline-name-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcepipeline_project.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckSourcePipelineProjectDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcePipelineProjectClassicConfig(name),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(false)),
				),
			},
			{
				Config: testAccResourceNcloudSourcePipelineProjectClassicConfigUpdateTaskName(name, "updated_task_name"),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "task.0.name", "updated_task_name"),
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

func TestAccResourceNcloudSourcePipelineProject_classic_updateDescription(t *testing.T) {
	var project devtools.PipelineProject
	name := fmt.Sprintf("test-pipeline-name-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcepipeline_project.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckSourcePipelineProjectDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcePipelineProjectClassicConfig(name),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(false)),
				),
			},
			{
				Config: testAccResourceNcloudSourcePipelineProjectClassicConfigUpdateDescription(name, "updatedDescription"),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(false)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "updatedDescription"),
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

func TestAccResourceNcloudSourcePipelineProject_vpc_basic(t *testing.T) {
	var project devtools.PipelineProject
	name := fmt.Sprintf("test-pipeline-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcepipeline_project.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckSourcePipelineProjectDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcePipelineProjectVpcConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(true)),
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

func TestAccResourceNcloudSourcePipelineProject_vpc_updateTaskName(t *testing.T) {
	var project devtools.PipelineProject
	name := fmt.Sprintf("test-pipeline-name-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcepipeline_project.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckSourcePipelineProjectDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcePipelineProjectVpcConfig(name),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(true)),
				),
			},
			{
				Config: testAccResourceNcloudSourcePipelineProjectVpcConfigUpdateTaskName(name, "updated_task_name"),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "task.0.name", "updated_task_name"),
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

func TestAccResourceNcloudSourcePipelineProject_vpc_updateDescription(t *testing.T) {
	var project devtools.PipelineProject
	name := fmt.Sprintf("test-pipeline-name-%s", acctest.RandString(5))
	resourceName := "ncloud_sourcepipeline_project.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckSourcePipelineProjectDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudSourcePipelineProjectVpcConfig(name),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(true)),
				),
			},
			{
				Config: testAccResourceNcloudSourcePipelineProjectVpcConfigUpdateDescription(name, "updatedDescription"),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourcePipelineProjectExists(resourceName, &project, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "updatedDescription"),
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

func testAccResourceNcloudSourcePipelineProjectClassicConfig(name string) string {
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

resource "ncloud_sourcepipeline_project" "foo" {
	name               = "%[1]s"
	description        = "test pipeline project"
	task {
		name 		   = "task_name"
		type 		   = "SourceBuild"
		config {
		  project_id   = ncloud_sourcebuild_project.test-project.id
		}
		linked_tasks   = []
	}
	triggers {
		repository {
			type = "sourcecommit"
			name = "sourceCommit"
			branch = "master"
		}
	}
}
`, name)
}

func testAccResourceNcloudSourcePipelineProjectClassicConfigUpdateTaskName(name, taskName string) string {
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
			branch     = "master"
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

resource "ncloud_sourcepipeline_project" "foo" {
	name               = "%[1]s"
	description        = "test pipeline project"
	task {
		name 		   = "%[2]s"
		type 		   = "SourceBuild"
		config {
			project_id   = ncloud_sourcebuild_project.test-project.id
		}
		linked_tasks   = []
	}
	triggers {
		repository {
			type = "sourcecommit"
			name = "sourceCommit"
			branch = "master"
		}
	}
}
`, name, taskName)
}

func testAccResourceNcloudSourcePipelineProjectClassicConfigUpdateDescription(name, description string) string {
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
			branch     = "master"
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

resource "ncloud_sourcepipeline_project" "foo" {
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
		repository {
			type = "sourcecommit"
			name = "sourceCommit"
			branch = "master"
		}
	}
}
`, name, description)
}

func testAccResourceNcloudSourcePipelineProjectVpcConfig(name string) string {
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
			branch     = "master"
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

resource "ncloud_sourcepipeline_project" "foo" {
	name               = "%[1]s"
	description        = "test pipeline project"
	task {
		name 		   = "task_name"
		type 		   = "SourceBuild"
		config {
			project_id   = ncloud_sourcebuild_project.test-project.id
		}
		linked_tasks   = []
	}
	triggers {
		repository {
			type = "sourcecommit"
			name = "sourceCommit"
			branch = "master"
		}
	}
}
`, name)
}

func testAccResourceNcloudSourcePipelineProjectVpcConfigUpdateTaskName(name, taskName string) string {
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
			branch     = "master"
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

resource "ncloud_sourcepipeline_project" "foo" {
	name               = "%[1]s"
	description        = "test pipeline project"
	task {
		name 		   = "%[2]s"
		type 		   = "SourceBuild"
		config {
			project_id   = ncloud_sourcebuild_project.test-project.id
		}
		linked_tasks   = []
	}
	triggers {
		repository {
			type = "sourcecommit"
			name = "sourceCommit"
			branch = "master"
		}
	}
}
`, name, taskName)
}

func testAccResourceNcloudSourcePipelineProjectVpcConfigUpdateDescription(name, description string) string {
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
			branch     = "master"
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

resource "ncloud_sourcepipeline_project" "foo" {
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
		repository {
			type = "sourcecommit"
			name = "sourceCommit"
			branch = "master"
		}
	}
}
`, name, description)
}

func testAccCheckSourcePipelineProjectExists(n string, project *devtools.PipelineProject, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		pipelineProject, err := devtools.GetPipelineProject(context.Background(), config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*project = *pipelineProject

		return nil
	}
}

func testAccCheckSourcePipelineProjectDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_sourcepipeline_project" {
			continue
		}

		pipelineProject, _ := devtools.GetPipelineProject(context.Background(), config, rs.Primary.ID)

		if pipelineProject != nil {
			return errors.New("SourcePipeline project still exists")
		}
	}

	// FIXME: When testing, time gap is required between each test to create
	// resource without any issues. It should guarantee to deletion and creation
	// of the resource without problem from the PipelineProject API
	time.Sleep(1 * time.Minute)

	return nil
}
