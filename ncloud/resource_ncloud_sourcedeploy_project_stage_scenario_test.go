package ncloud

import (
	"context"
	"fmt"
	"testing"
	"strconv"
	"log"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcedeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)
// Create Load Balancer target group Before SourceDeploy BlueGreen Test
const TF_TEST_SD_LOAD_BALANCER_TARGET_GROUP_NO = "0"
// Setting up prometheus in NKS Before SourceDeploy-Canary-Auto Test
const TF_TEST_SD_PROMETHEUS_URL = "http://prometheus-example.com"

func TestAccResourceNcloudSourceDeployScenario_basic(t *testing.T) {
	var scenario vsourcedeploy.GetScenarioDetailResponse
	scenarioNameSvrNormal := getTestSourceDeployScenarioName() + "-server-normal"
	scenarioNameAsgNoraml := getTestSourceDeployScenarioName() + "-asg-normal"
	scenarioNameAsgBg := getTestSourceDeployScenarioName() + "-asg-bg"
	scenarioNameNksRolling := getTestSourceDeployScenarioName() + "-nks-rolling"
	scenarioNameNksBg := getTestSourceDeployScenarioName() + "-nks-bg"
	scenarioNameNksCanaryManual := getTestSourceDeployScenarioName() + "-nks-canary-manual"
	scenarioNameNksCanaryAuto := getTestSourceDeployScenarioName() + "-nks-canary-auto"
	scenarioNameObjNormal := getTestSourceDeployScenarioName() + "-obj-normal"

	resourceNameSvrNormal := "ncloud_sourcedeploy_project_stage_scenario.server_normal"
	resourceNameAsgNormal := "ncloud_sourcedeploy_project_stage_scenario.asg_normal"
	resourceNameAsgBg := "ncloud_sourcedeploy_project_stage_scenario.asg_bg"
	resourceNameNksRolling := "ncloud_sourcedeploy_project_stage_scenario.nks_rolling"
	resourceNameNksBg := "ncloud_sourcedeploy_project_stage_scenario.nks_bg"
	resourceNameNksCanaryManual := "ncloud_sourcedeploy_project_stage_scenario.nks_canary_manual"
	resourceNameNksCanaryAuto := "ncloud_sourcedeploy_project_stage_scenario.nks_canary_auto"
	resourceNameObjNormal := "ncloud_sourcedeploy_project_stage_scenario.obj_normal"


	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) }, 
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceDeployScenarioDestroy,
		Steps: []resource.TestStep{ 
			{
				Config: testAccResourceNcloudSourceDeployScenarioConfig(
					scenarioNameSvrNormal, 
					scenarioNameAsgNoraml, 
					scenarioNameAsgBg, 
					scenarioNameNksRolling,
					scenarioNameNksBg, 
					scenarioNameNksCanaryManual, 
					scenarioNameNksCanaryAuto,
					scenarioNameObjNormal,
				),
				Check: resource.ComposeTestCheckFunc( 
					testAccCheckSourceDeployScenarioExists(resourceNameSvrNormal, &scenario),
					testAccCheckSourceDeployScenarioExists(resourceNameAsgNormal, &scenario),
					testAccCheckSourceDeployScenarioExists(resourceNameAsgBg, &scenario),
					testAccCheckSourceDeployScenarioExists(resourceNameNksRolling, &scenario),
					testAccCheckSourceDeployScenarioExists(resourceNameNksBg, &scenario),
					testAccCheckSourceDeployScenarioExists(resourceNameObjNormal, &scenario),
					testAccCheckSourceDeployScenarioExists(resourceNameNksCanaryManual, &scenario),
					testAccCheckSourceDeployScenarioExists(resourceNameNksCanaryAuto, &scenario),
					resource.TestCheckResourceAttr(resourceNameSvrNormal, "name", scenarioNameSvrNormal),
					resource.TestCheckResourceAttr(resourceNameAsgNormal, "name", scenarioNameAsgNoraml),
					resource.TestCheckResourceAttr(resourceNameAsgBg, "name", scenarioNameAsgBg),
					resource.TestCheckResourceAttr(resourceNameNksRolling, "name", scenarioNameNksRolling),
					resource.TestCheckResourceAttr(resourceNameNksBg, "name", scenarioNameNksBg),
					resource.TestCheckResourceAttr(resourceNameNksCanaryManual, "name", scenarioNameNksCanaryManual),
					resource.TestCheckResourceAttr(resourceNameNksCanaryAuto, "name", scenarioNameNksCanaryAuto),
					resource.TestCheckResourceAttr(resourceNameObjNormal, "name", scenarioNameObjNormal),


				),
			},
		},
	})
}

