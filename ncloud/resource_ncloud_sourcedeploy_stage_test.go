package ncloud

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vsourcedeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Create Server Before SourceDeploy Test
const TF_TEST_SD_SERVER_NAME = "test-w-wwf"
// Create AutoScalingGroup Before SourceDeploy Test
const TF_TEST_SD_ASG_NAME = "sourcedeploy-bluegreen-12860"
// Create KubernetesService cluster Before SourceDeploy Test
const TF_TEST_SD_NKS_CLUSTER_UUID = "9c9e529a-cc10-42fe-8a34-8f6156456b15"
// Create ObjectStorage bucket cluster Before SourceDeploy Test
const TF_TEST_SD_OBJECTSTORAGE_BUCKET_NAME = "dev"

func TestAccResourceNcloudSourceDeployStage_basic(t *testing.T) {
	var stage vsourcedeploy.GetStageDetailResponse
	stageNameSvr := getTestSourceDeployStageName() + "-svr"
	stageNameAsg := getTestSourceDeployStageName() + "-asg"
	stageNameNks := getTestSourceDeployStageName() + "-nks"
	stageNameObj := getTestSourceDeployStageName() + "-obj"
	resourceNameSvr := "ncloud_sourcedeploy_stage.svr_stage"
	resourceNameAsg := "ncloud_sourcedeploy_stage.asg_stage"
	resourceNameNks := "ncloud_sourcedeploy_stage.nks_stage"
	resourceNameObj := "ncloud_sourcedeploy_stage.obj_stage"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) }, 
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceDeployStageDestroy,
		Steps: []resource.TestStep{ 
			{
				Config: testAccResourceNcloudSourceDeployStageSvrConfig(stageNameSvr, stageNameAsg, stageNameNks, stageNameObj),
				Check: resource.ComposeTestCheckFunc( 
					testAccCheckSourceDeployStageExists(resourceNameSvr, &stage),
					testAccCheckSourceDeployStageExists(resourceNameAsg, &stage),
					testAccCheckSourceDeployStageExists(resourceNameNks, &stage),
					testAccCheckSourceDeployStageExists(resourceNameObj, &stage),
					resource.TestCheckResourceAttr(resourceNameSvr, "name", stageNameSvr), 
					resource.TestCheckResourceAttr(resourceNameAsg, "name", stageNameAsg), 
					resource.TestCheckResourceAttr(resourceNameNks, "name", stageNameNks),
					resource.TestCheckResourceAttr(resourceNameObj, "name", stageNameObj),
				),
			},
		},
	})
}



func testAccCheckSourceDeployStageExists(n string, stage *vsourcedeploy.GetStageDetailResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No stage no is set")
		}
		projectId := ncloud.String(rs.Primary.Attributes["project_id"])
		stageId := &rs.Primary.ID
		resp, err := getSourceDeployStageById(context.Background(), config, projectId, stageId )
		if err != nil {
			return err
		}
		stage = resp
		return nil
	}
} 

func testAccCheckSourceDeployStageDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig) 

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_sourcedeploy_stage" {
			continue
		}

		projectId := ncloud.String(rs.Primary.Attributes["project_id"])
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
	}

	return nil
}


func testAccResourceNcloudSourceDeployStageSvrConfig(stageNameSvr string, stageNameAsg string, stageNameNks string, stageNameObj string) string {
	return fmt.Sprintf(`
data "ncloud_server" "server" {
	filter {
		name = "name"
		values = ["%[1]s"]
	}
}

data "ncloud_auto_scaling_group" "asg" {
	filter{
		name    = "name"
		values  = ["%[2]s"]
	}
}
resource "ncloud_sourcedeploy_project" "project" {
	name    							= "tf-test-project"
}

resource "ncloud_sourcedeploy_stage" "svr_stage" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	name    							= "%[5]s"
	type    							= "Server"
	config {
		server_no  						= [data.ncloud_server.server.id]
	}
}
resource "ncloud_sourcedeploy_stage" "asg_stage" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	name    							= "%[6]s"
	type    							= "AutoScalingGroup"
	config {
		auto_scaling_group_no  			= data.ncloud_auto_scaling_group.asg.id
	}
}
resource "ncloud_sourcedeploy_stage" "nks_stage" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	name    							= "%[7]s"
	type    							= "KubernetesService"
	config {
		cluster_uuid   					= "%[3]s"
	}
}
resource "ncloud_sourcedeploy_stage" "obj_stage" {
	project_id  						= ncloud_sourcedeploy_project.project.id
	name    							= "%[8]s"
	type    							= "ObjectStorage"
	config {
	  bucket_name  						= "%[4]s"
	}
}
`, TF_TEST_SD_SERVER_NAME, TF_TEST_SD_ASG_NAME, TF_TEST_SD_NKS_CLUSTER_UUID, TF_TEST_SD_OBJECTSTORAGE_BUCKET_NAME, 
stageNameSvr, stageNameAsg, stageNameNks, stageNameObj)
}

func getTestSourceDeployStageName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testStageName := fmt.Sprintf("tf-%d-stage", rInt)
	return testStageName
}