func testAccResourceNcloudSourceDeployScenarioConfig(
	scenarioNameSvrNormal string, 
	scenarioNameAsgNoraml string, 
	scenarioNameAsgBg string, 
	scenarioNameNksRolling string,
	scenarioNameNksBg string, 
	scenarioNameNksCanaryManual string, 
	scenarioNameNksCanaryAuto string,
	scenarioNameObjNormal string ) string {
	return fmt.Sprintf(`
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

data "ncloud_server" "server" {
	filter {
		name 				= "name"
		values				= ["%[1]s"]
	}
}

data "ncloud_auto_scaling_group" "asg" {
	filter{
		name  				= "name"
		values  			= ["%[2]s"]
	}
}

resource "ncloud_sourcedeploy_project" "project" {
	name    								= "tf-test-project"
}

resource "ncloud_sourcedeploy_project_stage" "svr_stage" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	name    								= "svr"
	type    								= "Server"
	config {
		server_no  							= [data.ncloud_server.server.id]
	}
}
resource "ncloud_sourcedeploy_project_stage" "asg_stage" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	name    								= "asg"
	type    								= "AutoScalingGroup"
	config {
		auto_scaling_group_no  				= data.ncloud_auto_scaling_group.asg.id
	}
}
resource "ncloud_sourcedeploy_project_stage" "nks_stage" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	name    								= "nks"
	type    								= "KubernetesService"
	config {
		cluster_uuid   						= "%[3]s"
	}
}
resource "ncloud_sourcedeploy_project_stage" "obj_stage" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	name    								= "obj"
	type    								= "ObjectStorage"
	config {
		bucket_name  						= "%[4]s"
	}
}

resource "ncloud_sourcedeploy_project_stage_scenario" "server_normal" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.svr_stage.id
	name    								= "%[5]s"
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


resource "ncloud_sourcedeploy_project_stage_scenario" "asg_normal" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.asg_stage.id
	name    								= "%[6]s"
	description   	 						= "test"
	config {
		strategy  							= "normal"
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

resource "ncloud_sourcedeploy_project_stage_scenario" "asg_bg" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.asg_stage.id
	name    								= "%[7]s"
	description   	 						= "test"
	config {
		strategy  							= "blueGreen"
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
		load_balancer{
			load_balancer_target_group_no 	= "%[8]s"
			delete_server 					= true
		}
	}
}

resource "ncloud_sourcedeploy_project_stage_scenario" "nks_rolling" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.nks_stage.id
	name    								= "%[9]s"
	description    							= "test"
	config {
		strategy  							= "rolling"
		manifest {
			type    						= "SourceCommit"
			repository 						= ncloud_sourcecommit_repository.test-repo.name
			branch    						= "master"
			path      						= ["/deployment/prod.yaml"]
		}
	}
}

resource "ncloud_sourcedeploy_project_stage_scenario" "nks_bg" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.nks_stage.id
	name    								= "%[10]s"
	description    							= "test"
	config {	
		strategy  							= "blueGreen"
		manifest {
			type     						= "SourceCommit"
			repository 						= ncloud_sourcecommit_repository.test-repo.name
			branch    						= "master"
			path      						= ["/deployment/canary.yaml"]
		}
	}
}

resource "ncloud_sourcedeploy_project_stage_scenario" "nks_canary_manual" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id   								= ncloud_sourcedeploy_project_stage.nks_stage.id
	name    								= "%[11]s"
	description    							= "test"
	config {
		strategy  							= "canary"
		manifest {
			type     						= "SourceCommit"
			repository	 					= ncloud_sourcecommit_repository.test-repo.name
			branch    						= "master"
			path      						= ["/deployment/canary.yaml"]
		}
		canary_config{
			analysis_type  					=   "manual"
			timeout       					=   10
			canary_count  					=   1
		}
	}
}

 resource "ncloud_sourcedeploy_project_stage_scenario" "nks_canary_auto" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	stage_id   							= ncloud_sourcedeploy_project_stage.nks_stage.id
	name    							= "%[12]s"
	description    						= "test"
	config {
		strategy  						= "canary"
		manifest {
			type     					= "SourceCommit"
			repository 					= ncloud_sourcecommit_repository.test-repo.name
			branch    					= "master"
			path     					= ["test.yaml"]
		}
		canary_config{
			analysis_type  				= "auto"
			canary_count  				= 1
			prometheus    				= "%[13]s"
			env{
				baseline 				= "baselineenv"
				canary    				= "canaryenv"
			}
			metrics{
				name      				= "success_rate"
				success_criteria  		= "base"
				query_type     	 		= "promQL"
				weight    				= 100
				query   				= "test"
			}
			analysis_config{
				duration  				= 10
				delay    				= 1
				interval  				= 1
				step      				= 10
			}
			pass_score					= 90
		}
	}
 }

resource "ncloud_sourcedeploy_project_stage_scenario" "obj_normal" {
	project_id  							= ncloud_sourcedeploy_project.project.id
	stage_id    							= ncloud_sourcedeploy_project_stage.obj_stage.id
	name    								= "%[14]s"
	description    							= "test"
	config {
		file {
			type     						= "SourceBuild"
			source_build {
				id 							= ncloud_sourcebuild_project.test-build-project.id
			}
		}
		path {
			source_path       				=   "/"
			deploy_path       				=   "/terraform"
		}
	}
}
`, TF_TEST_SD_SERVER_NAME, TF_TEST_SD_ASG_NAME, TF_TEST_SD_NKS_CLUSTER_UUID, TF_TEST_SD_OBJECTSTORAGE_BUCKET_NAME, 
scenarioNameSvrNormal, scenarioNameAsgNoraml, scenarioNameAsgBg, TF_TEST_SD_LOAD_BALANCER_TARGET_GROUP_NO, scenarioNameNksRolling, 
scenarioNameNksBg, scenarioNameNksCanaryManual, scenarioNameNksCanaryAuto, TF_TEST_SD_PROMETHEUS_URL, 
scenarioNameObjNormal )
}


func testAccCheckSourceDeployScenarioExists(n string, scenario *vsourcedeploy.GetScenarioDetailResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No scenario no is set")
		}
		projectId := ncloud.String(rs.Primary.Attributes["project_id"])
		stageId := ncloud.String(rs.Primary.Attributes["stage_id"])
		scenarioId := &rs.Primary.ID
		resp, err := getSourceDeployScenarioById(context.Background(), config, projectId, stageId, scenarioId)
		if err != nil {
			return err
		}
		scenario = resp
		return nil
	}
} 

func testAccCheckSourceDeployScenarioDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		log.Printf(rs.Type)
		if rs.Type != "ncloud_sourcedeploy_project_stage_scenario" {
			continue
		}
		projectId := ncloud.String(rs.Primary.Attributes["project_id"])
		stageId := ncloud.String(rs.Primary.Attributes["stage_id"])
		project, projectErr := getSourceDeployProjectById(context.Background(), config, rs.Primary.ID)
		if projectErr != nil {
			return projectErr
		}

		if project == nil{
			return nil
		}
		
		stages, stageErr := getStages(context.Background(), config, projectId )
		if stageErr != nil {
			return stageErr
		}

		for _, stage := range stages.StageList {
			if strconv.Itoa(int(ncloud.Int32Value(stage.Id))) == rs.Primary.ID {
				return fmt.Errorf("stage still exists")
			}
		}

		scenarios, scenarioErr := GetScenarioes(context.Background(), config, projectId, stageId )
		if scenarioErr != nil {
			return scenarioErr
		}

		for _, scenario := range scenarios.ScenarioList {
			if strconv.Itoa(int(ncloud.Int32Value(scenario.Id))) == rs.Primary.ID {
				return fmt.Errorf("scenario still exists")
			}
		}
	}

	return nil
}

func getTestSourceDeployScenarioName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testScenarioName := fmt.Sprintf("tf-%d-scenario", rInt)
	return testScenarioName
}
